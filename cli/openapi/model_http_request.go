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

// HTTPRequest struct for HTTPRequest
type HTTPRequest struct {
	Url     *string      `json:"url,omitempty"`
	Method  *string      `json:"method,omitempty"`
	Headers []HTTPHeader `json:"headers,omitempty"`
	Body    *string      `json:"body,omitempty"`
	Auth    *HTTPAuth    `json:"auth,omitempty"`
}

// NewHTTPRequest instantiates a new HTTPRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHTTPRequest() *HTTPRequest {
	this := HTTPRequest{}
	return &this
}

// NewHTTPRequestWithDefaults instantiates a new HTTPRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHTTPRequestWithDefaults() *HTTPRequest {
	this := HTTPRequest{}
	return &this
}

// GetUrl returns the Url field value if set, zero value otherwise.
func (o *HTTPRequest) GetUrl() string {
	if o == nil || o.Url == nil {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPRequest) GetUrlOk() (*string, bool) {
	if o == nil || o.Url == nil {
		return nil, false
	}
	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *HTTPRequest) HasUrl() bool {
	if o != nil && o.Url != nil {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
func (o *HTTPRequest) SetUrl(v string) {
	o.Url = &v
}

// GetMethod returns the Method field value if set, zero value otherwise.
func (o *HTTPRequest) GetMethod() string {
	if o == nil || o.Method == nil {
		var ret string
		return ret
	}
	return *o.Method
}

// GetMethodOk returns a tuple with the Method field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPRequest) GetMethodOk() (*string, bool) {
	if o == nil || o.Method == nil {
		return nil, false
	}
	return o.Method, true
}

// HasMethod returns a boolean if a field has been set.
func (o *HTTPRequest) HasMethod() bool {
	if o != nil && o.Method != nil {
		return true
	}

	return false
}

// SetMethod gets a reference to the given string and assigns it to the Method field.
func (o *HTTPRequest) SetMethod(v string) {
	o.Method = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise.
func (o *HTTPRequest) GetHeaders() []HTTPHeader {
	if o == nil || o.Headers == nil {
		var ret []HTTPHeader
		return ret
	}
	return o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPRequest) GetHeadersOk() ([]HTTPHeader, bool) {
	if o == nil || o.Headers == nil {
		return nil, false
	}
	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *HTTPRequest) HasHeaders() bool {
	if o != nil && o.Headers != nil {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given []HTTPHeader and assigns it to the Headers field.
func (o *HTTPRequest) SetHeaders(v []HTTPHeader) {
	o.Headers = v
}

// GetBody returns the Body field value if set, zero value otherwise.
func (o *HTTPRequest) GetBody() string {
	if o == nil || o.Body == nil {
		var ret string
		return ret
	}
	return *o.Body
}

// GetBodyOk returns a tuple with the Body field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPRequest) GetBodyOk() (*string, bool) {
	if o == nil || o.Body == nil {
		return nil, false
	}
	return o.Body, true
}

// HasBody returns a boolean if a field has been set.
func (o *HTTPRequest) HasBody() bool {
	if o != nil && o.Body != nil {
		return true
	}

	return false
}

// SetBody gets a reference to the given string and assigns it to the Body field.
func (o *HTTPRequest) SetBody(v string) {
	o.Body = &v
}

// GetAuth returns the Auth field value if set, zero value otherwise.
func (o *HTTPRequest) GetAuth() HTTPAuth {
	if o == nil || o.Auth == nil {
		var ret HTTPAuth
		return ret
	}
	return *o.Auth
}

// GetAuthOk returns a tuple with the Auth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPRequest) GetAuthOk() (*HTTPAuth, bool) {
	if o == nil || o.Auth == nil {
		return nil, false
	}
	return o.Auth, true
}

// HasAuth returns a boolean if a field has been set.
func (o *HTTPRequest) HasAuth() bool {
	if o != nil && o.Auth != nil {
		return true
	}

	return false
}

// SetAuth gets a reference to the given HTTPAuth and assigns it to the Auth field.
func (o *HTTPRequest) SetAuth(v HTTPAuth) {
	o.Auth = &v
}

func (o HTTPRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Url != nil {
		toSerialize["url"] = o.Url
	}
	if o.Method != nil {
		toSerialize["method"] = o.Method
	}
	if o.Headers != nil {
		toSerialize["headers"] = o.Headers
	}
	if o.Body != nil {
		toSerialize["body"] = o.Body
	}
	if o.Auth != nil {
		toSerialize["auth"] = o.Auth
	}
	return json.Marshal(toSerialize)
}

type NullableHTTPRequest struct {
	value *HTTPRequest
	isSet bool
}

func (v NullableHTTPRequest) Get() *HTTPRequest {
	return v.value
}

func (v *NullableHTTPRequest) Set(val *HTTPRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableHTTPRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableHTTPRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHTTPRequest(val *HTTPRequest) *NullableHTTPRequest {
	return &NullableHTTPRequest{value: val, isSet: true}
}

func (v NullableHTTPRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHTTPRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
