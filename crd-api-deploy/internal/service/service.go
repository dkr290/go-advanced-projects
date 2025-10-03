// Package service
package service

import (
	"context"
	"fmt"
	"strings"

	"model-image-deployer/config"
	"model-image-deployer/internal/apierror"
	"model-image-deployer/internal/k8s"
	"model-image-deployer/internal/models"
	"model-image-deployer/internal/template"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/util/validation"
)

// APIServiceInterface defines the interface for the APIService.
type APIServiceInterface interface {
	CreateAPP(ctx context.Context, req *models.CreateCrdRequest) (*models.CreateCrdResponse, error)
	ListAPPs(ctx context.Context) (*models.ListAPIResponse, error)
	DeleteAPP(ctx context.Context, req *models.DeleteCrdInput) (*models.DeleteCrdResponse, error)
}

// APIService handles SimpleAPI operations
type APIService struct {
	config         *config.Config
	k8sClient      k8s.K8sClientInterface
	templateEngine template.TemplateEngineInterface
}

// NewAPIService creates a new SimpleAPI service
func NewAPIService(
	cfg *config.Config,
	k8sClient k8s.K8sClientInterface,
	templateEngine template.TemplateEngineInterface,
) (*APIService, error) {
	return &APIService{
		config:         cfg,
		k8sClient:      k8sClient,
		templateEngine: templateEngine,
	}, nil
}

// CreateAPP creates  CRD in the cluster
func (s *APIService) CreateAPP(
	ctx context.Context, req *models.CreateCrdRequest,
) (*models.CreateCrdResponse, error) {
	log.Info().Interface("request", req).Msg("Creating CRD")
	if err := s.validateCreateRequest(req); err != nil {
		log.Error().Err(err).Msg("Validation failed")
		return nil, err
	}

	defaultValues := s.setDefaultValues(req)
	log.Debug().Interface("defaults", defaultValues).Msg("Default values set")

	crdYAML, err := s.templateEngine.GenerateCRD(defaultValues)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate CRD YAML")
		return nil, fmt.Errorf("failed to generate CRD YAML: %w", err)
	}
	log.Debug().Str("yaml", crdYAML).Msg("Generated CRD YAML")

	resourceName, err := s.k8sClient.ResourceNameForKind(
		defaultValues.Kind,
		defaultValues.Group,
		defaultValues.CrdVersion,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get resource from kind")
		return nil, fmt.Errorf("failed to get Resource from kind %v", err)
	}
	log.Debug().Str("resourceName", resourceName).Msg("Got resource name")

	err = s.k8sClient.ApplyCRD(
		ctx,
		crdYAML,
		defaultValues.Namespace,
		resourceName,
		defaultValues.Group,
		defaultValues.CrdVersion,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to apply CRD to cluster")
		return nil, err
	}

	response := &models.CreateCrdResponse{
		Message:   "CRD resource created successfully",
		Name:      defaultValues.Name,
		Namespace: defaultValues.Namespace,
		Kind:      defaultValues.Kind,
	}
	log.Info().Interface("response", response).Msg("CRD created successfully")
	return response, nil
}

// ListAPPs lists SimpleAPP resources in a namespace
func (s *APIService) ListAPPs(
	ctx context.Context,
) (*models.ListAPIResponse, error) {
	log.Info().Msg("Listing CRDs")
	conf := s.config

	resourceName, err := s.k8sClient.ResourceNameForKind(conf.Kind, conf.Group, conf.CrdVersion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get resource from kind")
		return nil, fmt.Errorf("failed to get Resource from kind %v", err)
	}
	log.Debug().Str("resourceName", resourceName).Msg("Got resource name")

	response, err := s.k8sClient.ListAllAPPs(
		ctx,
		resourceName,
		conf.Group,
		conf.Kind,
		conf.CrdVersion,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list CRDs")
		return nil, err
	}
	log.Info().Int("count", len(response.Items)).Msg("CRDs listed successfully")
	return response, nil
}

