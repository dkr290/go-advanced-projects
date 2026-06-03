package controller

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/remotecommand"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	masterStatus     = "master"
	replicaStatus    = "slave"
	failoverCooldown = 60 * time.Second
)

func redisPodName(crName string, idx int) string {
	return fmt.Sprintf("%s-redis-%d-0", crName, idx)
}

func redisServiceName(crName string, idx int) string {
	return fmt.Sprintf("%s-redis-%d", crName, idx)
}

func masterIndex(bcredis *bcredisv1alpha1.Bcredis) int {
	if bcredis.Status.MasterPod == "" {
		return 0
	}
	prefix := fmt.Sprintf("%s-redis-", bcredis.Name)
	s := strings.TrimPrefix(bcredis.Status.MasterPod, prefix)
	s = strings.TrimSuffix(s, "-0")
	idx, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return idx
}

func redisCmd(spec bcredisv1alpha1.BcredisSpec, args ...string) []string {
	if spec.RedisPasswordSecret == "" {
		return append([]string{"redis-cli"}, args...)
	}
	return []string{
		"sh", "-c",
		fmt.Sprintf(`redis-cli -a "$REDIS_PASSWORD" %s`, strings.Join(args, " ")),
	}
}

func findReplicaCandidate(
	ctx context.Context,
	r *BcredisReconciler,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	currentMasterIdx int,
) (int, bool, error) {
	replicas := max(spec.Replicas, int32(3))
	for i := 0; i < int(replicas); i++ {
		if i == currentMasterIdx {
			continue
		}
		ready, err := r.isPodExecReady(ctx, bcredis.Namespace, redisPodName(bcredis.Name, i))
		if err != nil {
			return -1, false, err
		}
		if ready {
			return i, true, nil
		}
	}
	return -1, false, nil
}

// reconcileRoles checks Redis replication status and triggers automatic failover
func (r *BcredisReconciler) reconcileRoles(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
) (bool, error) {
	logger := logf.FromContext(ctx).WithValues("bcredis", bcredis.Name)
	mIdx := masterIndex(bcredis)
	masterReady, err := r.isPodExecReady(ctx, bcredis.Namespace, redisPodName(bcredis.Name, mIdx))
	if err != nil {
		return false, err
	}
	if !masterReady {
		logger.Info(
			"Current master pod is not exec-ready, scanning for existing master",
			"masterIdx",
			mIdx,
		)
		// Sentinel may have already promoted someone — check before acting
		existingMaster := r.findCurrentMaster(ctx, bcredis, spec, logger)
		if existingMaster >= 0 && existingMaster != mIdx {
			logger.Info(
				"Sentinel already promoted a new master, updating status",
				"newMasterIdx",
				existingMaster,
			)
			bcredis.Status.MasterPod = redisPodName(bcredis.Name, existingMaster)
			bcredis.Status.CurrentMasterService = redisServiceName(bcredis.Name, existingMaster)
			bcredis.Status.FailedOver = true
			bcredis.Status.LastFailoverTime = metav1.Now()
			return true, nil
		}
		if existingMaster < 0 {
			logger.Info("No master found anywhere, triggering operator failover", "masterIdx", mIdx)
			return true, r.performFailover(ctx, bcredis, spec, logger)
		}
		// existingMaster == mIdx means pod is coming back up, just requeue
		return true, nil
	}

	masterRole, err := r.getRedisRole(ctx, bcredis, spec, mIdx, logger)
	if err != nil {
		// Role check failed transiently — log and requeue, don't failover
		logger.Info("Could not get role for master, will retry", "masterIdx", mIdx, "err", err)
		return true, err
	}
	if masterRole != masterStatus {
		logger.Info(
			"Tracked master is no longer master, scanning for current master",
			"masterIdx",
			mIdx,
			"role",
			masterRole,
		)
		existingMaster := r.findCurrentMaster(ctx, bcredis, spec, logger)
		if existingMaster >= 0 && existingMaster != mIdx {
			// Sentinel already elected a new master — just track it
			logger.Info(
				"Sentinel already elected a new master, syncing status",
				"newMasterIdx",
				existingMaster,
			)
			bcredis.Status.MasterPod = redisPodName(bcredis.Name, existingMaster)
			bcredis.Status.CurrentMasterService = redisServiceName(bcredis.Name, existingMaster)
			bcredis.Status.FailedOver = true
			bcredis.Status.LastFailoverTime = metav1.Now()
			return true, nil
		}
		if existingMaster < 0 {
			logger.Info("No master found, triggering operator failover")
			return true, r.performFailover(ctx, bcredis, spec, logger)
		}
		// existingMaster == mIdx: role check was transient, do nothing
		return true, nil
	}
	// Initialize status if not set yet
	if bcredis.Status.MasterPod == "" || bcredis.Status.CurrentMasterService == "" {
		bcredis.Status.MasterPod = redisPodName(bcredis.Name, mIdx)
		bcredis.Status.CurrentMasterService = redisServiceName(bcredis.Name, mIdx)
		logger.Info("Initialized master status", "masterPod", bcredis.Status.MasterPod)
	}

	// Ensure all replicas are pointing to the correct master
	masterSvc := redisServiceName(bcredis.Name, mIdx)
	replicas := max(spec.Replicas, int32(3))
	for i := 0; i < int(replicas); i++ {
		if i == mIdx {
			continue
		}
		ready, _ := r.isPodExecReady(ctx, bcredis.Namespace, redisPodName(bcredis.Name, i))
		if !ready {
			continue
		}
		role, err := r.getRedisRole(ctx, bcredis, spec, i, logger)
		if err != nil {
			logger.Info("Could not get role for replica, skipping", "idx", i)
			continue
		}
		if role == masterStatus {
			// This replica thinks it's a master — point it to the real master
			logger.Info("Replica is acting as master, correcting replication", "idx", i)
			p := redisPodName(bcredis.Name, i)
			_, _ = r.execPodCommand(ctx, bcredis.Namespace, p, "redis",
				redisCmd(spec, "REPLICAOF", masterSvc, "6379"))
		}
	}

	return false, nil
}

