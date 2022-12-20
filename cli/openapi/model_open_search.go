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

// OpenSearch struct for OpenSearch
type OpenSearch struct {
	Addresses []string `json:"addresses,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Index *string `json:"index,omitempty"`
}

// NewOpenSearch instantiates a new OpenSearch object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewOpenSearch() *OpenSearch {
	this := OpenSearch{}
	return &this
}

// NewOpenSearchWithDefaults instantiates a new OpenSearch object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewOpenSearchWithDefaults() *OpenSearch {
	this := OpenSearch{}
	return &this
}

// GetAddresses returns the Addresses field value if set, zero value otherwise.
func (o *OpenSearch) GetAddresses() []string {
	if o == nil || o.Addresses == nil {
		var ret []string
		return ret
	}
	return o.Addresses
}

// GetAddressesOk returns a tuple with the Addresses field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OpenSearch) GetAddressesOk() ([]string, bool) {
	if o == nil || o.Addresses == nil {
		return nil, false
	}
	return o.Addresses, true
}

// HasAddresses returns a boolean if a field has been set.
func (o *OpenSearch) HasAddresses() bool {
	if o != nil && o.Addresses != nil {
		return true
	}

	return false
}

// SetAddresses gets a reference to the given []string and assigns it to the Addresses field.
func (o *OpenSearch) SetAddresses(v []string) {
	o.Addresses = v
}

// GetUsername returns the Username field value if set, zero value otherwise.
func (o *OpenSearch) GetUsername() string {
	if o == nil || o.Username == nil {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OpenSearch) GetUsernameOk() (*string, bool) {
	if o == nil || o.Username == nil {
		return nil, false
	}
	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *OpenSearch) HasUsername() bool {
	if o != nil && o.Username != nil {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *OpenSearch) SetUsername(v string) {
	o.Username = &v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *OpenSearch) GetPassword() string {
	if o == nil || o.Password == nil {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OpenSearch) GetPasswordOk() (*string, bool) {
	if o == nil || o.Password == nil {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *OpenSearch) HasPassword() bool {
	if o != nil && o.Password != nil {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *OpenSearch) SetPassword(v string) {
	o.Password = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *OpenSearch) GetIndex() string {
	if o == nil || o.Index == nil {
		var ret string
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *OpenSearch) GetIndexOk() (*string, bool) {
	if o == nil || o.Index == nil {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *OpenSearch) HasIndex() bool {
	if o != nil && o.Index != nil {
		return true
	}

	return false
}

// SetIndex gets a reference to the given string and assigns it to the Index field.
func (o *OpenSearch) SetIndex(v string) {
	o.Index = &v
}

func (o OpenSearch) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Addresses != nil {
		toSerialize["addresses"] = o.Addresses
	}
	if o.Username != nil {
		toSerialize["username"] = o.Username
	}
	if o.Password != nil {
		toSerialize["password"] = o.Password
	}
	if o.Index != nil {
		toSerialize["index"] = o.Index
	}
	return json.Marshal(toSerialize)
}

type NullableOpenSearch struct {
	value *OpenSearch
	isSet bool
}

func (v NullableOpenSearch) Get() *OpenSearch {
	return v.value
}

func (v *NullableOpenSearch) Set(val *OpenSearch) {
	v.value = val
	v.isSet = true
}

func (v NullableOpenSearch) IsSet() bool {
	return v.isSet
}

func (v *NullableOpenSearch) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableOpenSearch(val *OpenSearch) *NullableOpenSearch {
	return &NullableOpenSearch{value: val, isSet: true}
}

func (v NullableOpenSearch) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableOpenSearch) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


