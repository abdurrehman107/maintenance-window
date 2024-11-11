/*
Copyright 2024.

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

package v1

import (
	// "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MaintenanceWindowSpec defines the desired state of MaintenanceWindow.
type MaintenanceWindowSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	StartTime metav1.Time `json:"startTime"`
	EndTime   metav1.Time `json:"endTime"`
	Enabled   bool        `json:"enabled"`
	// +optional
	Resources []*metav1.APIResource `json:"resources"`
}

// ResourceType defines the type of the resource
type ResourceType struct {
	Kind      string
	Name      string
	Namespace string
}

// MaintenanceWindowStatus defines the observed state of MaintenanceWindow.
type MaintenanceWindowStatus struct {
	Active bool `json:"active"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MaintenanceWindow is the Schema for the maintenancewindows API.
type MaintenanceWindow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MaintenanceWindowSpec   `json:"spec,omitempty"`
	Status MaintenanceWindowStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MaintenanceWindowList contains a list of MaintenanceWindow.
type MaintenanceWindowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MaintenanceWindow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MaintenanceWindow{}, &MaintenanceWindowList{})
}
