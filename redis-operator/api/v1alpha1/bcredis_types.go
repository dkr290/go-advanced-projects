/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BcredisSpec defines the desired state of Bcredis.
type BcredisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// REdis image is the deocker image to use for Redis pods
	RedisImage string `json:"redisImage,omitempty"`
	// StorageClassName is the storage class in azure for PVC
	StorageClassName string `json:"storageClassName,omitempty"`
	// EnvoyGatewayClassName is the name of the Envoy GatewayClass to use.
	// +kubebuilder:default="envoy-gateway"
	EnvoyGatewayClassName string `json:"envoyGatewayClassName,omitempty"`
	// StorageSize is the size of the Azure PVC for each redis instance
	StorageSize resource.Quantity `json:"storageSize,omitempty"`
	// RedisPasswordSecret is the name of the k8s secret containing the Redis password
	RedisPasswordSecret string `json:"redisPasswordSecret,omitempty"`
	// ServicePort is the port exposed by the Envoy Gateway for Redis TCP traffic.
	// +kubebuilder:default=6379
	ServicePort int32 `json:"servicePort,omitempty"`
}

// RedisRole represents the role of a Redis instance.
// +kubebuilder:validation:Enum=master;replica
type RedisRole string

const (
	RoleMaster  RedisRole = "master"
	RoleReplica RedisRole = "replica"
)

// BcredisStatus defines the observed state of Bcredis.
type BcredisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// MasterPod is tha name of the pod currently acting as master.
	MasterPod string `json:"masterPod,omitempty"`
	// ReplicaPod is the pod currently acting as replica.
	ReplicaPod string `json:"replicaPod,omitempty"`
	// CurrentMasterService is the service name for the master pod.
	CurrentMasterService string `json:"currentMasterService,omitempty"`
	// Phase is the current phase of the Bcredis resource
	Phase string `json:"phase,omitempty"`
	// Conditions repesent the lates available obervation of the Bcredis state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Master",type=string,JSONPath=`.status.masterPod`
// +kubebuilder:printcolumn:name="Replica",type=string,JSONPath=`.status.replicaPod`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// Bcredis is the Schema for the bcredis API.
type Bcredis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BcredisSpec   `json:"spec,omitempty"`
	Status BcredisStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BcredisList contains a list of Bcredis.
type BcredisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bcredis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bcredis{}, &BcredisList{})
}
