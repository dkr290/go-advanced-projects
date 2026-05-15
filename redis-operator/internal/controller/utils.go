package controller

import (
	"bytes"
	"context"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// reconcileRoles checks master/replica health and performs failover if needed.
// Returns true if a requeue is needed.
func (r *BcredisReconciler) reconcileRoles(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
) (bool, error) {
	log := logf.FromContext(ctx)
	ns := bcredis.Namespace

	pod0Name := fmt.Sprintf("%s-redis-0-0", bcredis.Name)
	pod1Name := fmt.Sprintf("%s-redis-1-0", bcredis.Name)

	pod0 := &corev1.Pod{}
	pod1 := &corev1.Pod{}

	err0 := r.Get(ctx, types.NamespacedName{Name: pod0Name, Namespace: ns}, pod0)
	err1 := r.Get(ctx, types.NamespacedName{Name: pod1Name, Namespace: ns}, pod1)

	// Pods not yet created
	if err0 != nil || err1 != nil {
		bcredis.Status.Phase = "Initializing"
		return true, nil
	}

	pod0Ready := isPodReady(pod0)
	pod1Ready := isPodReady(pod1)

	// Both not ready yet
	if !pod0Ready && !pod1Ready {
		bcredis.Status.Phase = "Initializing"
		return true, nil
	}

	currentMaster := bcredis.Status.MasterPod
	currentReplica := bcredis.Status.ReplicaPod

	// Initial setup: assign roles
	if currentMaster == "" {
		if pod0Ready {
			log.Info("Initial role assignment: pod0=master, pod1=replica")
			if err := r.promoteToMaster(ctx, bcredis, pod0Name, spec); err != nil {
				return true, err
			}
			if pod1Ready {
				if err := r.configureReplica(ctx, bcredis, pod1Name, pod0Name, spec); err != nil {
					return true, err
				}
				bcredis.Status.ReplicaPod = pod1Name
			}
			bcredis.Status.MasterPod = pod0Name
			bcredis.Status.CurrentMasterService = fmt.Sprintf("%s-redis-0", bcredis.Name)
			bcredis.Status.Phase = "Running"
			if err := r.updateTCPRoute(
				ctx,
				bcredis,
				bcredis.Status.CurrentMasterService,
				spec,
			); err != nil {
				return true, err
			}
			return false, nil
		} else if pod1Ready {
			log.Info("Initial role assignment: pod1=master, pod0=replica (pod0 not ready)")
			if err := r.promoteToMaster(ctx, bcredis, pod1Name, spec); err != nil {
				return true, err
			}
			bcredis.Status.MasterPod = pod1Name
			bcredis.Status.CurrentMasterService = fmt.Sprintf("%s-redis-1", bcredis.Name)
			bcredis.Status.Phase = "Running"
			if err := r.updateTCPRoute(
				ctx,
				bcredis,
				bcredis.Status.CurrentMasterService,
				spec,
			); err != nil {
				return true, err
			}
			return true, nil
		}
	}

	// Ongoing health check: detect master failure and failover
	masterReady := (currentMaster == pod0Name && pod0Ready) ||
		(currentMaster == pod1Name && pod1Ready)
	replicaReady := (currentReplica == pod0Name && pod0Ready) ||
		(currentReplica == pod1Name && pod1Ready)

	if !masterReady && replicaReady {
		log.Info(
			"Master is not ready, promoting replica",
			"master",
			currentMaster,
			"replica",
			currentReplica,
		)
		bcredis.Status.Phase = "Failover"

		// Promote current replica to master
		if err := r.promoteToMaster(ctx, bcredis, currentReplica, spec); err != nil {
			return true, fmt.Errorf("failed to promote replica: %w", err)
		}

		// Update TCPRoute to point to new master service
		newMasterSvc := r.podNameToServiceName(bcredis, currentReplica)
		if err := r.updateTCPRoute(ctx, bcredis, newMasterSvc, spec); err != nil {
			return true, fmt.Errorf("failed to update TCPRoute: %w", err)
		}

		// Swap roles
		bcredis.Status.MasterPod = currentReplica
		bcredis.Status.ReplicaPod = currentMaster
		bcredis.Status.CurrentMasterService = newMasterSvc
		bcredis.Status.Phase = "Running"

		log.Info(
			"Failover complete",
			"newMaster",
			bcredis.Status.MasterPod,
			"newReplica",
			bcredis.Status.ReplicaPod,
		)
		return true, nil
	}

	// If old master came back, configure it as replica of current master
	if !masterReady && !replicaReady {
		bcredis.Status.Phase = "Degraded"
		return true, nil
	}

	// Old master recovered — set it as replica if it's running but not assigned a role properly
	masterPod := bcredis.Status.MasterPod
	replicaPod := bcredis.Status.ReplicaPod

	if replicaPod != "" {
		replicaPodObj := &corev1.Pod{}
		if err := r.Get(
			ctx,
			types.NamespacedName{Name: replicaPod, Namespace: ns},
			replicaPodObj,
		); err == nil {
			if isPodReady(replicaPodObj) {
				// Check if replica is already replicating from master
				role, err := r.getRedisRole(ctx, replicaPod, ns)
				if err == nil && role != "slave" && role != "replica" {
					log.Info(
						"Reconfiguring recovered pod as replica",
						"pod",
						replicaPod,
						"master",
						masterPod,
					)
					if err := r.configureReplica(
						ctx,
						bcredis,
						replicaPod,
						masterPod,
						spec,
					); err != nil {
						log.Error(err, "failed to reconfigure recovered pod as replica")
					}
				}
			}
		}
	}

	bcredis.Status.Phase = "Running"
	return false, nil
}

// promoteToMaster runs REPLICAOF NO ONE on the target pod.
func (r *BcredisReconciler) promoteToMaster(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	podName string,
	spec bcredisv1alpha1.BcredisSpec,
) error {
	cmd := []string{"redis-cli", "REPLICAOF", "NO", "ONE"}
	if spec.RedisPasswordSecret != "" {
		cmd = append([]string{"redis-cli", "-a", "$(REDIS_PASSWORD)"}, "REPLICAOF", "NO", "ONE")
	}
	_, _, err := r.execInPod(ctx, bcredis.Namespace, podName, "redis", cmd)
	return err
}

// configureReplica runs REPLICAOF <masterHost> 6379 on the target pod.
func (r *BcredisReconciler) configureReplica(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	podName, masterPodName string,
	spec bcredisv1alpha1.BcredisSpec,
) error {
	masterSvcHost := r.podNameToServiceName(bcredis, masterPodName)
	cmd := []string{"redis-cli", "REPLICAOF", masterSvcHost, "6379"}
	_, _, err := r.execInPod(ctx, bcredis.Namespace, podName, "redis", cmd)
	return err
}

// getRedisRole returns the role of a Redis pod ("master", "slave", "replica").
func (r *BcredisReconciler) getRedisRole(ctx context.Context, podName, ns string) (string, error) {
	cmd := []string{"redis-cli", "INFO", "replication"}
	stdout, _, err := r.execInPod(ctx, ns, podName, "redis", cmd)
	if err != nil {
		return "", err
	}
	for _, line := range splitLines(stdout) {
		if len(line) > 5 && line[:5] == "role:" {
			return line[5:], nil
		}
	}
	return "", fmt.Errorf("could not determine role from INFO replication output")
}

// execInPod executes a command in a container within a pod and returns stdout, stderr.
func (r *BcredisReconciler) execInPod(
	ctx context.Context,
	ns, podName, containerName string,
	cmd []string,
) (string, string, error) {
	req := r.KubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(ns).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(r.RestConfig, "POST", req.URL())
	if err != nil {
		return "", "", err
	}

	var stdout, stderr bytes.Buffer
	if err := exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		return "", stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

// updateTCPRoute patches the TCPRoute to point to the given service name.
func (r *BcredisReconciler) updateTCPRoute(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	masterSvcName string,
	spec bcredisv1alpha1.BcredisSpec,
) error {
	bcredis.Status.CurrentMasterService = masterSvcName
	return r.reconcileGateway(ctx, bcredis, spec)
}

// podNameToServiceName derives the service name from a StatefulSet pod name.
// Pod names follow the pattern: <bcredisName>-redis-<idx>-0
func (r *BcredisReconciler) podNameToServiceName(
	bcredis *bcredisv1alpha1.Bcredis,
	podName string,
) string {
	// Pod: myredis-redis-0-0 -> Service: myredis-redis-0
	// Pod: myredis-redis-1-0 -> Service: myredis-redis-1
	if len(podName) > 2 {
		// Strip trailing "-0" (the StatefulSet ordinal)
		return podName[:len(podName)-2]
	}
	return podName
}

// isPodReady returns true if a pod is in the Ready condition.
func isPodReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// sectionNamePtr is a helper to get a pointer to a SectionName.
func sectionNamePtr(s string) *gatewayv1.SectionName {
	sn := gatewayv1.SectionName(s)
	return &sn
}

// splitLines splits a string into lines, trimming carriage returns.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
