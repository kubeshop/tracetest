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

// Selector struct for Selector
type Selector struct {
	Query *string `json:"query,omitempty"`
	Structure []SpanSelector `json:"structure,omitempty"`
}

// NewSelector instantiates a new Selector object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSelector() *Selector {
	this := Selector{}
	return &this
}

// NewSelectorWithDefaults instantiates a new Selector object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSelectorWithDefaults() *Selector {
	this := Selector{}
	return &this
}

// GetQuery returns the Query field value if set, zero value otherwise.
func (o *Selector) GetQuery() string {
	if o == nil || o.Query == nil {
		var ret string
		return ret
	}
	return *o.Query
}

// GetQueryOk returns a tuple with the Query field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Selector) GetQueryOk() (*string, bool) {
	if o == nil || o.Query == nil {
		return nil, false
	}
	return o.Query, true
}

// HasQuery returns a boolean if a field has been set.
func (o *Selector) HasQuery() bool {
	if o != nil && o.Query != nil {
		return true
	}

	return false
}

// SetQuery gets a reference to the given string and assigns it to the Query field.
func (o *Selector) SetQuery(v string) {
	o.Query = &v
}

// GetStructure returns the Structure field value if set, zero value otherwise.
func (o *Selector) GetStructure() []SpanSelector {
	if o == nil || o.Structure == nil {
		var ret []SpanSelector
		return ret
	}
	return o.Structure
}

// GetStructureOk returns a tuple with the Structure field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Selector) GetStructureOk() ([]SpanSelector, bool) {
	if o == nil || o.Structure == nil {
		return nil, false
	}
	return o.Structure, true
}

// HasStructure returns a boolean if a field has been set.
func (o *Selector) HasStructure() bool {
	if o != nil && o.Structure != nil {
		return true
	}

	return false
}

// SetStructure gets a reference to the given []SpanSelector and assigns it to the Structure field.
func (o *Selector) SetStructure(v []SpanSelector) {
	o.Structure = v
}

func (o Selector) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Query != nil {
		toSerialize["query"] = o.Query
	}
	if o.Structure != nil {
		toSerialize["structure"] = o.Structure
	}
	return json.Marshal(toSerialize)
}

type NullableSelector struct {
	value *Selector
	isSet bool
}

func (v NullableSelector) Get() *Selector {
	return v.value
}

func (v *NullableSelector) Set(val *Selector) {
	v.value = val
	v.isSet = true
}

func (v NullableSelector) IsSet() bool {
	return v.isSet
}

func (v *NullableSelector) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSelector(val *Selector) *NullableSelector {
	return &NullableSelector{value: val, isSet: true}
}

func (v NullableSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSelector) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


