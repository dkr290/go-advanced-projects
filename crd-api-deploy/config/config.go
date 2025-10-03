// Package config to set some static values
package config

import (
	"os"
	"strconv"

	"model-image-deployer/internal/models"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Port                  int
	Group                 string
	CrdVersion            string
	Namespace             string
	ImageRepo             string
	ImagePullPolicy       string
	ImagePullSecret       string
	IngressType           string
	EnvoyGateway          string
	EnvoyGatewayNamespace string
	Limits                models.ResourceList
	Requests              models.ResourceList
	HTTPGetStartup        models.HTTPGetAction
	FailureThreshold      int
	PeriodSeconds         int
	RunAsUser             int
	RunAsNonRoot          bool
	RunAsGroup            int
	FSGroup               int
	Affinity              map[string]any
	Tolerations           []map[string]any
	Kind                  string
	Env                   string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", 8000),
		Group:           getEnv("GROUP", "apps.api.test"),
		CrdVersion:      getEnv("CRD_VERSION", "v1alpha1"),
		Namespace:       getEnv("NAMESPACE", "catwalk"),
		ImageRepo:       getEnv("DOCKER_REPO", "aacatwalkregistry.azurecr.io"),
		ImagePullPolicy: getEnv("IMAGE_PULL_POLICY", "IfNotPresent"),
		Limits: models.ResourceList{
			CPU:              getEnv("CPU_LIMIT", "4000m"),
			Memory:           getEnv("MEMORY_LIMIT", "4Gi"),
			EphemeralStorage: getEnv("EPHEMERAL_STORAGE", "6Gi"),
		},
		Requests: models.ResourceList{
			CPU:    "512m",
			Memory: "512Mi",
		},
		HTTPGetStartup: models.HTTPGetAction{
			Path: getEnv("PROBE_PATH", "/"),
			Port: int32(getEnv("PROBE_PATH", 8000)),
		},
		IngressType:           getEnv("INGRESS_TYPE", "httproute"),
		EnvoyGateway:          getEnv("ENVOY_GATEWAY", "default-gateway"),
		EnvoyGatewayNamespace: getEnv("ENVOY_GATEWAY_NAMESPACE", "envoy-gateway-system"),
		FailureThreshold:      getEnv("FAILURE_THRESHOLD", 50),
		PeriodSeconds:         getEnv("FAILURE_PERIOD_SECONDS", 10),
		RunAsUser:             getEnv("RUN_AS_USER", 1000),
		RunAsNonRoot:          getEnv("RUN_AS_NON_ROOT", true),
		RunAsGroup:            3000,
		FSGroup:               2000,
		ImagePullSecret:       getEnv("IMAGEPULL_SECRET", "regcred"),
		Kind:                  getEnv("KIND", "Simpleapi"),
		Env:                   getEnv("ENV", "test"),
		Affinity: map[string]any{
			"nodeAffinity": map[string]any{
				"requiredDuringSchedulingIgnoredDuringExecution": map[string]any{
					"nodeSelectorTerms": []map[string]any{
						{
							"matchExpressions": []map[string]any{
								{
									"key":      "app",
									"operator": "In",
									"values":   []string{"catwalk"},
								},
							},
						},
					},
				},
			},
		},
		Tolerations: []map[string]any{
			{
				"key":      "purpose",
				"operator": "Equal",
				"value":    "catwalk",
				"effect":   "NoSchedule",
			},
		},
	}
}

type EnvType interface {
	string | int | bool
}

func getEnv[T EnvType](key string, defaultValue T) T {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case string:
		return any(rawValue).(T)

	case int:
		if i, err := strconv.Atoi(rawValue); err == nil {
			return any(i).(T)
		} else {
			log.Fatal().Err(err).Msg("Error converting the value in config")
		}
	case bool:
		if b, err := strconv.ParseBool(rawValue); err == nil {
			// Return the parsed bool if successful
			return any(b).(T)
		} else {
			log.Fatal().Err(err).Msg("Error converting the value in config")
		}

	}
	return defaultValue
}
