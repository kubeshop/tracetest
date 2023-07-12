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

// checks if the GetTransactions200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetTransactions200Response{}

// GetTransactions200Response struct for GetTransactions200Response
type GetTransactions200Response struct {
	Count *int32                `json:"count,omitempty"`
	Items []TransactionResource `json:"items,omitempty"`
}

// NewGetTransactions200Response instantiates a new GetTransactions200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetTransactions200Response() *GetTransactions200Response {
	this := GetTransactions200Response{}
	return &this
}

// NewGetTransactions200ResponseWithDefaults instantiates a new GetTransactions200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetTransactions200ResponseWithDefaults() *GetTransactions200Response {
	this := GetTransactions200Response{}
	return &this
}

// GetCount returns the Count field value if set, zero value otherwise.
func (o *GetTransactions200Response) GetCount() int32 {
	if o == nil || isNil(o.Count) {
		var ret int32
		return ret
	}
	return *o.Count
}

// GetCountOk returns a tuple with the Count field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetTransactions200Response) GetCountOk() (*int32, bool) {
	if o == nil || isNil(o.Count) {
		return nil, false
	}
	return o.Count, true
}

// HasCount returns a boolean if a field has been set.
func (o *GetTransactions200Response) HasCount() bool {
	if o != nil && !isNil(o.Count) {
		return true
	}

	return false
}

// SetCount gets a reference to the given int32 and assigns it to the Count field.
func (o *GetTransactions200Response) SetCount(v int32) {
	o.Count = &v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *GetTransactions200Response) GetItems() []TransactionResource {
	if o == nil || isNil(o.Items) {
		var ret []TransactionResource
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetTransactions200Response) GetItemsOk() ([]TransactionResource, bool) {
	if o == nil || isNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *GetTransactions200Response) HasItems() bool {
	if o != nil && !isNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []TransactionResource and assigns it to the Items field.
func (o *GetTransactions200Response) SetItems(v []TransactionResource) {
	o.Items = v
}

func (o GetTransactions200Response) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetTransactions200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Count) {
		toSerialize["count"] = o.Count
	}
	if !isNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableGetTransactions200Response struct {
	value *GetTransactions200Response
	isSet bool
}

func (v NullableGetTransactions200Response) Get() *GetTransactions200Response {
	return v.value
}

func (v *NullableGetTransactions200Response) Set(val *GetTransactions200Response) {
	v.value = val
	v.isSet = true
}

func (v NullableGetTransactions200Response) IsSet() bool {
	return v.isSet
}

func (v *NullableGetTransactions200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetTransactions200Response(val *GetTransactions200Response) *NullableGetTransactions200Response {
	return &NullableGetTransactions200Response{value: val, isSet: true}
}

func (v NullableGetTransactions200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetTransactions200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
