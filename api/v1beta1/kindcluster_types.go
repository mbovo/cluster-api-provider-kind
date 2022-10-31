/*
Copyright 2022 Manuel Bovo.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// Be sure to be called before object removal from the apiserver.
	KindClusterFinalizer = "kindcluster.infrastructure.cluster.x-k8s.io"
)

// KindClusterSpec defines the desired state of KindCluster
type KindClusterSpec struct {

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	//+optional
	//+nullable
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:default=1
	ControlPlaneCount int32 `json:"controlPlaneCount,omitempty"`

	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:default=0
	WorkerCount int32 `json:"workerCount,omitempty"`

	K8sVersion string `json:"k8sVersion,omitempty"`

	// KIND image to use, see https://github.com/kubernetes-sigs/kind/releases for a list

	//+kubebuilder:default="kindest/node:v1.25.2@sha256:9be91e9e9cdf116809841fc77ebdb8845443c4c72fe5218f3ae9eb57fdb4bace"
	Image string `json:"image,omitempty"`
}

// KindClusterStatus defines the observed state of KindCluster
type KindClusterStatus struct {

	// Cluster readiness
	//+kubebuilder:default=false
	Ready bool `json:"ready"`

	// Number of nodes ready in the cluster
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:default=0
	Nodes int32 `json:"nodes,omitempty"`
}

// KindCluster is the Schema for the kindclusters API
// +kubebuilder:printcolumn:name="ready",type=boolean,JSONPath=`.status.ready`,description="cluster readiness"
// +kubebuilder:printcolumn:name="created",type=date,JSONPath=`.metadata.creationTimestamp`,description="Creatiion timestamp"
// +kubebuilder:printcolumn:name="workers",type=integer,JSONPath=`.spec.workerCount`,priority=10,description="Number of workers nodes "
// +kubebuilder:printcolumn:name="controlplane",type=integer,JSONPath=`.spec.controlPlaneCount`,priority=15,description="Number of nodes in control plane"
// +kubebuilder:printcolumn:name="version",type=string,JSONPath=`.spec.k8sVersion`,priority=20,description="Kubernetes version"
// +kubebuilder:printcolumn:name="image",type=string,JSONPath=`.spec.image`,priority=25,description="Kind image used"
// +kubebuilder:resource:shortName={kc,kcl}
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type KindCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KindClusterSpec   `json:"spec,omitempty"`
	Status KindClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// KindClusterList contains a list of KindCluster
type KindClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KindCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KindCluster{}, &KindClusterList{})
}
