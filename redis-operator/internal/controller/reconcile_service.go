package controller

import (
	"context"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// reconcileService creates/updates a ClusterIP Service for a Redis instance.
func (r *BcredisReconciler) reconcileService(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
	idx int,
) error {
	svcName := fmt.Sprintf("%s-redis-%d", bcredis.Name, idx)
	labels := map[string]string{
		"app":      bcredis.Name + "-redis",
		"instance": fmt.Sprintf("%d", idx),
		"bcredis":  bcredis.Name,
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: bcredis.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
		if err := controllerutil.SetControllerReference(bcredis, svc, r.Scheme); err != nil {
			return err
		}
		svc.Spec = corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "redis",
					Port:       int32(spec.ServicePort),
					TargetPort: intstr.FromInt(int(spec.ServicePort)),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "sentinel",
					Port:       26379,
					TargetPort: intstr.FromInt(26379),
					Protocol:   corev1.ProtocolTCP,
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		}
		return nil
	})
	return err
}
