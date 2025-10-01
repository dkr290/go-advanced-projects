// Package k8s deploying the crd after
package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	"crd-api-deploy/internal/models"

	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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

// ApplyCRD applies a CRD YAML to the cluster
func (c *Client) ApplyCRD(ctx context.Context, crdYAML string, namespace string) error {
	var obj unstructured.Unstructured
	err := yaml.Unmarshal([]byte(crdYAML), &obj.Object)
	if err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if namespace != "" {
		obj.SetNamespace(namespace)
	}

	gvr := schema.GroupVersionResource{
		Group:    "apps.api.test",
		Version:  "v1alpha1",
		Resource: "simpleapis",
	}

	if namespace != "" && namespace != "default" {
		err = c.ensureNamespace(ctx, namespace)
		if err != nil {
			return fmt.Errorf("failed to ensure namespace exists: %w", err)
		}
	}

	if obj.GetNamespace() != "" {
		_, err = c.dynamicClient.Resource(gvr).
			Namespace(obj.GetNamespace()).
			Create(ctx, &obj, metav1.CreateOptions{})
		if err != nil {
			_, err = c.dynamicClient.Resource(gvr).
				Namespace(obj.GetNamespace()).
				Update(ctx, &obj, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create/update resource: %w", err)
			}
		}
	} else {
		_, err = c.dynamicClient.Resource(gvr).Create(ctx, &obj, metav1.CreateOptions{})
		if err != nil {
			_, err = c.dynamicClient.Resource(gvr).Update(ctx, &obj, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create/update resource: %w", err)
			}
		}
	}

	return nil
}

// ensureNamespace ensures that a namespace exists
func (c *Client) ensureNamespace(ctx context.Context, namespace string) error {
	_, err := c.clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		_, err = c.clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create namespace %s: %w", namespace, err)
		}
	}
	return nil
}

// GetSimpleAPI retrieves a SimpleAPI resource
func (c *Client) GetSimpleAPI(
	ctx context.Context,
	name, namespace string,
) (*models.GetAPIResponse, error) {
	gvr := schema.GroupVersionResource{
		Group:    "apps.api.test",
		Version:  "v1alpha1",
		Resource: "simpleapis",
	}

	obj, err := c.dynamicClient.Resource(gvr).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get SimpleAPI %s/%s: %w", namespace, name, err)
	}

	response := &models.GetAPIResponse{
		APIVersion: obj.GetAPIVersion(),
		Kind:       obj.GetKind(),
		Metadata:   obj.Object["metadata"].(map[string]any),
		Spec:       obj.Object["spec"].(map[string]any),
	}

	return response, nil
}

// ListSimpleAPIs lists all SimpleAPI resources in a namespace
func (c *Client) ListSimpleAPIs(
	ctx context.Context,
	namespace string,
) (*models.ListSimpleAPIResponse, error) {
	gvr := schema.GroupVersionResource{
		Group:    "apps.api.test",
		Version:  "v1alpha1",
		Resource: "simpleapis",
	}

	list, err := c.dynamicClient.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list SimpleAPIs in namespace %s: %w", namespace, err)
	}

	response := &models.ListSimpleAPIResponse{
		APIVersion: "apps.api.test/v1alpha1",
		Kind:       "SimpleAPIList",
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

	return response, nil
}
