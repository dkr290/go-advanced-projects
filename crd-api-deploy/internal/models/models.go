// Package models - all huma requests and responces
package models

// Label represents a key-value pair for Kubernetes labels
type Label struct {
	Key   string `json:"key"   example:"app"    doc:"Label key"`
	Value string `json:"value" example:"my-app" doc:"Label value"`
}

// ResourceRequirements defines resource limits and requests
type ResourceRequirements struct {
	Limits   ResourceList `json:"limits"   doc:"Resource limits"`
	Requests ResourceList `json:"requests" doc:"Resource requests"`
}

// ResourceList defines resource quantities
type ResourceList struct {
	CPU              string `json:"cpu"                         example:"1000m" doc:"CPU quantity"`
	Memory           string `json:"memory"                      example:"2Gi"   doc:"Memory quantity"`
	EphemeralStorage string `json:"ephemeral-storage,omitempty" example:"6Gi"   doc:"Ephemeral storage quantity"`
}

// StartupProbe defines startup probe configuration
type StartupProbe struct {
	HTTPGet          HTTPGetAction `json:"httpGet"          doc:"HTTP GET action for probe"`
	FailureThreshold int32         `json:"failureThreshold" doc:"Failure threshold"         example:"20"`
	PeriodSeconds    int32         `json:"periodSeconds"    doc:"Period in seconds"         example:"10"`
}

// HTTPGetAction defines HTTP GET action for probes
type HTTPGetAction struct {
	Path string `json:"path" example:"/"    doc:"HTTP path"`
	Port int32  `json:"port" example:"8000" doc:"Port number"`
}

// PodSecurityContext defines security context for pod
type PodSecurityContext struct {
	RunAsNonRoot bool  `json:"runAsNonRoot" example:"true" doc:"Run as non-root user"`
	RunAsUser    int64 `json:"runAsUser"    example:"1000" doc:"User ID to run as"`
	RunAsGroup   int64 `json:"runAsGroup"   example:"3000" doc:"Group ID to run as"`
	FSGroup      int64 `json:"fsGroup"      example:"2000" doc:"Filesystem group ID"`
}

// CreateAPIRequest represents the request payload for creating a SimpleAPI CRD
type CreateAPIRequest struct {
	Kind               string               `json:"kind"                         example:"Simpleapi"             doc:"Kind of the CRD"`
	Group              string               `json:"group"                        example:"apps.api.test"         doc:"The Group name"`
	CrdVersion         string               `json:"crdversion,omitempty"         example:"v1alpha1"              doc:"crd version"`
	Port               int32                `json:"port,omitempty"               example:"8000"                  doc:"Port number"`
	Name               string               `json:"name"                         example:"simpleapi-new"         doc:"Name of the SimpleAPI resource"`
	Namespace          string               `json:"namespace"                    example:"default"               doc:"Namespace for the resource"`
	Labels             []Label              `json:"labels,omitempty"                                             doc:"Labels to be applied to the resource"`
	Replicas           int32                `json:"replicas,omitempty"           example:"1"                     doc:"Number of replicas"`
	Image              string               `json:"image"                        example:"nginx"                 doc:"Container image"`
	Version            string               `json:"version"                      example:"latest"                doc:"Image version/tag"`
	IngressHostName    string               `json:"ingressHostName,omitempty"    example:"simpleapi.example.com" doc:"Ingress hostname"`
	ImagePullSecret    string               `json:"imagePullSecret,omitempty"    example:"regcred"               doc:"Image pull secret name"`
	ServiceAccount     string               `json:"serviceAccount,omitempty"     example:"simpleapi-sa"          doc:"Service account name"`
	Resources          ResourceRequirements `json:"resources,omitempty"                                          doc:"Resource requirements"`
	PodSecurityContext PodSecurityContext   `json:"podSecurityContext,omitempty"                                 doc:"Pod security context"`
	StartupProbe       StartupProbe         `json:"startupProbe,omitempty"                                       doc:"Startup probe configuration"`
	Affinity           map[string]any       `json:"affinity,omitempty"                                           doc:"Affinity configuration"`
	Tolerations        []map[string]any     `json:"tolerations,omitempty"                                        doc:"Tolerations"`
}

// CreateAPIResponse represents the response after creating a SimpleAPI CRD
type CreateAPIResponse struct {
	Message   string `json:"message"   example:"SimpleAPI resource created successfully" doc:"Response message"`
	Name      string `json:"name"      example:"simpleapi-new"                           doc:"Name of the created resource"`
	Namespace string `json:"namespace" example:"default"                                 doc:"Namespace of the created resource"`
	Kind      string `json:"kind"      example:"Simpleapi"                               doc:"Kind of the created resource"`
}

// GetAPIResponse represents a SimpleAPI resource response
type GetAPIResponse struct {
	APIVersion string         `json:"apiVersion" example:"apps.api.test/v1alpha1" doc:"API version"`
	Kind       string         `json:"kind"       example:"Simpleapi"              doc:"Resource kind"`
	Metadata   map[string]any `json:"metadata"                                    doc:"Resource metadata"`
	Spec       map[string]any `json:"spec"                                        doc:"Resource specification"`
}

// ListSimpleAPIResponse represents a list of SimpleAPI resources
type ListAPIResponse struct {
	APIVersion string           `json:"apiVersion" example:"apps.api.test/v1alpha1" doc:"API version"`
	Kind       string           `json:"kind"       example:"SimpleAPIList"          doc:"Resource kind"`
	Items      []GetAPIResponse `json:"items"                                       doc:"List of SimpleAPI resources"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"             example:"Failed to create resource" doc:"Error message"`
	Details string `json:"details,omitempty" example:"Invalid namespace"         doc:"Error details"`
}

type GetAPIInput struct {
	Kind       string `json:"kind"                 example:"Simpleapi"     doc:"Kind of the CRD"`
	Group      string `json:"group"                example:"apps.api.test" doc:"The Group name"`
	CrdVersion string `json:"crdversion,omitempty" example:"v1alpha1"      doc:"crd version"`
	Name       string `json:"name"                                         doc:"Name of the SimpleAPI resource"`
	Namespace  string `json:"namespace"                                    doc:"Namespace of the SimpleAPI resource"`
}
type ListAPIInput struct {
	Kind       string `json:"kind"                 example:"Simpleapi"     doc:"Kind of the CRD"`
	Group      string `json:"group"                example:"apps.api.test" doc:"The Group name"`
	CrdVersion string `json:"crdversion,omitempty" example:"v1alpha1"      doc:"crd version"`
	Namespace  string `json:"namespace"                                    doc:"Namespace of the SimpleAPI resource"`
}
