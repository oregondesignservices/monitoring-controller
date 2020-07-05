/*
Copyright 2020 Raising The Floor.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/url"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type FromType string

var (
	FromTypeBodyYaml FromType = "body_yaml"
	FromTypeBodyJson FromType = "body_json"
	FromTypeBodyRaw  FromType = "body_raw"
	FromTypeHeaders  FromType = "headers"
)

type Variable struct {
	// The variable name
	Name string `json:"name"`

	// Where to extract the variable from
	// +kubebuilder:validation:Enum=body_yaml;body_json;body_raw;headers
	From FromType `json:"from"`

	// The JSON path to the data.
	JsonPath string `json:"json_path,omitempty"`

	// The final value of the variable, after its been extracted
	Value string `json:"-"`
}

type VariableList []*Variable

type HttpRequest struct {
	// Name of the HTTP request. Used for debugging and metrics
	Name string `json:"name"`

	// A target service, to be used in metrics
	TargetService string `json:"target_service"`

	// The request timeout. Default is 5 seconds
	Timeout string `json:"timeout,omitempty"`

	// The HTTP method
	// +kubebuilder:validation:Enum=HEAD;GET;POST;PUT;PATCH;DELETE;OPTIONS
	Method string `json:"method"`

	// HTTP(S) URL to make the request
	Url string `json:"url"`

	// Any potential query parameters
	QueryParams url.Values `json:"query_params,omitempty"`

	// The request body
	Body string `json:"body,omitempty"`

	// Request headers
	Headers http.Header `json:"headers,omitempty"`

	// Extract variables for later requests to utilize
	Variables VariableList `json:"variables,omitempty"`

	// Expected response codes. By default, this will be anything seen as "ok"
	ExpectedResponseCodes []int `json:"expected_response_codes,omitempty"`

	// Variables available from previous requests
	availableVariables VariableList `json:"-"`
}

// HttpMonitorSpec defines the desired state of HttpMonitor
type HttpMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Requests []HttpRequest `json:"requests"`

	// How frequently to execute the monitor requests
	Period string `json:"period"`
}

// HttpMonitorStatus defines the observed state of HttpMonitor
type HttpMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// HttpMonitor is the Schema for the httpmonitors API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:spec
type HttpMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HttpMonitorSpec   `json:"spec,omitempty"`
	Status HttpMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HttpMonitorList contains a list of HttpMonitor
type HttpMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HttpMonitor{}, &HttpMonitorList{})
}
