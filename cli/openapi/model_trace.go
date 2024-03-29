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

// checks if the Trace type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Trace{}

// Trace struct for Trace
type Trace struct {
	TraceId *string `json:"traceId,omitempty"`
	Tree    *Span   `json:"tree,omitempty"`
	// flattened version, mapped as spanId -> span{}
	Flat *map[string]Span `json:"flat,omitempty"`
}

// NewTrace instantiates a new Trace object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTrace() *Trace {
	this := Trace{}
	return &this
}

// NewTraceWithDefaults instantiates a new Trace object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTraceWithDefaults() *Trace {
	this := Trace{}
	return &this
}

// GetTraceId returns the TraceId field value if set, zero value otherwise.
func (o *Trace) GetTraceId() string {
	if o == nil || isNil(o.TraceId) {
		var ret string
		return ret
	}
	return *o.TraceId
}

// GetTraceIdOk returns a tuple with the TraceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Trace) GetTraceIdOk() (*string, bool) {
	if o == nil || isNil(o.TraceId) {
		return nil, false
	}
	return o.TraceId, true
}

// HasTraceId returns a boolean if a field has been set.
func (o *Trace) HasTraceId() bool {
	if o != nil && !isNil(o.TraceId) {
		return true
	}

	return false
}

// SetTraceId gets a reference to the given string and assigns it to the TraceId field.
func (o *Trace) SetTraceId(v string) {
	o.TraceId = &v
}

// GetTree returns the Tree field value if set, zero value otherwise.
func (o *Trace) GetTree() Span {
	if o == nil || isNil(o.Tree) {
		var ret Span
		return ret
	}
	return *o.Tree
}

// GetTreeOk returns a tuple with the Tree field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Trace) GetTreeOk() (*Span, bool) {
	if o == nil || isNil(o.Tree) {
		return nil, false
	}
	return o.Tree, true
}

// HasTree returns a boolean if a field has been set.
func (o *Trace) HasTree() bool {
	if o != nil && !isNil(o.Tree) {
		return true
	}

	return false
}

// SetTree gets a reference to the given Span and assigns it to the Tree field.
func (o *Trace) SetTree(v Span) {
	o.Tree = &v
}

// GetFlat returns the Flat field value if set, zero value otherwise.
func (o *Trace) GetFlat() map[string]Span {
	if o == nil || isNil(o.Flat) {
		var ret map[string]Span
		return ret
	}
	return *o.Flat
}

// GetFlatOk returns a tuple with the Flat field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Trace) GetFlatOk() (*map[string]Span, bool) {
	if o == nil || isNil(o.Flat) {
		return nil, false
	}
	return o.Flat, true
}

// HasFlat returns a boolean if a field has been set.
func (o *Trace) HasFlat() bool {
	if o != nil && !isNil(o.Flat) {
		return true
	}

	return false
}

// SetFlat gets a reference to the given map[string]Span and assigns it to the Flat field.
func (o *Trace) SetFlat(v map[string]Span) {
	o.Flat = &v
}

func (o Trace) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Trace) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.TraceId) {
		toSerialize["traceId"] = o.TraceId
	}
	if !isNil(o.Tree) {
		toSerialize["tree"] = o.Tree
	}
	if !isNil(o.Flat) {
		toSerialize["flat"] = o.Flat
	}
	return toSerialize, nil
}

type NullableTrace struct {
	value *Trace
	isSet bool
}

func (v NullableTrace) Get() *Trace {
	return v.value
}

func (v *NullableTrace) Set(val *Trace) {
	v.value = val
	v.isSet = true
}

func (v NullableTrace) IsSet() bool {
	return v.isSet
}

func (v *NullableTrace) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTrace(val *Trace) *NullableTrace {
	return &NullableTrace{value: val, isSet: true}
}

func (v NullableTrace) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTrace) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
