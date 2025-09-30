package service

import (
	"context"
	"fmt"

	"crd-api-deploy/internal/k8s"
	"crd-api-deploy/internal/models"
	"crd-api-deploy/internal/template"
)

// SimpleAPIService handles SimpleAPI operations
type SimpleAPIService struct {
	k8sClient      *k8s.Client
	templateEngine *template.Engine
}

// NewSimpleAPIService creates a new SimpleAPI service
func NewSimpleAPIService() (*SimpleAPIService, error) {
	k8sClient, err := k8s.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	templateEngine, err := template.NewEngine()
	if err != nil {
		return nil, fmt.Errorf("failed to create template engine: %w", err)
	}

	return &SimpleAPIService{
		k8sClient:      k8sClient,
		templateEngine: templateEngine,
	}, nil
}

// CreateSimpleAPI creates a SimpleAPI CRD in the cluster
func (s *SimpleAPIService) CreateSimpleAPI(ctx context.Context, req *models.CreateSimpleAPIRequest) (*models.CreateSimpleAPIResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	s.setDefaultValues(req)

	crdYAML, err := s.templateEngine.GenerateCRD(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CRD YAML: %w", err)
	}

	err = s.k8sClient.ApplyCRD(ctx, crdYAML, req.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to apply CRD to cluster: %w", err)
	}

	return &models.CreateSimpleAPIResponse{
		Message:   "SimpleAPI resource created successfully",
		Name:      req.Name,
		Namespace: req.Namespace,
		Kind:      req.Kind,
	}, nil
}

// GetSimpleAPI retrieves a SimpleAPI resource
func (s *SimpleAPIService) GetSimpleAPI(ctx context.Context, name, namespace string) (*models.GetSimpleAPIResponse, error) {
	return s.k8sClient.GetSimpleAPI(ctx, name, namespace)
}

// ListSimpleAPIs lists SimpleAPI resources in a namespace
func (s *SimpleAPIService) ListSimpleAPIs(ctx context.Context, namespace string) (*models.ListSimpleAPIResponse, error) {
	return s.k8sClient.ListSimpleAPIs(ctx, namespace)
}

// validateCreateRequest validates the create request
func (s *SimpleAPIService) validateCreateRequest(req *models.CreateSimpleAPIRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if req.Kind == "" {
		return fmt.Errorf("kind is required")
	}
	if req.Image == "" {
		return fmt.Errorf("image is required")
	}
	if req.Version == "" {
		return fmt.Errorf("version is required")
	}

	return nil
}

// setDefaultValues sets default values for optional fields
func (s *SimpleAPIService) setDefaultValues(req *models.CreateSimpleAPIRequest) {
	if len(req.Labels) == 0 {
		req.Labels = []models.Label{
			{Key: "app.kubernetes.io/name", Value: "my-api"},
			{Key: "app.kubernetes.io/managed-by", Value: "operator"},
			{Key: "app", Value: "my-new-api"},
		}
	}

	if req.Port == 0 {
		req.Port = 8000
	}
	if req.Replicas == 0 {
		req.Replicas = 1
	}
	if req.IngressType == "" {
		req.IngressType = "ingress"
	}
	if req.Resources.Limits.CPU == "" {
		req.Resources.Limits.CPU = "1000m"
	}
	if req.Resources.Limits.Memory == "" {
		req.Resources.Limits.Memory = "2Gi"
	}
	if req.Resources.Limits.EphemeralStorage == "" {
		req.Resources.Limits.EphemeralStorage = "6Gi"
	}
	if req.Resources.Requests.CPU == "" {
		req.Resources.Requests.CPU = "500m"
	}
	if req.Resources.Requests.Memory == "" {
		req.Resources.Requests.Memory = "256Mi"
	}
	if req.StartupProbe.HTTPGet.Path == "" {
		req.StartupProbe.HTTPGet.Path = "/"
	}
	if req.StartupProbe.HTTPGet.Port == 0 {
		req.StartupProbe.HTTPGet.Port = req.Port
	}
	if req.StartupProbe.FailureThreshold == 0 {
		req.StartupProbe.FailureThreshold = 20
	}
	if req.StartupProbe.PeriodSeconds == 0 {
		req.StartupProbe.PeriodSeconds = 10
	}
	if req.PodSecurityContext.RunAsUser == 0 {
		req.PodSecurityContext.RunAsNonRoot = true
		req.PodSecurityContext.RunAsUser = 1000
		req.PodSecurityContext.RunAsGroup = 3000
		req.PodSecurityContext.FSGroup = 2000
	}
	if req.ServiceAccount == "" {
		req.ServiceAccount = "simpleapi-sa"
	}
	if req.ImagePullSecret == "" {
		req.ImagePullSecret = "regcred"
	}
	if req.Affinity == nil {
		req.Affinity = make(map[string]interface{})
	}
	if req.Tolerations == nil {
		req.Tolerations = []map[string]interface{}{}
	}
}