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

// checks if the WebhookResultResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WebhookResultResponse{}

// WebhookResultResponse struct for WebhookResultResponse
type WebhookResultResponse struct {
	StatusCode *int32       `json:"statusCode,omitempty"`
	Status     *string      `json:"status,omitempty"`
	Body       *string      `json:"body,omitempty"`
	Headers    []HTTPHeader `json:"headers,omitempty"`
	Error      *string      `json:"error,omitempty"`
}

// NewWebhookResultResponse instantiates a new WebhookResultResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWebhookResultResponse() *WebhookResultResponse {
	this := WebhookResultResponse{}
	return &this
}

// NewWebhookResultResponseWithDefaults instantiates a new WebhookResultResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWebhookResultResponseWithDefaults() *WebhookResultResponse {
	this := WebhookResultResponse{}
	return &this
}

// GetStatusCode returns the StatusCode field value if set, zero value otherwise.
func (o *WebhookResultResponse) GetStatusCode() int32 {
	if o == nil || isNil(o.StatusCode) {
		var ret int32
		return ret
	}
	return *o.StatusCode
}

// GetStatusCodeOk returns a tuple with the StatusCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultResponse) GetStatusCodeOk() (*int32, bool) {
	if o == nil || isNil(o.StatusCode) {
		return nil, false
	}
	return o.StatusCode, true
}

// HasStatusCode returns a boolean if a field has been set.
func (o *WebhookResultResponse) HasStatusCode() bool {
	if o != nil && !isNil(o.StatusCode) {
		return true
	}

	return false
}

// SetStatusCode gets a reference to the given int32 and assigns it to the StatusCode field.
func (o *WebhookResultResponse) SetStatusCode(v int32) {
	o.StatusCode = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *WebhookResultResponse) GetStatus() string {
	if o == nil || isNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultResponse) GetStatusOk() (*string, bool) {
	if o == nil || isNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *WebhookResultResponse) HasStatus() bool {
	if o != nil && !isNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *WebhookResultResponse) SetStatus(v string) {
	o.Status = &v
}

// GetBody returns the Body field value if set, zero value otherwise.
func (o *WebhookResultResponse) GetBody() string {
	if o == nil || isNil(o.Body) {
		var ret string
		return ret
	}
	return *o.Body
}

// GetBodyOk returns a tuple with the Body field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultResponse) GetBodyOk() (*string, bool) {
	if o == nil || isNil(o.Body) {
		return nil, false
	}
	return o.Body, true
}

// HasBody returns a boolean if a field has been set.
func (o *WebhookResultResponse) HasBody() bool {
	if o != nil && !isNil(o.Body) {
		return true
	}

	return false
}

// SetBody gets a reference to the given string and assigns it to the Body field.
func (o *WebhookResultResponse) SetBody(v string) {
	o.Body = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise.
func (o *WebhookResultResponse) GetHeaders() []HTTPHeader {
	if o == nil || isNil(o.Headers) {
		var ret []HTTPHeader
		return ret
	}
	return o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultResponse) GetHeadersOk() ([]HTTPHeader, bool) {
	if o == nil || isNil(o.Headers) {
		return nil, false
	}
	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *WebhookResultResponse) HasHeaders() bool {
	if o != nil && !isNil(o.Headers) {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given []HTTPHeader and assigns it to the Headers field.
func (o *WebhookResultResponse) SetHeaders(v []HTTPHeader) {
	o.Headers = v
}

// GetError returns the Error field value if set, zero value otherwise.
func (o *WebhookResultResponse) GetError() string {
	if o == nil || isNil(o.Error) {
		var ret string
		return ret
	}
	return *o.Error
}

// GetErrorOk returns a tuple with the Error field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WebhookResultResponse) GetErrorOk() (*string, bool) {
	if o == nil || isNil(o.Error) {
		return nil, false
	}
	return o.Error, true
}

// HasError returns a boolean if a field has been set.
func (o *WebhookResultResponse) HasError() bool {
	if o != nil && !isNil(o.Error) {
		return true
	}

	return false
}

// SetError gets a reference to the given string and assigns it to the Error field.
func (o *WebhookResultResponse) SetError(v string) {
	o.Error = &v
}

func (o WebhookResultResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WebhookResultResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.StatusCode) {
		toSerialize["statusCode"] = o.StatusCode
	}
	if !isNil(o.Status) {
		toSerialize["status"] = o.Status
	}
	if !isNil(o.Body) {
		toSerialize["body"] = o.Body
	}
	if !isNil(o.Headers) {
		toSerialize["headers"] = o.Headers
	}
	if !isNil(o.Error) {
		toSerialize["error"] = o.Error
	}
	return toSerialize, nil
}

type NullableWebhookResultResponse struct {
	value *WebhookResultResponse
	isSet bool
}

func (v NullableWebhookResultResponse) Get() *WebhookResultResponse {
	return v.value
}

func (v *NullableWebhookResultResponse) Set(val *WebhookResultResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWebhookResultResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWebhookResultResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWebhookResultResponse(val *WebhookResultResponse) *NullableWebhookResultResponse {
	return &NullableWebhookResultResponse{value: val, isSet: true}
}

func (v NullableWebhookResultResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWebhookResultResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
