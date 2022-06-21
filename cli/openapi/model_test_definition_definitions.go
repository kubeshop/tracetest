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

// TestDefinitionDefinitions struct for TestDefinitionDefinitions
type TestDefinitionDefinitions struct {
	Selector   *Selector   `json:"selector,omitempty"`
	Assertions []Assertion `json:"assertions,omitempty"`
}

// NewTestDefinitionDefinitions instantiates a new TestDefinitionDefinitions object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTestDefinitionDefinitions() *TestDefinitionDefinitions {
	this := TestDefinitionDefinitions{}
	return &this
}

// NewTestDefinitionDefinitionsWithDefaults instantiates a new TestDefinitionDefinitions object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTestDefinitionDefinitionsWithDefaults() *TestDefinitionDefinitions {
	this := TestDefinitionDefinitions{}
	return &this
}

// GetSelector returns the Selector field value if set, zero value otherwise.
func (o *TestDefinitionDefinitions) GetSelector() Selector {
	if o == nil || o.Selector == nil {
		var ret Selector
		return ret
	}
	return *o.Selector
}

// GetSelectorOk returns a tuple with the Selector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestDefinitionDefinitions) GetSelectorOk() (*Selector, bool) {
	if o == nil || o.Selector == nil {
		return nil, false
	}
	return o.Selector, true
}

// HasSelector returns a boolean if a field has been set.
func (o *TestDefinitionDefinitions) HasSelector() bool {
	if o != nil && o.Selector != nil {
		return true
	}

	return false
}

// SetSelector gets a reference to the given Selector and assigns it to the Selector field.
func (o *TestDefinitionDefinitions) SetSelector(v Selector) {
	o.Selector = &v
}

// GetAssertions returns the Assertions field value if set, zero value otherwise.
func (o *TestDefinitionDefinitions) GetAssertions() []Assertion {
	if o == nil || o.Assertions == nil {
		var ret []Assertion
		return ret
	}
	return o.Assertions
}

// GetAssertionsOk returns a tuple with the Assertions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestDefinitionDefinitions) GetAssertionsOk() ([]Assertion, bool) {
	if o == nil || o.Assertions == nil {
		return nil, false
	}
	return o.Assertions, true
}

// HasAssertions returns a boolean if a field has been set.
func (o *TestDefinitionDefinitions) HasAssertions() bool {
	if o != nil && o.Assertions != nil {
		return true
	}

	return false
}

// SetAssertions gets a reference to the given []Assertion and assigns it to the Assertions field.
func (o *TestDefinitionDefinitions) SetAssertions(v []Assertion) {
	o.Assertions = v
}

func (o TestDefinitionDefinitions) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Selector != nil {
		toSerialize["selector"] = o.Selector
	}
	if o.Assertions != nil {
		toSerialize["assertions"] = o.Assertions
	}
	return json.Marshal(toSerialize)
}

type NullableTestDefinitionDefinitions struct {
	value *TestDefinitionDefinitions
	isSet bool
}

func (v NullableTestDefinitionDefinitions) Get() *TestDefinitionDefinitions {
	return v.value
}

func (v *NullableTestDefinitionDefinitions) Set(val *TestDefinitionDefinitions) {
	v.value = val
	v.isSet = true
}

func (v NullableTestDefinitionDefinitions) IsSet() bool {
	return v.isSet
}

func (v *NullableTestDefinitionDefinitions) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTestDefinitionDefinitions(val *TestDefinitionDefinitions) *NullableTestDefinitionDefinitions {
	return &NullableTestDefinitionDefinitions{value: val, isSet: true}
}

func (v NullableTestDefinitionDefinitions) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTestDefinitionDefinitions) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
