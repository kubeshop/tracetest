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

// HTTPAuthApiKey struct for HTTPAuthApiKey
type HTTPAuthApiKey struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
	In    *string `json:"in,omitempty"`
}

// NewHTTPAuthApiKey instantiates a new HTTPAuthApiKey object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHTTPAuthApiKey() *HTTPAuthApiKey {
	this := HTTPAuthApiKey{}
	return &this
}

// NewHTTPAuthApiKeyWithDefaults instantiates a new HTTPAuthApiKey object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHTTPAuthApiKeyWithDefaults() *HTTPAuthApiKey {
	this := HTTPAuthApiKey{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *HTTPAuthApiKey) GetKey() string {
	if o == nil || o.Key == nil {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPAuthApiKey) GetKeyOk() (*string, bool) {
	if o == nil || o.Key == nil {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *HTTPAuthApiKey) HasKey() bool {
	if o != nil && o.Key != nil {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *HTTPAuthApiKey) SetKey(v string) {
	o.Key = &v
}

// GetValue returns the Value field value if set, zero value otherwise.
func (o *HTTPAuthApiKey) GetValue() string {
	if o == nil || o.Value == nil {
		var ret string
		return ret
	}
	return *o.Value
}

// GetValueOk returns a tuple with the Value field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPAuthApiKey) GetValueOk() (*string, bool) {
	if o == nil || o.Value == nil {
		return nil, false
	}
	return o.Value, true
}

// HasValue returns a boolean if a field has been set.
func (o *HTTPAuthApiKey) HasValue() bool {
	if o != nil && o.Value != nil {
		return true
	}

	return false
}

// SetValue gets a reference to the given string and assigns it to the Value field.
func (o *HTTPAuthApiKey) SetValue(v string) {
	o.Value = &v
}

// GetIn returns the In field value if set, zero value otherwise.
func (o *HTTPAuthApiKey) GetIn() string {
	if o == nil || o.In == nil {
		var ret string
		return ret
	}
	return *o.In
}

// GetInOk returns a tuple with the In field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPAuthApiKey) GetInOk() (*string, bool) {
	if o == nil || o.In == nil {
		return nil, false
	}
	return o.In, true
}

// HasIn returns a boolean if a field has been set.
func (o *HTTPAuthApiKey) HasIn() bool {
	if o != nil && o.In != nil {
		return true
	}

	return false
}

// SetIn gets a reference to the given string and assigns it to the In field.
func (o *HTTPAuthApiKey) SetIn(v string) {
	o.In = &v
}

func (o HTTPAuthApiKey) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Key != nil {
		toSerialize["key"] = o.Key
	}
	if o.Value != nil {
		toSerialize["value"] = o.Value
	}
	if o.In != nil {
		toSerialize["in"] = o.In
	}
	return json.Marshal(toSerialize)
}

type NullableHTTPAuthApiKey struct {
	value *HTTPAuthApiKey
	isSet bool
}

func (v NullableHTTPAuthApiKey) Get() *HTTPAuthApiKey {
	return v.value
}

func (v *NullableHTTPAuthApiKey) Set(val *HTTPAuthApiKey) {
	v.value = val
	v.isSet = true
}

func (v NullableHTTPAuthApiKey) IsSet() bool {
	return v.isSet
}

func (v *NullableHTTPAuthApiKey) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHTTPAuthApiKey(val *HTTPAuthApiKey) *NullableHTTPAuthApiKey {
	return &NullableHTTPAuthApiKey{value: val, isSet: true}
}

func (v NullableHTTPAuthApiKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHTTPAuthApiKey) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
