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

// EnvironmentValue struct for EnvironmentValue
type EnvironmentValue struct {
	Key *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

// NewEnvironmentValue instantiates a new EnvironmentValue object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEnvironmentValue() *EnvironmentValue {
	this := EnvironmentValue{}
	return &this
}

// NewEnvironmentValueWithDefaults instantiates a new EnvironmentValue object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEnvironmentValueWithDefaults() *EnvironmentValue {
	this := EnvironmentValue{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *EnvironmentValue) GetKey() string {
	if o == nil || o.Key == nil {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnvironmentValue) GetKeyOk() (*string, bool) {
	if o == nil || o.Key == nil {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *EnvironmentValue) HasKey() bool {
	if o != nil && o.Key != nil {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *EnvironmentValue) SetKey(v string) {
	o.Key = &v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *EnvironmentValue) GetValue() string {
	if o == nil || o.Value == nil {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EnvironmentValue) GetValueOk() (*string, bool) {
	if o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *EnvironmentValue) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *EnvironmentValue) SetValue(v string) {
	o.Value = &v
}

func (o EnvironmentValue) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Key != nil {
		toSerialize["key"] = o.Key
	}
	if o.Value != nil {
		toSerialize["value"] = o.Value
	}
	return json.Marshal(toSerialize)
}

type NullableEnvironmentValue struct {
	value *EnvironmentValue
	isSet bool
}

func (v NullableEnvironmentValue) Get() *EnvironmentValue {
	return v.value
}

func (v *NullableEnvironmentValue) Set(val *EnvironmentValue) {
	v.value = val
	v.isSet = true
}

func (v NullableEnvironmentValue) IsSet() bool {
	return v.isSet
}

func (v *NullableEnvironmentValue) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEnvironmentValue(val *EnvironmentValue) *NullableEnvironmentValue {
	return &NullableEnvironmentValue{value: val, isSet: true}
}

func (v NullableEnvironmentValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEnvironmentValue) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


