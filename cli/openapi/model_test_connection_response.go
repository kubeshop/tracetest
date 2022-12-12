/*
TraceTest

OpenAPI definition for TraceTest endpoint and resources

API version: 0.2.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// TestConnectionResponse struct for TestConnectionResponse
type TestConnectionResponse struct {
	Successful   *bool   `json:"successful,omitempty"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// NewTestConnectionResponse instantiates a new TestConnectionResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTestConnectionResponse() *TestConnectionResponse {
	this := TestConnectionResponse{}
	return &this
}

// NewTestConnectionResponseWithDefaults instantiates a new TestConnectionResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTestConnectionResponseWithDefaults() *TestConnectionResponse {
	this := TestConnectionResponse{}
	return &this
}

// GetSuccessful returns the Successful field value if set, zero value otherwise.
func (o *TestConnectionResponse) GetSuccessful() bool {
	if o == nil || o.Successful == nil {
		var ret bool
		return ret
	}
	return *o.Successful
}

// GetSuccessfulOk returns a tuple with the Successful field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestConnectionResponse) GetSuccessfulOk() (*bool, bool) {
	if o == nil || o.Successful == nil {
		return nil, false
	}
	return o.Successful, true
}

// HasSuccessful returns a boolean if a field has been set.
func (o *TestConnectionResponse) HasSuccessful() bool {
	if o != nil && o.Successful != nil {
		return true
	}

	return false
}

// SetSuccessful gets a reference to the given bool and assigns it to the Successful field.
func (o *TestConnectionResponse) SetSuccessful(v bool) {
	o.Successful = &v
}

// GetErrorMessage returns the ErrorMessage field value if set, zero value otherwise.
func (o *TestConnectionResponse) GetErrorMessage() string {
	if o == nil || o.ErrorMessage == nil {
		var ret string
		return ret
	}
	return *o.ErrorMessage
}

// GetErrorMessageOk returns a tuple with the ErrorMessage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestConnectionResponse) GetErrorMessageOk() (*string, bool) {
	if o == nil || o.ErrorMessage == nil {
		return nil, false
	}
	return o.ErrorMessage, true
}

// HasErrorMessage returns a boolean if a field has been set.
func (o *TestConnectionResponse) HasErrorMessage() bool {
	if o != nil && o.ErrorMessage != nil {
		return true
	}

	return false
}

// SetErrorMessage gets a reference to the given string and assigns it to the ErrorMessage field.
func (o *TestConnectionResponse) SetErrorMessage(v string) {
	o.ErrorMessage = &v
}

func (o TestConnectionResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Successful != nil {
		toSerialize["successful"] = o.Successful
	}
	if o.ErrorMessage != nil {
		toSerialize["errorMessage"] = o.ErrorMessage
	}
	return json.Marshal(toSerialize)
}

type NullableTestConnectionResponse struct {
	value *TestConnectionResponse
	isSet bool
}

func (v NullableTestConnectionResponse) Get() *TestConnectionResponse {
	return v.value
}

func (v *NullableTestConnectionResponse) Set(val *TestConnectionResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTestConnectionResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTestConnectionResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTestConnectionResponse(val *TestConnectionResponse) *NullableTestConnectionResponse {
	return &NullableTestConnectionResponse{value: val, isSet: true}
}

func (v NullableTestConnectionResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTestConnectionResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
