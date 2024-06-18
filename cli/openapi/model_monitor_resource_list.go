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

// checks if the MonitorResourceList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MonitorResourceList{}

// MonitorResourceList struct for MonitorResourceList
type MonitorResourceList struct {
	Count *int32            `json:"count,omitempty"`
	Items []MonitorResource `json:"items,omitempty"`
}

// NewMonitorResourceList instantiates a new MonitorResourceList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMonitorResourceList() *MonitorResourceList {
	this := MonitorResourceList{}
	return &this
}

// NewMonitorResourceListWithDefaults instantiates a new MonitorResourceList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMonitorResourceListWithDefaults() *MonitorResourceList {
	this := MonitorResourceList{}
	return &this
}

// GetCount returns the Count field value if set, zero value otherwise.
func (o *MonitorResourceList) GetCount() int32 {
	if o == nil || isNil(o.Count) {
		var ret int32
		return ret
	}
	return *o.Count
}

// GetCountOk returns a tuple with the Count field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MonitorResourceList) GetCountOk() (*int32, bool) {
	if o == nil || isNil(o.Count) {
		return nil, false
	}
	return o.Count, true
}

// HasCount returns a boolean if a field has been set.
func (o *MonitorResourceList) HasCount() bool {
	if o != nil && !isNil(o.Count) {
		return true
	}

	return false
}

// SetCount gets a reference to the given int32 and assigns it to the Count field.
func (o *MonitorResourceList) SetCount(v int32) {
	o.Count = &v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *MonitorResourceList) GetItems() []MonitorResource {
	if o == nil || isNil(o.Items) {
		var ret []MonitorResource
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MonitorResourceList) GetItemsOk() ([]MonitorResource, bool) {
	if o == nil || isNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *MonitorResourceList) HasItems() bool {
	if o != nil && !isNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []MonitorResource and assigns it to the Items field.
func (o *MonitorResourceList) SetItems(v []MonitorResource) {
	o.Items = v
}

func (o MonitorResourceList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o MonitorResourceList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Count) {
		toSerialize["count"] = o.Count
	}
	if !isNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableMonitorResourceList struct {
	value *MonitorResourceList
	isSet bool
}

func (v NullableMonitorResourceList) Get() *MonitorResourceList {
	return v.value
}

func (v *NullableMonitorResourceList) Set(val *MonitorResourceList) {
	v.value = val
	v.isSet = true
}

func (v NullableMonitorResourceList) IsSet() bool {
	return v.isSet
}

func (v *NullableMonitorResourceList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMonitorResourceList(val *MonitorResourceList) *NullableMonitorResourceList {
	return &NullableMonitorResourceList{value: val, isSet: true}
}

func (v NullableMonitorResourceList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMonitorResourceList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
