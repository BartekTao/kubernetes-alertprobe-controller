/*
Copyright 2023.

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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AlertProbeSpec defines the desired state of AlertProbe
type AlertProbeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// URL to check
	URL string `json:"url"`

	// Check period in seconds
	PeriodSeconds int32 `json:"periodSeconds"`
}

// AlertProbeStatus defines the observed state of AlertProbe
type AlertProbeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The result of the last check
	LastCheckResult string `json:"lastCheckResult"`

	// The time of the last check
	LastCheckTime metav1.Time `json:"lastCheckTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AlertProbe is the Schema for the alertprobes API
type AlertProbe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AlertProbeSpec   `json:"spec,omitempty"`
	Status AlertProbeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AlertProbeList contains a list of AlertProbe
type AlertProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlertProbe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlertProbe{}, &AlertProbeList{})
}
