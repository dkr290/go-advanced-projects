package controller

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/remotecommand"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	masterStatus  = "master"
	replicaStatus = "slave"
)

// reconcileRoles checks Redis replication status and triggers automatic failover
func (r *BcredisReconciler) reconcileRoles(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
) (bool, error) {
	logger := logf.FromContext(ctx).WithValues("bcredis", bcredis.Name)

	// Check if master (redis-0) is reachable
	masterReachable, err := r.checkRedisReachable(ctx, bcredis, 0, logger)
	if err != nil {
		logger.Error(err, "failed to check master reachability")
		return false, err
	}

	if !masterReachable {
		logger.Info("Master is unreachable, triggering automatic failover")
		return true, r.performFailover(ctx, bcredis, spec, logger)
	}

	// Check master role
	masterRole, err := r.getRedisRole(ctx, bcredis, 0, logger)
	if err != nil {
		logger.Error(err, "failed to get master role")
		return false, err
	}

	// Check replica (redis-1) role
	replicaRole, err := r.getRedisRole(ctx, bcredis, 1, logger)
	if err != nil {
		logger.Error(err, "failed to get replica role")
		return false, err
	}

	logger.Info("Redis roles", "master", masterRole, "replica", replicaRole)

	// If master is not master, trigger failover
	if masterRole != masterStatus {
		logger.Info("Master is not master, triggering automatic failover")
		return true, r.performFailover(ctx, bcredis, spec, logger)
	}

	// If replica is not replica, reconfigure it
	if replicaRole != replicaStatus {
		logger.Info("Replica is not replica, reconfiguring")
		return true, r.reconfigureReplica(ctx, bcredis, spec, logger)
	}

	return false, nil
}

// checkRedisReachable checks if Redis is reachable via ping
func (r *BcredisReconciler) checkRedisReachable(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	idx int,
	logger logr.Logger,
) (bool, error) {
	podName := fmt.Sprintf("%s-redis-%d", bcredis.Name, idx)
	cmd := []string{"redis-cli", "ping"}

	output, err := r.execPodCommand(ctx, podName, "redis", cmd)
	if err != nil {
		logger.Error(err, "failed to ping redis", "pod", podName)
		return false, err
	}

	return output == "PONG", nil
}

// getRedisRole returns the Redis role (master or slave)
func (r *BcredisReconciler) getRedisRole(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	idx int,
	logger logr.Logger,
) (string, error) {
	podName := fmt.Sprintf("%s-redis-%d", bcredis.Name, idx)
	cmd := []string{"redis-cli", "info", "replication"}

	output, err := r.execPodCommand(ctx, podName, "redis", cmd)
	if err != nil {
		logger.Error(err, "failed to get redis info", "pod", podName)
		return "", err
	}

	// Parse role from output
	for _, line := range splitLines(output) {
		if len(line) >= 6 && line[:6] == "role:" {
			return line[7:], nil
		}
	}
	return "", fmt.Errorf("role not found in output")
}

// performFailover promotes replica to master when master is unreachable
func (r *BcredisReconciler) performFailover(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	logger logr.Logger,
) error {
	logger.Info("Starting automatic failover")

	// Step 1: Promote replica (redis-1) to master
	replicaPod := fmt.Sprintf("%s-redis-1", bcredis.Name)
	promoteCmd := []string{"redis-cli", "SLAVEOF", "NO", "ONE"}

	_, err := r.execPodCommand(ctx, replicaPod, "redis", promoteCmd)
	if err != nil {
		logger.Error(err, "failed to promote replica to master")
		return err
	}
	logger.Info("Replica promoted to master")

	// Step 2: Wait for Sentinel to detect the change
	logger.Info("Waiting for Sentinel to detect failover...")
	time.Sleep(10 * time.Second)

	// Step 3: Update status to indicate failover occurred
	bcredis.Status.MasterPod = replicaPod
	bcredis.Status.ReplicaPod = fmt.Sprintf("%s-redis-0", bcredis.Name)
	bcredis.Status.FailedOver = true
	bcredis.Status.LastFailoverTime = metav1.Now()

	return nil
}

// reconfigureReplica reconfigures a non-replica pod to be a replica
func (r *BcredisReconciler) reconfigureReplica(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	logger logr.Logger,
) error {
	replicaPod := fmt.Sprintf("%s-redis-1", bcredis.Name)
	masterService := fmt.Sprintf("%s-redis-0.%s.svc.cluster.local", bcredis.Name, bcredis.Namespace)

	// Reconfigure replica to point to master
	slaveOfCmd := []string{"redis-cli", "SLAVEOF", masterService, "6379"}

	_, err := r.execPodCommand(ctx, replicaPod, "redis", slaveOfCmd)
	if err != nil {
		logger.Error(err, "failed to reconfigure replica")
		return err
	}
	logger.Info("Replica reconfigured to point to master")
	return nil
}

// execPodCommand executes a command inside a pod container
func (r *BcredisReconciler) execPodCommand(
	ctx context.Context,
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
		Namespace("default").
		SubResource("exec")

	// Add command as path parameters
	for _, cmd := range command {
		req = req.Suffix(cmd)
	}
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             nil,
		Stdout:            stdout,
		Stderr:            stderr,
		Tty:               false,
		TerminalSizeQueue: nil,
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
