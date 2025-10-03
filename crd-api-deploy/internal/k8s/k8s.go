// Package k8s - responsible for all k8s operations
package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/yaml"

	"model-image-deployer/internal/apierror"
	"model-image-deployer/internal/models"
)

// K8sClientInterface defines the interface for k8s client operations.
type K8sClientInterface interface {
	ApplyCRD(ctx context.Context, crdYAML string, namespace, resource, group, version string) error
	DeleteCrd(
		ctx context.Context,
		name, resource, group, kind, version, namespace string,
	) (*models.DeleteCrdResponse, error)
	ListAllAPPs(
		ctx context.Context,
		resource, group, kind, version string,
	) (*models.ListAPIResponse, error)
	ResourceNameForKind(kind, group, version string) (string, error)
}

// Client wraps Kubernetes clients
type Client struct {
	clientset     *kubernetes.Clientset
	dynamicClient dynamic.Interface
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Client{
		clientset:     clientset,
		dynamicClient: dynamicClient,
	}, nil
}

// getKubernetesConfig gets the Kubernetes configuration
func getKubernetesConfig() (*rest.Config, error) {
	// Try in-cluster config first
	// before that the service account with exact cluster toles and bindings is needed
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fall back to kubeconfig
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	return config, nil
}

// ApplyCRD applies a CRD YAML to the cluster using server-side apply.
func (c *Client) ApplyCRD(
	ctx context.Context,
	crdYAML string,
	namespace, resource, group, version string,
) error {
	log.Info().
		Str("namespace", namespace).
		Str("resource", resource).
		Str("group", group).
		Str("version", version).
		Msg("Applying CRD with server-side apply")

	var objMap map[string]any
	if err := yaml.Unmarshal([]byte(crdYAML), &objMap); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal YAML")
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	obj := unstructured.Unstructured{Object: objMap}
	if namespace != "" {
		obj.SetNamespace(namespace)
	}

	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	patchOptions := metav1.PatchOptions{
		FieldManager: "model-image-deployer",
		Force:        func() *bool { b := true; return &b }(),
	}

	var dr dynamic.ResourceInterface
	if obj.GetNamespace() != "" {
		dr = c.dynamicClient.Resource(gvr).Namespace(obj.GetNamespace())
	} else {
		dr = c.dynamicClient.Resource(gvr)
	}

	_, err := dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, []byte(crdYAML), patchOptions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to apply CRD with server-side apply")
		if apierrors.IsConflict(err) {
			return apierror.ErrK8sConflict
		}
		return apierror.ErrK8sAPIFailure
	}

	log.Info().Str("name", obj.GetName()).Msg("CRD applied successfully using server-side apply")
	return nil
}

func (c *Client) DeleteCrd(
	ctx context.Context,
	name, resource, group, kind, version, namespace string,
) (*models.DeleteCrdResponse, error) {
	log.Info().Str("name", name).Str("namespace", namespace).Msg("Deleting Crd")
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	err := c.dynamicClient.Resource(gvr).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete Crd")
		if apierrors.IsNotFound(err) {
			return nil, apierror.ErrNotFound
		}
		return nil, apierror.ErrK8sAPIFailure
	}

	response := &models.DeleteCrdResponse{
		Name:       name,
		APIVersion: group + "/" + version,
		Kind:       kind,
	}
	log.Info().Interface("response", response).Msg("Crd deleted successfully")
	return response, nil
}

// ListAllAPPs lists all SimpleAPI resources in a namespace
func (c *Client) ListAllAPPs(
	ctx context.Context,
	resource, group, kind, version string,
) (*models.ListAPIResponse, error) {
	log.Info().
		Str("resource", resource).
		Str("group", group).
		Str("kind", kind).
		Str("version", version).
		Msg("Listing all APPs")
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}

	list, err := c.dynamicClient.Resource(gvr).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to list Crds")
		return nil, fmt.Errorf("failed to list Crds in namespace  %w", err)
	}

	response := &models.ListAPIResponse{
		APIVersion: group + "/" + version,
		Kind:       kind,
		Items:      make([]models.GetAPIResponse, len(list.Items)),
	}

	for i, item := range list.Items {
		response.Items[i] = models.GetAPIResponse{
			APIVersion: item.GetAPIVersion(),
			Kind:       item.GetKind(),
			Metadata:   item.Object["metadata"].(map[string]any),
			Spec:       item.Object["spec"].(map[string]any),
		}
	}
	log.Info().Int("count", len(response.Items)).Msg("APPs listed successfully")
	return response, nil
}

// ResourceNameForKind return the plural resource name for a given Kind, Group, and Version
func (c *Client) ResourceNameForKind(kind, group, version string) (string, error) {
	log.Info().
		Str("kind", kind).
		Str("group", group).
		Str("version", version).
		Msg("Discovering resource name for kind")
	gv := group + "/" + version
	apiResourceList, err := c.clientset.Discovery().ServerResourcesForGroupVersion(gv)
	if err != nil {
		log.Error().Err(err).Msg("Failed to discover resources")
		return "", fmt.Errorf("failed to discover resources for %s: %w", gv, err)
	}
	for _, apiResource := range apiResourceList.APIResources {
		if apiResource.Kind == kind {
			log.Info().Str("resourceName", apiResource.Name).Msg("Resource name discovered")
			return apiResource.Name, nil // Plural resource name
		}
	}
	log.Error().Msg("Resource name not found")
	return "", fmt.Errorf("resource name for kind %s not found in %s", kind, gv)
}
