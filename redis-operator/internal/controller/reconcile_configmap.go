package controller

import (
	"context"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var masterConf = `bind 0.0.0.0
protected-mode no
appendonly yes
`

var replicaConf = `bind 0.0.0.0
protected-mode no
appendonly yes
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
	return err
}
