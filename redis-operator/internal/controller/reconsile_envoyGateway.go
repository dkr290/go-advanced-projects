package controller

import (
	"context"
	"fmt"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// reconcileGateway creates/updates the Envoy Gateway and TCPRoute.
func (r *BcredisReconciler) reconcileGateway(
	ctx context.Context,
	bcredis *bcredisv1alpha1.Bcredis,
	spec bcredisv1alpha1.BcredisSpec,
) error {
	gwName := bcredis.Name + "-gateway"
	gwClass := gatewayv1.ObjectName(spec.EnvoyGatewayClassName)
	ns := bcredis.Namespace
	port := gatewayv1.PortNumber(spec.ServicePort)

	gw := &gatewayv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gwName,
			Namespace: ns,
		},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, gw, func() error {
		if err := controllerutil.SetControllerReference(bcredis, gw, r.Scheme); err != nil {
			return err
		}

		var infra *gatewayv1.GatewayInfrastructure
		if len(spec.GatewayAnnotations) > 0 {
			infraAnnotations := make(
				map[gatewayv1.AnnotationKey]gatewayv1.AnnotationValue,
				len(spec.GatewayAnnotations),
			)
			for k, v := range spec.GatewayAnnotations {
				infraAnnotations[gatewayv1.AnnotationKey(k)] = gatewayv1.AnnotationValue(v)
			}
			infra = &gatewayv1.GatewayInfrastructure{
				Annotations: infraAnnotations,
			}

		}


		gw.Spec = gatewayv1.GatewaySpec{
			GatewayClassName: gwClass,
			Infrastructure:   infra,
			Listeners: []gatewayv1.Listener{
				{
					Name:     "redis-tcp",
					Protocol: gatewayv1.TCPProtocolType,
					Port:     port,
				},
			},
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Determine master service name from status
	masterSvc := bcredis.Status.CurrentMasterService
	if masterSvc == "" {
		// Default to instance 0 on first creation
		masterSvc = fmt.Sprintf("%s-redis-0", bcredis.Name)
	}

	nsPtr := gatewayv1.Namespace(ns)
	portPtr := gatewayv1.PortNumber(6379)
	gwNs := gatewayv1.Namespace(ns)

	tcpRoute := &gatewayv1alpha2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcredis.Name + "-tcproute",
			Namespace: ns,
		},
	}
	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, tcpRoute, func() error {
		if err := controllerutil.SetControllerReference(bcredis, tcpRoute, r.Scheme); err != nil {
			return err
		}
		tcpRoute.Spec = gatewayv1alpha2.TCPRouteSpec{
			CommonRouteSpec: gatewayv1.CommonRouteSpec{
				ParentRefs: []gatewayv1.ParentReference{
					{
						Name:        gatewayv1.ObjectName(gwName),
						Namespace:   &gwNs,
						SectionName: sectionNamePtr("redis-tcp"),
					},
				},
			},
			Rules: []gatewayv1alpha2.TCPRouteRule{
				{
					BackendRefs: []gatewayv1.BackendRef{
						{
							BackendObjectReference: gatewayv1.BackendObjectReference{
								Name:      gatewayv1.ObjectName(masterSvc),
								Port:      &portPtr,
								Namespace: &nsPtr,
							},
						},
					},
				},
			},
		}
		return nil
	})
	return err
}
