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

// ResolveContext struct for ResolveContext
type ResolveContext struct {
	TestId        *string `json:"testId,omitempty"`
	RunId         *string `json:"runId,omitempty"`
	SpanId        *string `json:"spanId,omitempty"`
	Selector      *string `json:"selector,omitempty"`
	EnvironmentId *string `json:"environmentId,omitempty"`
}

// NewResolveContext instantiates a new ResolveContext object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewResolveContext() *ResolveContext {
	this := ResolveContext{}
	return &this
}

// NewResolveContextWithDefaults instantiates a new ResolveContext object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewResolveContextWithDefaults() *ResolveContext {
	this := ResolveContext{}
	return &this
}

// GetTestId returns the TestId field value if set, zero value otherwise.
func (o *ResolveContext) GetTestId() string {
	if o == nil || o.TestId == nil {
		var ret string
		return ret
	}
	return *o.TestId
}

// GetTestIdOk returns a tuple with the TestId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResolveContext) GetTestIdOk() (*string, bool) {
	if o == nil || o.TestId == nil {
		return nil, false
	}
	return o.TestId, true
}

// HasTestId returns a boolean if a field has been set.
func (o *ResolveContext) HasTestId() bool {
	if o != nil && o.TestId != nil {
		return true
	}

	return false
}

// SetTestId gets a reference to the given string and assigns it to the TestId field.
func (o *ResolveContext) SetTestId(v string) {
	o.TestId = &v
}

// GetRunId returns the RunId field value if set, zero value otherwise.
func (o *ResolveContext) GetRunId() string {
	if o == nil || o.RunId == nil {
		var ret string
		return ret
	}
	return *o.RunId
}

// GetRunIdOk returns a tuple with the RunId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResolveContext) GetRunIdOk() (*string, bool) {
	if o == nil || o.RunId == nil {
		return nil, false
	}
	return o.RunId, true
}

// HasRunId returns a boolean if a field has been set.
func (o *ResolveContext) HasRunId() bool {
	if o != nil && o.RunId != nil {
		return true
	}

	return false
}

// SetRunId gets a reference to the given string and assigns it to the RunId field.
func (o *ResolveContext) SetRunId(v string) {
	o.RunId = &v
}

// GetSpanId returns the SpanId field value if set, zero value otherwise.
func (o *ResolveContext) GetSpanId() string {
	if o == nil || o.SpanId == nil {
		var ret string
		return ret
	}
	return *o.SpanId
}

// GetSpanIdOk returns a tuple with the SpanId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResolveContext) GetSpanIdOk() (*string, bool) {
	if o == nil || o.SpanId == nil {
		return nil, false
	}
	return o.SpanId, true
}

// HasSpanId returns a boolean if a field has been set.
func (o *ResolveContext) HasSpanId() bool {
	if o != nil && o.SpanId != nil {
		return true
	}

	return false
}

// SetSpanId gets a reference to the given string and assigns it to the SpanId field.
func (o *ResolveContext) SetSpanId(v string) {
	o.SpanId = &v
}

// GetSelector returns the Selector field value if set, zero value otherwise.
func (o *ResolveContext) GetSelector() string {
	if o == nil || o.Selector == nil {
		var ret string
		return ret
	}
	return *o.Selector
}

// GetSelectorOk returns a tuple with the Selector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResolveContext) GetSelectorOk() (*string, bool) {
	if o == nil || o.Selector == nil {
		return nil, false
	}
	return o.Selector, true
}

// HasSelector returns a boolean if a field has been set.
func (o *ResolveContext) HasSelector() bool {
	if o != nil && o.Selector != nil {
		return true
	}

	return false
}

// SetSelector gets a reference to the given string and assigns it to the Selector field.
func (o *ResolveContext) SetSelector(v string) {
	o.Selector = &v
}

// GetEnvironmentId returns the EnvironmentId field value if set, zero value otherwise.
func (o *ResolveContext) GetEnvironmentId() string {
	if o == nil || o.EnvironmentId == nil {
		var ret string
		return ret
	}
	return *o.EnvironmentId
}

// GetEnvironmentIdOk returns a tuple with the EnvironmentId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ResolveContext) GetEnvironmentIdOk() (*string, bool) {
	if o == nil || o.EnvironmentId == nil {
		return nil, false
	}
	return o.EnvironmentId, true
}

// HasEnvironmentId returns a boolean if a field has been set.
func (o *ResolveContext) HasEnvironmentId() bool {
	if o != nil && o.EnvironmentId != nil {
		return true
	}

	return false
}

// SetEnvironmentId gets a reference to the given string and assigns it to the EnvironmentId field.
func (o *ResolveContext) SetEnvironmentId(v string) {
	o.EnvironmentId = &v
}

func (o ResolveContext) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.TestId != nil {
		toSerialize["testId"] = o.TestId
	}
	if o.RunId != nil {
		toSerialize["runId"] = o.RunId
	}
	if o.SpanId != nil {
		toSerialize["spanId"] = o.SpanId
	}
	if o.Selector != nil {
		toSerialize["selector"] = o.Selector
	}
	if o.EnvironmentId != nil {
		toSerialize["environmentId"] = o.EnvironmentId
	}
	return json.Marshal(toSerialize)
}

type NullableResolveContext struct {
	value *ResolveContext
	isSet bool
}

func (v NullableResolveContext) Get() *ResolveContext {
	return v.value
}

func (v *NullableResolveContext) Set(val *ResolveContext) {
	v.value = val
	v.isSet = true
}

func (v NullableResolveContext) IsSet() bool {
	return v.isSet
}

func (v *NullableResolveContext) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResolveContext(val *ResolveContext) *NullableResolveContext {
	return &NullableResolveContext{value: val, isSet: true}
}

func (v NullableResolveContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResolveContext) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
