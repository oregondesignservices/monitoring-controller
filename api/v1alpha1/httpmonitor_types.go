/*
Copyright 2020 Raising the Floor - International

Licensed under the New BSD license. You may not use this file except in
compliance with this License.

You may obtain a copy of the License at
https://github.com/GPII/universal/blob/master/LICENSE.txt

The R&D leading to these results received funding from the:
* Rehabilitation Services Administration, US Dept. of Education under
  grant H421A150006 (APCP)
* National Institute on Disability, Independent Living, and
  Rehabilitation Research (NIDILRR)
* Administration for Independent Living & Dept. of Education under grants
  H133E080022 (RERC-IT) and H133E130028/90RE5003-01-00 (UIITA-RERC)
* European Union's Seventh Framework Programme (FP7/2007-2013) grant
  agreement nos. 289016 (Cloud4all) and 610510 (Prosperity4All)
* William and Flora Hewlett Foundation
* Ontario Ministry of Research and Innovation
* Canadian Foundation for Innovation
* Adobe Foundation
* Consumer Electronics Association Foundation
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
	FromTypeBodyRaw  FromType = "body_raw" //
	FromTypeHeaders  FromType = "headers"  // extract the variable from Headers
	FromTypeProvided FromType = "provided" // provided by the user
)

type Variable struct {
	// The variable name
	Name string `json:"name"`

	// Where to extract the variable from
	// +kubebuilder:validation:Enum=body_yaml;body_json;body_raw;headers;provided
	From FromType `json:"from"`

	// The JSON path to the data.
	JsonPath string `json:"json_path,omitempty"`

	// The final value of the variable, after its been extracted
	Value string `json:"value"`
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
	VariablesFromResponse VariableList `json:"vars_from_response,omitempty"`

	// Expected response codes. By default, this will be anything seen as "ok"
	ExpectedResponseCodes []int `json:"expected_response_codes,omitempty"`

	// VariablesFromResponse available from previous requests
	AvailableVariables VariableList `json:"-"`
}

// HttpMonitorSpec defines the desired state of HttpMonitor
type HttpMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Variables available to all requests from the start
	Globals map[string]string `json:"globals,omitempty"`

	Requests []HttpRequest `json:"requests"`

	// Optional requests to be run after `requests`.
	Cleanup []HttpRequest `json:"cleanup,omitempty"`

	// How frequently to execute the monitor requests
	Period *metav1.Duration `json:"period"`
}

// HttpMonitorStatus defines the observed state of HttpMonitor
type HttpMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastExecution *metav1.Time `json:"last_execution"`
	LastFailure   *metav1.Time `json:"last_failure"`
}

// HttpMonitor is the Schema for the httpmonitors API
// +kubebuilder:object:root=true
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HttpMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HttpMonitorSpec   `json:"spec,omitempty"`
	Status HttpMonitorStatus `json:"status,omitempty"`
}

// HttpMonitorList contains a list of HttpMonitor
// +kubebuilder:object:root=true
type HttpMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HttpMonitor{}, &HttpMonitorList{})
}
