package controller

import (
	"context"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var masterConf = `bind 0.0.0.0
protected-mode no
appendonly yes
daemonize no
port 6379
`

var replicaConf = `bind 0.0.0.0
protected-mode no
appendonly yes
daemonize no
port 6379
`
var sentinelConf = `bind 0.0.0.0
protected-mode no
daemonize no
port 26379
sentinel monitor mymaster %s 6379 2
sentinel down-after-milliseconds mymaster 5000
sentinel failover-timeout mymaster 60000
sentinel parallel-syncs mymaster 1
`


// reconcileConfigMap creates/updates the Redis configuration ConfigMap.
func (r *BcredisReconciler) reconcileConfigMap(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcredis.Name + "-redis-config",
			Namespace: bcredis.Namespace,
			Labels: map[string]string{
				"app": bcredis.Name,
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		if err := controllerutil.SetControllerReference(bcredis, cm, r.Scheme); err != nil {
			return err
		}
		cm.Data = map[string]string{
			"master.conf":  masterConf,
			"replica.conf": replicaConf,
		}
		return nil
	})
	if err != nil {
		return err
	}
// reconcile sentinel ConfigMap with master hostname placeholder 
sentinelCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcredis.Name + "-sentinel-config",
			Namespace: bcredis.Namespace,
			Labels: map[string]string{
				"app": bcredis.Name,
			},
		},
	}
_, err = controllerutil.CreateOrUpdate(ctx, r.Client, sentinelCM, func() error {
		if err := controllerutil.SetControllerReference(bcredis, sentinelCM, r.Scheme); err != nil {
			return err
		}
	// Use master-0 as the monitored master hostname
		masterHostname := fmt.Sprintf("%s-redis-0.%s.svc.cluster.local", bcredis.Name, bcredis.Namespace)
		sentinelCM.Data = map[string]string{
			"sentinel.conf": fmt.Sprintf(sentinelConf, masterHostname),
		}
		return nil
	})
	return err


}
