package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"
	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// reconcileHeadlessService creates a headless service for Sentinel pod discovery
func (r *BcredisReconciler) reconcileHeadlessService(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
) error {
	svcName := fmt.Sprintf("%s-redis-headless", bcredis.Name)
	labels := map[string]string{
		"app":     bcredis.Name + "-redis",
		"bcredis": bcredis.Name,
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
			ClusterIP: corev1.ClusterIPNone, // Headless service
			Selector:  labels,
			PublishNotReadyAddresses: true,                 // optional, helps early DNS for pod discovery
			Ports: []corev1.ServicePort{
				{
					Name:       "redis",
					Port:       6379,
					TargetPort: intstr.FromInt(6379),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "sentinel",
					Port:       26379,
					TargetPort: intstr.FromInt(26379),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		}
		return nil
	})
	return err
}

