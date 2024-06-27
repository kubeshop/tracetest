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

// checks if the WebhookResultRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WebhookResultRequest{}

// WebhookResultRequest struct for WebhookResultRequest
type WebhookResultRequest struct {
	Url     *string      `json:"url,omitempty"`
	Headers []HTTPHeader `json:"headers,omitempty"`
	Body    *string      `json:"body,omitempty"`
}

// NewWebhookResultRequest instantiates a new WebhookResultRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhookResultRequest() *WebhookResultRequest {
	this := WebhookResultRequest{}
	return &this
}

// NewWebhookResultRequestWithDefaults instantiates a new WebhookResultRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhookResultRequestWithDefaults() *WebhookResultRequest {
	this := WebhookResultRequest{}
	return &this
}

// GetUrl returns the Url field value if set, zero value otherwise.
func (o *WebhookResultRequest) GetUrl() string {
	if o == nil || isNil(o.Url) {
		var ret string
		return ret
	}
	return *o.Url
}

// GetUrlOk returns a tuple with the Url field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultRequest) GetUrlOk() (*string, bool) {
	if o == nil || isNil(o.Url) {
		return nil, false
	}
	return o.Url, true
}

// HasUrl returns a boolean if a field has been set.
func (o *WebhookResultRequest) HasUrl() bool {
	if o != nil && !isNil(o.Url) {
		return true
	}

	return false
}

// SetUrl gets a reference to the given string and assigns it to the Url field.
func (o *WebhookResultRequest) SetUrl(v string) {
	o.Url = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise.
func (o *WebhookResultRequest) GetHeaders() []HTTPHeader {
	if o == nil || isNil(o.Headers) {
		var ret []HTTPHeader
		return ret
	}
	return o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultRequest) GetHeadersOk() ([]HTTPHeader, bool) {
	if o == nil || isNil(o.Headers) {
		return nil, false
	}
	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *WebhookResultRequest) HasHeaders() bool {
	if o != nil && !isNil(o.Headers) {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given []HTTPHeader and assigns it to the Headers field.
func (o *WebhookResultRequest) SetHeaders(v []HTTPHeader) {
	o.Headers = v
}

// GetBody returns the Body field value if set, zero value otherwise.
func (o *WebhookResultRequest) GetBody() string {
	if o == nil || isNil(o.Body) {
		var ret string
		return ret
	}
	return *o.Body
}

// GetBodyOk returns a tuple with the Body field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultRequest) GetBodyOk() (*string, bool) {
	if o == nil || isNil(o.Body) {
		return nil, false
	}
	return o.Body, true
}

// HasBody returns a boolean if a field has been set.
func (o *WebhookResultRequest) HasBody() bool {
	if o != nil && !isNil(o.Body) {
		return true
	}

	return false
}

// SetBody gets a reference to the given string and assigns it to the Body field.
func (o *WebhookResultRequest) SetBody(v string) {
	o.Body = &v
}

func (o WebhookResultRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WebhookResultRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Url) {
		toSerialize["url"] = o.Url
	}
	if !isNil(o.Headers) {
		toSerialize["headers"] = o.Headers
	}
	if !isNil(o.Body) {
		toSerialize["body"] = o.Body
	}
	return toSerialize, nil
}

type NullableWebhookResultRequest struct {
	value *WebhookResultRequest
	isSet bool
}

func (v NullableWebhookResultRequest) Get() *WebhookResultRequest {
	return v.value
}

func (v *NullableWebhookResultRequest) Set(val *WebhookResultRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhookResultRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhookResultRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhookResultRequest(val *WebhookResultRequest) *NullableWebhookResultRequest {
	return &NullableWebhookResultRequest{value: val, isSet: true}
}

func (v NullableWebhookResultRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhookResultRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}