func (s *APIService) DeleteAPP(
	ctx context.Context, req *models.DeleteCrdInput,
) (*models.DeleteCrdResponse, error) {
	log.Info().Interface("request", req).Msg("Deleting CRD")
	conf := s.config

	resourceName, err := s.k8sClient.ResourceNameForKind(conf.Kind, conf.Group, conf.CrdVersion)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get resource from kind")
		return nil, fmt.Errorf("failed to get Resource from kind %v", err)
	}
	log.Debug().Str("resourceName", resourceName).Msg("Got resource name")

	response, err := s.k8sClient.DeleteCrd(
		ctx,
		req.Name,
		resourceName,
		conf.Group,
		conf.Kind,
		conf.CrdVersion,
		conf.Namespace,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete CRD")
		return nil, err
	}
	log.Info().Interface("response", response).Msg("CRD deleted successfully")
	return response, nil
}

// validateCreateRequest validates the create request
func (s *APIService) validateCreateRequest(req *models.CreateCrdRequest) error {
	if req.Name == "" {
		return fmt.Errorf("%w: name is required", apierror.ErrInvalidInput)
	}

	if errs := validation.IsDNS1123Subdomain(req.Name); len(errs) > 0 {
		return fmt.Errorf("%w: invalid name: %s", apierror.ErrInvalidInput, strings.Join(errs, ", "))
	}

	if req.Version == "" {
		return fmt.Errorf("%w: version is required", apierror.ErrInvalidInput)
	}

	return nil
}

// setDefaultValues sets default values for optional fields
func (s *APIService) setDefaultValues(req *models.CreateCrdRequest) *models.SetDefaultValues {
	conf := s.config

	if req.Replicas == nil {
		defaultReplicas := 1
		req.Replicas = &defaultReplicas

	}

	httpHostName := strings.ToLower(fmt.Sprintf("%s.k8s-%s.bankingcircle.net", req.Name, conf.Env))

	defaults := &models.SetDefaultValues{
		Labels: []models.Label{
			{Key: "app.kubernetes.io/name", Value: req.Name},
			{Key: "app.kubernetes.io/managed-by", Value: "operator"},
			{Key: "app", Value: req.Name},
		},
		Group:                 conf.Group,
		CrdVersion:            conf.CrdVersion,
		Kind:                  conf.Kind,
		Env:                   conf.Env,
		Namespace:             conf.Namespace,
		Name:                  strings.ToLower(req.Name),
		Image:                 conf.ImageRepo + "/" + req.Name,
		ServiceAccount:        strings.ToLower(req.Name + "-sa"),
		ImagePullPolicy:       conf.ImagePullPolicy,
		EnvoyGateway:          conf.EnvoyGateway,
		IngressType:           conf.IngressType,
		EnvoyGatewayNamespace: conf.EnvoyGatewayNamespace,
		Version:               req.Version,
		Replicas:              req.Replicas,
		Port:                  int32(conf.Port),
		Resources: models.ResourceRequirements{
			Limits:   conf.Limits,
			Requests: conf.Requests,
		},
		StartupProbe: models.StartupProbe{
			HTTPGet:          conf.HTTPGetStartup,
			FailureThreshold: int32(conf.FailureThreshold),
			PeriodSeconds:    int32(conf.PeriodSeconds),
		},
		PodSecurityContext: models.PodSecurityContext{
			RunAsUser:    int64(conf.RunAsUser),
			RunAsNonRoot: conf.RunAsNonRoot,
			RunAsGroup:   int64(conf.RunAsGroup),
			FSGroup:      int64(conf.FSGroup),
		},
	}
	if conf.Env == "prd" {
		defaults.Affinity = conf.Affinity
		defaults.Tolerations = conf.Tolerations
		defaults.IngressHostName = httpHostName
	}
	if conf.Env == "uat" || conf.Env == "dev" {
		defaults.IngressHostName = httpHostName
	}

	if conf.Env == "test" {
		defaults.ImagePullSecret = conf.ImagePullSecret
	}
	return defaults
}
