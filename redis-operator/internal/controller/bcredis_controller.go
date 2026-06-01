// Package controller
package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	bcredisv1alpha1 "github.com/example/redis-operator/api/v1alpha1"
)

const (
	finalizerName   = "bcredis.bankingcircle.net/finalizer"
	instanceZero    = "redis-0"
	instanceOne     = "redis-1"
	requeAfter      = 15 * time.Second
	heathCheckRetry = 3 * time.Second
)

// BcredisReconciler reconciles a Bcredis object

type BcredisReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	RestConfig *rest.Config
	KubeClient *kubernetes.Clientset
}

// +kubebuilder:rbac:groups=bcredis.bankingcircle.net,resources=bcredis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bcredis.bankingcircle.net,resources=bcredis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=bcredis.bankingcircle.net,resources=bcredis/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services;pods;persistentvolumeclaims;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=pods/exec,verbs=create
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways;tcproutes,verbs=get;list;watch;create;update;patch;delete
// BcredisReconciler reconciles a Bcredis object

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Bcredis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *BcredisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx).WithValues("bcredis", req.NamespacedName)
	ctx = logf.IntoContext(ctx, logger)
	logger.Info("Reconciling bcredis")

	// Fetch the Bcredis instance
	bcredis := &bcredisv1alpha1.Bcredis{}
	if err := r.Get(ctx, req.NamespacedName, bcredis); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("bcredis resource not found. Ignoring sice object must be deleted")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get bcredis")
		return ctrl.Result{}, err
	}
	logger = logger.WithValues(
		"name",
		bcredis.Name,
		"namespace",
		bcredis.Namespace,
		"image",
		bcredis.Spec.RedisImage,
	)
	// Delete all owned resources before removing the finalizer
if bcredis.DeletionTimestamp != nil {
    if controllerutil.ContainsFinalizer(bcredis, finalizerName) {
        // Delete all owned resources
        r.deleteOwnedResources(ctx, bcredis)
        controllerutil.RemoveFinalizer(bcredis, finalizerName)
        if err := r.Update(ctx, bcredis); err != nil {
            return ctrl.Result{}, err
        }
    }
}
	// Apply defaults
	spec := bcredis.Spec
	if spec.RedisImage == "" {
		spec.RedisImage = "redis:7.2"
	}
	if spec.StorageClassName == "" {
		spec.StorageClassName = "azurefile-csi"
	}
	if spec.StorageSize.IsZero() {
		spec.StorageSize = resource.MustParse("1Gi")
	}
	if spec.ServicePort == 0 {
		spec.ServicePort = 6379
	}

	name := bcredis.Name
	//	ns := bcredis.Namespace

	// Reconcile ConfigMap for Redis configuration

	currentMasterSvc := bcredis.Status.CurrentMasterService
	if currentMasterSvc == "" {
		currentMasterSvc = fmt.Sprintf("%s-redis-0", bcredis.Name)
	}
	if err := r.reconcileConfigMap(ctx, bcredis, currentMasterSvc); err != nil {
		logger.Error(err, "failed to reconcile ConfigMap")
		return ctrl.Result{}, err
	}
// Reconcile headless service for Sentinel discovery
	if err := r.reconcileHeadlessService(ctx, bcredis, spec); err != nil {
		logger.Error(err, "failed to reconcile headless service")
		return ctrl.Result{}, err
	}

	// Reconcile statefullsets for both instances
	for _, idx := range []int{0, 1} {

		if err := r.reconcileStatefulSet(ctx, bcredis, spec, idx, logger); err != nil {
			logger.Error(err, "failed to reconcile StatefulSet", "index", idx)
			return ctrl.Result{}, err

		}
		if err := r.reconcileService(ctx, bcredis, spec, idx); err != nil {
			logger.Error(err, "failed to reconcile Service", "index", idx)
			return ctrl.Result{}, err
		}

	}
	// Reconcile Envoy Gateway and TCPRoute
	if err := r.reconcileGateway(ctx, bcredis, spec); err != nil {
		logger.Error(err, "failed to reconcile Gateway")
		return ctrl.Result{}, err
	}
	// Reconcile PDB
	pdb := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bcredis.Name + "-redis-pdb",
			Namespace: bcredis.Namespace,
		},
	}
	_, pdbErr := controllerutil.CreateOrUpdate(ctx, r.Client, pdb, func() error {
		if err := controllerutil.SetControllerReference(bcredis, pdb, r.Scheme); err != nil {
			return err
		}
		pdb.Spec.MinAvailable = &intstr.IntOrString{IntVal: 1}
		pdb.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": bcredis.Name + "-redis",
			},
		}
		return nil
	})
	if pdbErr != nil {
		logger.Error(pdbErr, "failed to reconcile PDB")
		return ctrl.Result{}, pdbErr
	}

	// Determine current master/replica state and perform failover if needed
	requeue, err := r.reconcileRoles(ctx, bcredis, spec)
	if err != nil {
		logger.Error(err, "failed to reconcile roles")
		return ctrl.Result{RequeueAfter: heathCheckRetry}, err
	}

	// Update status
	if err := r.Status().Update(ctx, bcredis); err != nil {
		if errors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}



	_ = name
	if requeue {
		return ctrl.Result{RequeueAfter: requeAfter}, nil
	}
	return ctrl.Result{RequeueAfter: requeAfter}, nil


}

// SetupWithManager sets up the controller with the Manager.
func (r *BcredisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bcredisv1alpha1.Bcredis{}).
		Named("bcredis").
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}