func (r *BcredisReconciler) isPodExecReady(
	ctx context.Context,
	namespace, podName string,
) (bool, error) {
	pod := &corev1.Pod{}
	if err := r.Get(
		ctx,
		types.NamespacedName{Namespace: namespace, Name: podName},
		pod,
	); err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	if pod.Spec.NodeName == "" || pod.Status.Phase != corev1.PodRunning {
		return false, nil
	}

	// Check if the redis container specifically is ready
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Name == "redis" {
			// Redis container must be in Running state to accept exec
			if cs.State.Running != nil {
				return true, nil
			}
			// If redis container is in Waiting/Starting state, wait longer
			return false, nil
		}
	}
	// Redis container not found in status - wait for it to start
	return false, nil
}

// findCurrentMaster scans all pods and returns the index of whichever is currently master (-1 if none).
// Prefers the master that has connected replicas to avoid picking a stale standalone master after pod restart.
func (r *BcredisReconciler) findCurrentMaster(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	logger logr.Logger,
) int {
	replicas := max(spec.Replicas, int32(3))
	bestIdx := -1
	bestSlaves := -1

	for i := 0; i < int(replicas); i++ {
		ready, err := r.isPodExecReady(ctx, bcredis.Namespace, redisPodName(bcredis.Name, i))
		if err != nil || !ready {
			continue
		}
		role, slaves, err := r.getRedisRoleAndSlaves(ctx, bcredis, spec, i, logger)
		if err != nil {
			continue
		}
		if role != masterStatus {
			continue
		}
		// Prefer the master with the most connected replicas
		if slaves > bestSlaves {
			bestSlaves = slaves
			bestIdx = i
		}
	}
	return bestIdx
}

