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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LoremIpsumSpec defines the desired state of LoremIpsum
type LoremIpsumSpec struct {

	// Lines defines how many lines to generate.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5
	Lines int `json:"lines,omitempty"`
	// Capitalize defines whether to capitalize the generated output.
	Capitalize bool `json:"capitalize,omitempty"`
}

// LoremIpsumStatus defines the observed state of LoremIpsum
type LoremIpsumStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// LoremIpsum is the Schema for the loremipsums API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=li
type LoremIpsum struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoremIpsumSpec   `json:"spec,omitempty"`
	Status LoremIpsumStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LoremIpsumList contains a list of LoremIpsum
type LoremIpsumList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoremIpsum `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoremIpsum{}, &LoremIpsumList{})
}
