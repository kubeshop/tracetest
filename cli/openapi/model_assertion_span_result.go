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

// AssertionSpanResult struct for AssertionSpanResult
type AssertionSpanResult struct {
	SpanId *string `json:"spanId,omitempty"`
	ObservedValue *string `json:"observedValue,omitempty"`
	Passed *bool `json:"passed,omitempty"`
	Error *string `json:"error,omitempty"`
}

// NewAssertionSpanResult instantiates a new AssertionSpanResult object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAssertionSpanResult() *AssertionSpanResult {
	this := AssertionSpanResult{}
	return &this
}

// NewAssertionSpanResultWithDefaults instantiates a new AssertionSpanResult object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAssertionSpanResultWithDefaults() *AssertionSpanResult {
	this := AssertionSpanResult{}
	return &this
}

// GetSpanId returns the SpanId field value if set, zero value otherwise.
func (o *AssertionSpanResult) GetSpanId() string {
	if o == nil || o.SpanId == nil {
		var ret string
		return ret
	}
	return *o.SpanId
}

// GetSpanIdOk returns a tuple with the SpanId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssertionSpanResult) GetSpanIdOk() (*string, bool) {
	if o == nil || o.SpanId == nil {
		return nil, false
	}
	return o.SpanId, true
}

// HasSpanId returns a boolean if a field has been set.
func (o *AssertionSpanResult) HasSpanId() bool {
	if o != nil && o.SpanId != nil {
		return true
	}

	return false
}

// SetSpanId gets a reference to the given string and assigns it to the SpanId field.
func (o *AssertionSpanResult) SetSpanId(v string) {
	o.SpanId = &v
}

// GetObservedValue returns the ObservedValue field value if set, zero value otherwise.
func (o *AssertionSpanResult) GetObservedValue() string {
	if o == nil || o.ObservedValue == nil {
		var ret string
		return ret
	}
	return *o.ObservedValue
}

// GetObservedValueOk returns a tuple with the ObservedValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssertionSpanResult) GetObservedValueOk() (*string, bool) {
	if o == nil || o.ObservedValue == nil {
		return nil, false
	}
	return o.ObservedValue, true
}

// HasObservedValue returns a boolean if a field has been set.
func (o *AssertionSpanResult) HasObservedValue() bool {
	if o != nil && o.ObservedValue != nil {
		return true
	}

	return false
}

// SetObservedValue gets a reference to the given string and assigns it to the ObservedValue field.
func (o *AssertionSpanResult) SetObservedValue(v string) {
	o.ObservedValue = &v
}

// GetPassed returns the Passed field value if set, zero value otherwise.
func (o *AssertionSpanResult) GetPassed() bool {
	if o == nil || o.Passed == nil {
		var ret bool
		return ret
	}
	return *o.Passed
}

// GetPassedOk returns a tuple with the Passed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssertionSpanResult) GetPassedOk() (*bool, bool) {
	if o == nil || o.Passed == nil {
		return nil, false
	}
	return o.Passed, true
}

// HasPassed returns a boolean if a field has been set.
func (o *AssertionSpanResult) HasPassed() bool {
	if o != nil && o.Passed != nil {
		return true
	}

	return false
}

// SetPassed gets a reference to the given bool and assigns it to the Passed field.
func (o *AssertionSpanResult) SetPassed(v bool) {
	o.Passed = &v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *AssertionSpanResult) GetError() string {
	if o == nil || o.Error == nil {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AssertionSpanResult) GetErrorOk() (*string, bool) {
	if o == nil || o.Error == nil {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *AssertionSpanResult) HasError() bool {
	if o != nil && o.Error != nil {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *AssertionSpanResult) SetError(v string) {
	o.Error = &v
}

func (o AssertionSpanResult) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.SpanId != nil {
		toSerialize["spanId"] = o.SpanId
	}
	if o.ObservedValue != nil {
		toSerialize["observedValue"] = o.ObservedValue
	}
	if o.Passed != nil {
		toSerialize["passed"] = o.Passed
	}
	if o.Error != nil {
		toSerialize["error"] = o.Error
	}
	return json.Marshal(toSerialize)
}

type NullableAssertionSpanResult struct {
	value *AssertionSpanResult
	isSet bool
}

func (v NullableAssertionSpanResult) Get() *AssertionSpanResult {
	return v.value
}

func (v *NullableAssertionSpanResult) Set(val *AssertionSpanResult) {
	v.value = val
	v.isSet = true
}

func (v NullableAssertionSpanResult) IsSet() bool {
	return v.isSet
}

func (v *NullableAssertionSpanResult) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAssertionSpanResult(val *AssertionSpanResult) *NullableAssertionSpanResult {
	return &NullableAssertionSpanResult{value: val, isSet: true}
}

func (v NullableAssertionSpanResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAssertionSpanResult) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


