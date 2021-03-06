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

// HTTPAuthBasic struct for HTTPAuthBasic
type HTTPAuthBasic struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

// NewHTTPAuthBasic instantiates a new HTTPAuthBasic object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHTTPAuthBasic() *HTTPAuthBasic {
	this := HTTPAuthBasic{}
	return &this
}

// NewHTTPAuthBasicWithDefaults instantiates a new HTTPAuthBasic object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHTTPAuthBasicWithDefaults() *HTTPAuthBasic {
	this := HTTPAuthBasic{}
	return &this
}

// GetUsername returns the Username field value if set, zero value otherwise.
func (o *HTTPAuthBasic) GetUsername() string {
	if o == nil || o.Username == nil {
		var ret string
		return ret
	}
	return *o.Username
}

// GetUsernameOk returns a tuple with the Username field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPAuthBasic) GetUsernameOk() (*string, bool) {
	if o == nil || o.Username == nil {
		return nil, false
	}
	return o.Username, true
}

// HasUsername returns a boolean if a field has been set.
func (o *HTTPAuthBasic) HasUsername() bool {
	if o != nil && o.Username != nil {
		return true
	}

	return false
}

// SetUsername gets a reference to the given string and assigns it to the Username field.
func (o *HTTPAuthBasic) SetUsername(v string) {
	o.Username = &v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *HTTPAuthBasic) GetPassword() string {
	if o == nil || o.Password == nil {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPAuthBasic) GetPasswordOk() (*string, bool) {
	if o == nil || o.Password == nil {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *HTTPAuthBasic) HasPassword() bool {
	if o != nil && o.Password != nil {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *HTTPAuthBasic) SetPassword(v string) {
	o.Password = &v
}

func (o HTTPAuthBasic) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Username != nil {
		toSerialize["username"] = o.Username
	}
	if o.Password != nil {
		toSerialize["password"] = o.Password
	}
	return json.Marshal(toSerialize)
}

type NullableHTTPAuthBasic struct {
	value *HTTPAuthBasic
	isSet bool
}

func (v NullableHTTPAuthBasic) Get() *HTTPAuthBasic {
	return v.value
}

func (v *NullableHTTPAuthBasic) Set(val *HTTPAuthBasic) {
	v.value = val
	v.isSet = true
}

func (v NullableHTTPAuthBasic) IsSet() bool {
	return v.isSet
}

func (v *NullableHTTPAuthBasic) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHTTPAuthBasic(val *HTTPAuthBasic) *NullableHTTPAuthBasic {
	return &NullableHTTPAuthBasic{value: val, isSet: true}
}

func (v NullableHTTPAuthBasic) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHTTPAuthBasic) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
