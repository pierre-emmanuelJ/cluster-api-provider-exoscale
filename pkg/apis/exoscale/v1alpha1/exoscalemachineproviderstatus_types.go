/*
Copyright 2019 The Kubernetes Authors.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExoscaleMachineProviderStatus is the Schema for the exoscalemachineproviderstatuses API
// +k8s:openapi-gen=true
type ExoscaleMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Zone              string `json:"zone"`
	Template          string `json:"template"`
	User              string `json:"user"`
	Type              string `json:"type"`
	Disk              int64  `json:"disk"`
	SSHKey            string `json:"sshKey"`
	IPv6              bool   `json:"ipv6,omitempty"`
	SecurityGroup     string `json:"securityGroup"`
	AntiAffinityGroup string `json:"antiAffinityGroup,omitempty"`
	CloudInit         string `json:"cloudInit,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&ExoscaleMachineProviderStatus{})
}