// getRedisRoleAndSlaves returns role and number of connected_slaves from INFO replication
func (r *BcredisReconciler) getRedisRoleAndSlaves(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	idx int,
	logger logr.Logger,
) (string, int, error) {
	podName := redisPodName(bcredis.Name, idx)
	cmd := redisCmd(spec, "info", "replication")

	output, err := r.execPodCommand(ctx, bcredis.Namespace, podName, "redis", cmd)
	if err != nil {
		return "", 0, err
	}

	var role string
	slaves := 0
	for _, line := range splitLines(output) {
		if strings.HasPrefix(line, "role:") {
			role = strings.TrimSpace(strings.TrimPrefix(line, "role:"))
		}
		if strings.HasPrefix(line, "connected_slaves:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "connected_slaves:"))
			if n, err := strconv.Atoi(val); err == nil {
				slaves = n
			}
		}
	}
	if role == "" {
		return "", 0, fmt.Errorf("role not found in output")
	}
	return role, slaves, nil
}

// getRedisRole returns the Redis role (master or slave)
func (r *BcredisReconciler) getRedisRole(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	idx int,
	logger logr.Logger,
) (string, error) {
	role, _, err := r.getRedisRoleAndSlaves(ctx, bcredis, spec, idx, logger)
	return role, err
}

// performFailover promotes replica to master when master is unreachable
func (r *BcredisReconciler) performFailover(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	logger logr.Logger,
) error {
	// Cooldown guard — avoid rapid successive failovers
	if !bcredis.Status.LastFailoverTime.IsZero() {
		elapsed := time.Since(bcredis.Status.LastFailoverTime.Time)
		if elapsed < failoverCooldown {
			logger.Info(
				"Failover skipped — cooldown active",
				"elapsed", elapsed.Round(time.Second),
				"cooldown", failoverCooldown,
			)
			return nil
		}
	}

	logger.Info("Starting automatic failover")

	// Step 1: Promote replica (redis-1) to master
	currentMasterIdx := masterIndex(bcredis)
	candidateIdx, found, err := findReplicaCandidate(ctx, r, bcredis, spec, currentMasterIdx)
	if err != nil {
		return err
	}
	if !found {
		logger.Info("No replica candidate is exec-ready yet, will retry")
		return nil // soft requeue, not an error
	}

	candidatePod := redisPodName(bcredis.Name, candidateIdx)
	_, err = r.execPodCommand(
		ctx,
		bcredis.Namespace,
		candidatePod,
		"redis",
		redisCmd(spec, "SLAVEOF", "NO", "ONE"),
	)
	if err != nil {
		logger.Error(err, "failed to promote replica to master", "candidate", candidatePod)
		return err
	}

	time.Sleep(5 * time.Second)
	newMasterService := redisServiceName(bcredis.Name, candidateIdx)
	replicas := max(spec.Replicas, int32(3))
	for i := 0; i < int(replicas); i++ {
		if i == candidateIdx {
			continue
		}
		p := redisPodName(bcredis.Name, i)
		_, _ = r.execPodCommand(
			ctx,
			bcredis.Namespace,
			p,
			"redis",
			redisCmd(spec, "SLAVEOF", newMasterService, "6379"),
		)
	}

	// Step 3: Update status to indicate failover occurred
	bcredis.Status.MasterPod = candidatePod
	bcredis.Status.CurrentMasterService = newMasterService
	bcredis.Status.FailedOver = true
	bcredis.Status.LastFailoverTime = metav1.Now()

	if err := r.Status().Update(ctx, bcredis); err != nil {
		logger.Error(err, "failed to update status after failover")
		return err
	}

	return nil
}

// execPodCommand executes a command inside a pod container
func (r *BcredisReconciler) execPodCommand(
	ctx context.Context, namespace,
	podName, containerName string,
	command []string,
) (string, error) {
	config := r.RestConfig
	if config == nil {
		return "", fmt.Errorf("rest config is not initialized")
	}

	// Build URL with command as path segments
	req := r.KubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").VersionedParams(
		&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec,
	)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})
	if err != nil {
		return "", fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// splitLines splits a string by newlines and trims empty lines
func splitLines(s string) []string {
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
