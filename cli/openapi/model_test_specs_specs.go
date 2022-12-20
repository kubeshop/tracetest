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

// TestSpecsSpecs struct for TestSpecsSpecs
type TestSpecsSpecs struct {
	Name NullableString `json:"name,omitempty"`
	Selector *Selector `json:"selector,omitempty"`
	Assertions []string `json:"assertions,omitempty"`
}

// NewTestSpecsSpecs instantiates a new TestSpecsSpecs object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTestSpecsSpecs() *TestSpecsSpecs {
	this := TestSpecsSpecs{}
	return &this
}

// NewTestSpecsSpecsWithDefaults instantiates a new TestSpecsSpecs object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTestSpecsSpecsWithDefaults() *TestSpecsSpecs {
	this := TestSpecsSpecs{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TestSpecsSpecs) GetName() string {
	if o == nil || o.Name.Get() == nil {
		var ret string
		return ret
	}
	return *o.Name.Get()
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TestSpecsSpecs) GetNameOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return o.Name.Get(), o.Name.IsSet()
}

// HasName returns a boolean if a field has been set.
func (o *TestSpecsSpecs) HasName() bool {
	if o != nil && o.Name.IsSet() {
		return true
	}

	return false
}

// SetName gets a reference to the given NullableString and assigns it to the Name field.
func (o *TestSpecsSpecs) SetName(v string) {
	o.Name.Set(&v)
}
// SetNameNil sets the value for Name to be an explicit nil
func (o *TestSpecsSpecs) SetNameNil() {
	o.Name.Set(nil)
}

// UnsetName ensures that no value is present for Name, not even an explicit nil
func (o *TestSpecsSpecs) UnsetName() {
	o.Name.Unset()
}

// GetSelector returns the Selector field value if set, zero value otherwise.
func (o *TestSpecsSpecs) GetSelector() Selector {
	if o == nil || o.Selector == nil {
		var ret Selector
		return ret
	}
	return *o.Selector
}

// GetSelectorOk returns a tuple with the Selector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestSpecsSpecs) GetSelectorOk() (*Selector, bool) {
	if o == nil || o.Selector == nil {
		return nil, false
	}
	return o.Selector, true
}

// HasSelector returns a boolean if a field has been set.
func (o *TestSpecsSpecs) HasSelector() bool {
	if o != nil && o.Selector != nil {
		return true
	}

	return false
}

// SetSelector gets a reference to the given Selector and assigns it to the Selector field.
func (o *TestSpecsSpecs) SetSelector(v Selector) {
	o.Selector = &v
}

// GetAssertions returns the Assertions field value if set, zero value otherwise.
func (o *TestSpecsSpecs) GetAssertions() []string {
	if o == nil || o.Assertions == nil {
		var ret []string
		return ret
	}
	return o.Assertions
}

// GetAssertionsOk returns a tuple with the Assertions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestSpecsSpecs) GetAssertionsOk() ([]string, bool) {
	if o == nil || o.Assertions == nil {
		return nil, false
	}
	return o.Assertions, true
}

// HasAssertions returns a boolean if a field has been set.
func (o *TestSpecsSpecs) HasAssertions() bool {
	if o != nil && o.Assertions != nil {
		return true
	}

	return false
}

// SetAssertions gets a reference to the given []string and assigns it to the Assertions field.
func (o *TestSpecsSpecs) SetAssertions(v []string) {
	o.Assertions = v
}

func (o TestSpecsSpecs) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name.IsSet() {
		toSerialize["name"] = o.Name.Get()
	}
	if o.Selector != nil {
		toSerialize["selector"] = o.Selector
	}
	if o.Assertions != nil {
		toSerialize["assertions"] = o.Assertions
	}
	return json.Marshal(toSerialize)
}

type NullableTestSpecsSpecs struct {
	value *TestSpecsSpecs
	isSet bool
}

func (v NullableTestSpecsSpecs) Get() *TestSpecsSpecs {
	return v.value
}

func (v *NullableTestSpecsSpecs) Set(val *TestSpecsSpecs) {
	v.value = val
	v.isSet = true
}

func (v NullableTestSpecsSpecs) IsSet() bool {
	return v.isSet
}

func (v *NullableTestSpecsSpecs) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTestSpecsSpecs(val *TestSpecsSpecs) *NullableTestSpecsSpecs {
	return &NullableTestSpecsSpecs{value: val, isSet: true}
}

func (v NullableTestSpecsSpecs) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTestSpecsSpecs) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


