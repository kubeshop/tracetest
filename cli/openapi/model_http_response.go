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

// HTTPResponse struct for HTTPResponse
type HTTPResponse struct {
	Status *string `json:"status,omitempty"`
	StatusCode *int32 `json:"statusCode,omitempty"`
	Headers []HTTPHeader `json:"headers,omitempty"`
	Body *string `json:"body,omitempty"`
}

// NewHTTPResponse instantiates a new HTTPResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHTTPResponse() *HTTPResponse {
	this := HTTPResponse{}
	return &this
}

// NewHTTPResponseWithDefaults instantiates a new HTTPResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHTTPResponseWithDefaults() *HTTPResponse {
	this := HTTPResponse{}
	return &this
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *HTTPResponse) GetStatus() string {
	if o == nil || o.Status == nil {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPResponse) GetStatusOk() (*string, bool) {
	if o == nil || o.Status == nil {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *HTTPResponse) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *HTTPResponse) SetStatus(v string) {
	o.Status = &v
}

// GetStatusCode returns the StatusCode field value if set, zero value otherwise.
func (o *HTTPResponse) GetStatusCode() int32 {
	if o == nil || o.StatusCode == nil {
		var ret int32
		return ret
	}
	return *o.StatusCode
}

// GetStatusCodeOk returns a tuple with the StatusCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPResponse) GetStatusCodeOk() (*int32, bool) {
	if o == nil || o.StatusCode == nil {
		return nil, false
	}
	return o.StatusCode, true
}

// HasStatusCode returns a boolean if a field has been set.
func (o *HTTPResponse) HasStatusCode() bool {
	if o != nil && o.StatusCode != nil {
		return true
	}

	return false
}

// SetStatusCode gets a reference to the given int32 and assigns it to the StatusCode field.
func (o *HTTPResponse) SetStatusCode(v int32) {
	o.StatusCode = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise.
func (o *HTTPResponse) GetHeaders() []HTTPHeader {
	if o == nil || o.Headers == nil {
		var ret []HTTPHeader
		return ret
	}
	return o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPResponse) GetHeadersOk() ([]HTTPHeader, bool) {
	if o == nil || o.Headers == nil {
		return nil, false
	}
	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *HTTPResponse) HasHeaders() bool {
	if o != nil && o.Headers != nil {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given []HTTPHeader and assigns it to the Headers field.
func (o *HTTPResponse) SetHeaders(v []HTTPHeader) {
	o.Headers = v
}

// GetBody returns the Body field value if set, zero value otherwise.
func (o *HTTPResponse) GetBody() string {
	if o == nil || o.Body == nil {
		var ret string
		return ret
	}
	return *o.Body
}

// GetBodyOk returns a tuple with the Body field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HTTPResponse) GetBodyOk() (*string, bool) {
	if o == nil || o.Body == nil {
		return nil, false
	}
	return o.Body, true
}

// HasBody returns a boolean if a field has been set.
func (o *HTTPResponse) HasBody() bool {
	if o != nil && o.Body != nil {
		return true
	}

	return false
}

// SetBody gets a reference to the given string and assigns it to the Body field.
func (o *HTTPResponse) SetBody(v string) {
	o.Body = &v
}

func (o HTTPResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}
	if o.StatusCode != nil {
		toSerialize["statusCode"] = o.StatusCode
	}
	if o.Headers != nil {
		toSerialize["headers"] = o.Headers
	}
	if o.Body != nil {
		toSerialize["body"] = o.Body
	}
	return json.Marshal(toSerialize)
}

type NullableHTTPResponse struct {
	value *HTTPResponse
	isSet bool
}

func (v NullableHTTPResponse) Get() *HTTPResponse {
	return v.value
}

func (v *NullableHTTPResponse) Set(val *HTTPResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableHTTPResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableHTTPResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHTTPResponse(val *HTTPResponse) *NullableHTTPResponse {
	return &NullableHTTPResponse{value: val, isSet: true}
}

func (v NullableHTTPResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHTTPResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


