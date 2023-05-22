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

const (
	Running = "Running"
	Failed  = "Failed"
	Healthy = "Healthy"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProductSpec defines the desired state of Product
type ProductSpec struct {
	// ContainerOrchestration ContainerOrchestration `json:"containerOrchestration,omitempty"`
	// Infrastructure         Infrastructure         `json:"infrastructure,omitempty"`
}

// ProductStatus defines the observed state of Product
type ProductStatus struct {
	ObservedContainerOrchestration string    `json:"observedContainerOrchestraion"`
	ObservedInfrasturcture         string    `json:"observedInfrastructure"`
	Condition                      Condition `json:"condition"`
}

type Condition struct {
	ContainerOrchestration string `json:"containerOrchestration"`
	Infrastructure         string `json:"infrastructure"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Product is the Schema for the products API
type Product struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductSpec   `json:"spec,omitempty"`
	Status ProductStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProductList contains a list of Product
type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Product `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Product{}, &ProductList{})
}
