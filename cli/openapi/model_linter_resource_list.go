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

// checks if the LinterResourceList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LinterResourceList{}

// LinterResourceList struct for LinterResourceList
type LinterResourceList struct {
	Items []LinterResource `json:"items,omitempty"`
}

// NewLinterResourceList instantiates a new LinterResourceList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinterResourceList() *LinterResourceList {
	this := LinterResourceList{}
	return &this
}

// NewLinterResourceListWithDefaults instantiates a new LinterResourceList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinterResourceListWithDefaults() *LinterResourceList {
	this := LinterResourceList{}
	return &this
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *LinterResourceList) GetItems() []LinterResource {
	if o == nil || isNil(o.Items) {
		var ret []LinterResource
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceList) GetItemsOk() ([]LinterResource, bool) {
	if o == nil || isNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *LinterResourceList) HasItems() bool {
	if o != nil && !isNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []LinterResource and assigns it to the Items field.
func (o *LinterResourceList) SetItems(v []LinterResource) {
	o.Items = v
}

func (o LinterResourceList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LinterResourceList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableLinterResourceList struct {
	value *LinterResourceList
	isSet bool
}

func (v NullableLinterResourceList) Get() *LinterResourceList {
	return v.value
}

func (v *NullableLinterResourceList) Set(val *LinterResourceList) {
	v.value = val
	v.isSet = true
}

func (v NullableLinterResourceList) IsSet() bool {
	return v.isSet
}

func (v *NullableLinterResourceList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinterResourceList(val *LinterResourceList) *NullableLinterResourceList {
	return &NullableLinterResourceList{value: val, isSet: true}
}

func (v NullableLinterResourceList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinterResourceList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
