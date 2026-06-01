package controller

import (
	"bytes"
	"context"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

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
func (r *BcredisReconciler) deleteOwnedResources(ctx context.Context, bcredis *bcredisv1alpha1.Bcredis) {
    // Delete StatefulSets
    for _, idx := range []int{0, 1} {
        ss := &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s-redis-%d", bcredis.Name, idx), Namespace: bcredis.Namespace}}
        r.Delete(ctx, ss)
    }
    // Delete Services
    for _, idx := range []int{0, 1} {
        svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s-redis-%d", bcredis.Name, idx), Namespace: bcredis.Namespace}}
        r.Delete(ctx, svc)
    }
    // Delete headless service
    headless := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("%s-redis-headless", bcredis.Name), Namespace: bcredis.Namespace}}
    r.Delete(ctx, headless)
    // Delete ConfigMaps
    cm1 := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: bcredis.Name + "-redis-config", Namespace: bcredis.Namespace}}
    r.Delete(ctx, cm1)
    cm2 := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: bcredis.Name + "-sentinel-config", Namespace: bcredis.Namespace}}
    r.Delete(ctx, cm2)
    // Delete Gateway and TCPRoute
    gw := &gatewayv1.Gateway{ObjectMeta: metav1.ObjectMeta{Name: bcredis.Name + "-redis-gateway", Namespace: bcredis.Namespace}}
    r.Delete(ctx, gw)
		tr := &gatewayv1alpha2.TCPRoute{
      ObjectMeta: metav1.ObjectMeta{Name: bcredis.Name+"-tcproute",Namespace: bcredis.Namespace},
		}
    r.Delete(ctx, tr)
}
