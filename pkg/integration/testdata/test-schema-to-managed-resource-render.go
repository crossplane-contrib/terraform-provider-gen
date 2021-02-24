/*
	Copyright 2019 The Crossplane Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// +kubebuilder:object:root=true

// TestResource is a managed resource representing a resource mirrored in the cloud
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
type TestResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestResourceSpec   `json:"spec"`
	Status TestResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TestResource contains a list of TestResourceList
type TestResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestResource `json:"items"`
}

// A TestResourceSpec defines the desired state of a TestResource
type TestResourceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TestResourceParameters `json:"forProvider"`
}

// A TestResourceParameters defines the desired state of a TestResource
type TestResourceParameters struct {
	Labels                 map[string]string `json:"labels,omitempty"`
	NumberList             []int64           `json:"number_list,omitempty"`
	RequiredName           string            `json:"required_name"`
	DifferentResourceRefId string            `json:"different_resource_ref_id"`
	PerformOptionalAction  bool              `json:"perform_optional_action"`
}

// A TestResourceStatus defines the observed state of a TestResource
type TestResourceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          TestResourceObservation `json:"atProvider"`
}

// A TestResourceObservation records the observed state of a TestResource
type TestResourceObservation struct {
	ComputedOwnerId string `json:"computed_owner_id"`
}
