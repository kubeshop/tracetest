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

// GRPCClientSettings struct for GRPCClientSettings
type GRPCClientSettings struct {
	Endpoint        *string      `json:"endpoint,omitempty"`
	ReadBufferSize  *float32     `json:"readBufferSize,omitempty"`
	WriteBufferSize *float32     `json:"writeBufferSize,omitempty"`
	WaitForReady    *bool        `json:"waitForReady,omitempty"`
	Headers         []HTTPHeader `json:"headers,omitempty"`
	BalancerName    *string      `json:"balancerName,omitempty"`
	Compression     *string      `json:"compression,omitempty"`
	Tls             *TLS         `json:"tls,omitempty"`
	Auth            *HTTPAuth    `json:"auth,omitempty"`
}

// NewGRPCClientSettings instantiates a new GRPCClientSettings object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGRPCClientSettings() *GRPCClientSettings {
	this := GRPCClientSettings{}
	return &this
}

// NewGRPCClientSettingsWithDefaults instantiates a new GRPCClientSettings object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGRPCClientSettingsWithDefaults() *GRPCClientSettings {
	this := GRPCClientSettings{}
	return &this
}

// GetEndpoint returns the Endpoint field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetEndpoint() string {
	if o == nil || o.Endpoint == nil {
		var ret string
		return ret
	}
	return *o.Endpoint
}

// GetEndpointOk returns a tuple with the Endpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetEndpointOk() (*string, bool) {
	if o == nil || o.Endpoint == nil {
		return nil, false
	}
	return o.Endpoint, true
}

// HasEndpoint returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasEndpoint() bool {
	if o != nil && o.Endpoint != nil {
		return true
	}

	return false
}

// SetEndpoint gets a reference to the given string and assigns it to the Endpoint field.
func (o *GRPCClientSettings) SetEndpoint(v string) {
	o.Endpoint = &v
}

// GetReadBufferSize returns the ReadBufferSize field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetReadBufferSize() float32 {
	if o == nil || o.ReadBufferSize == nil {
		var ret float32
		return ret
	}
	return *o.ReadBufferSize
}

// GetReadBufferSizeOk returns a tuple with the ReadBufferSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetReadBufferSizeOk() (*float32, bool) {
	if o == nil || o.ReadBufferSize == nil {
		return nil, false
	}
	return o.ReadBufferSize, true
}

// HasReadBufferSize returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasReadBufferSize() bool {
	if o != nil && o.ReadBufferSize != nil {
		return true
	}

	return false
}

// SetReadBufferSize gets a reference to the given float32 and assigns it to the ReadBufferSize field.
func (o *GRPCClientSettings) SetReadBufferSize(v float32) {
	o.ReadBufferSize = &v
}

// GetWriteBufferSize returns the WriteBufferSize field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetWriteBufferSize() float32 {
	if o == nil || o.WriteBufferSize == nil {
		var ret float32
		return ret
	}
	return *o.WriteBufferSize
}

// GetWriteBufferSizeOk returns a tuple with the WriteBufferSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetWriteBufferSizeOk() (*float32, bool) {
	if o == nil || o.WriteBufferSize == nil {
		return nil, false
	}
	return o.WriteBufferSize, true
}

// HasWriteBufferSize returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasWriteBufferSize() bool {
	if o != nil && o.WriteBufferSize != nil {
		return true
	}

	return false
}

// SetWriteBufferSize gets a reference to the given float32 and assigns it to the WriteBufferSize field.
func (o *GRPCClientSettings) SetWriteBufferSize(v float32) {
	o.WriteBufferSize = &v
}

// GetWaitForReady returns the WaitForReady field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetWaitForReady() bool {
	if o == nil || o.WaitForReady == nil {
		var ret bool
		return ret
	}
	return *o.WaitForReady
}

// GetWaitForReadyOk returns a tuple with the WaitForReady field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetWaitForReadyOk() (*bool, bool) {
	if o == nil || o.WaitForReady == nil {
		return nil, false
	}
	return o.WaitForReady, true
}

// HasWaitForReady returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasWaitForReady() bool {
	if o != nil && o.WaitForReady != nil {
		return true
	}

	return false
}

// SetWaitForReady gets a reference to the given bool and assigns it to the WaitForReady field.
func (o *GRPCClientSettings) SetWaitForReady(v bool) {
	o.WaitForReady = &v
}

// GetHeaders returns the Headers field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetHeaders() []HTTPHeader {
	if o == nil || o.Headers == nil {
		var ret []HTTPHeader
		return ret
	}
	return o.Headers
}

// GetHeadersOk returns a tuple with the Headers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetHeadersOk() ([]HTTPHeader, bool) {
	if o == nil || o.Headers == nil {
		return nil, false
	}
	return o.Headers, true
}

// HasHeaders returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasHeaders() bool {
	if o != nil && o.Headers != nil {
		return true
	}

	return false
}

// SetHeaders gets a reference to the given []HTTPHeader and assigns it to the Headers field.
func (o *GRPCClientSettings) SetHeaders(v []HTTPHeader) {
	o.Headers = v
}

// GetBalancerName returns the BalancerName field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetBalancerName() string {
	if o == nil || o.BalancerName == nil {
		var ret string
		return ret
	}
	return *o.BalancerName
}

// GetBalancerNameOk returns a tuple with the BalancerName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetBalancerNameOk() (*string, bool) {
	if o == nil || o.BalancerName == nil {
		return nil, false
	}
	return o.BalancerName, true
}

// HasBalancerName returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasBalancerName() bool {
	if o != nil && o.BalancerName != nil {
		return true
	}

	return false
}

// SetBalancerName gets a reference to the given string and assigns it to the BalancerName field.
func (o *GRPCClientSettings) SetBalancerName(v string) {
	o.BalancerName = &v
}

// GetCompression returns the Compression field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetCompression() string {
	if o == nil || o.Compression == nil {
		var ret string
		return ret
	}
	return *o.Compression
}

// GetCompressionOk returns a tuple with the Compression field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetCompressionOk() (*string, bool) {
	if o == nil || o.Compression == nil {
		return nil, false
	}
	return o.Compression, true
}

// HasCompression returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasCompression() bool {
	if o != nil && o.Compression != nil {
		return true
	}

	return false
}

// SetCompression gets a reference to the given string and assigns it to the Compression field.
func (o *GRPCClientSettings) SetCompression(v string) {
	o.Compression = &v
}

// GetTls returns the Tls field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetTls() TLS {
	if o == nil || o.Tls == nil {
		var ret TLS
		return ret
	}
	return *o.Tls
}

// GetTlsOk returns a tuple with the Tls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetTlsOk() (*TLS, bool) {
	if o == nil || o.Tls == nil {
		return nil, false
	}
	return o.Tls, true
}

// HasTls returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasTls() bool {
	if o != nil && o.Tls != nil {
		return true
	}

	return false
}

// SetTls gets a reference to the given TLS and assigns it to the Tls field.
func (o *GRPCClientSettings) SetTls(v TLS) {
	o.Tls = &v
}

// GetAuth returns the Auth field value if set, zero value otherwise.
func (o *GRPCClientSettings) GetAuth() HTTPAuth {
	if o == nil || o.Auth == nil {
		var ret HTTPAuth
		return ret
	}
	return *o.Auth
}

// GetAuthOk returns a tuple with the Auth field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GRPCClientSettings) GetAuthOk() (*HTTPAuth, bool) {
	if o == nil || o.Auth == nil {
		return nil, false
	}
	return o.Auth, true
}

// HasAuth returns a boolean if a field has been set.
func (o *GRPCClientSettings) HasAuth() bool {
	if o != nil && o.Auth != nil {
		return true
	}

	return false
}

// SetAuth gets a reference to the given HTTPAuth and assigns it to the Auth field.
func (o *GRPCClientSettings) SetAuth(v HTTPAuth) {
	o.Auth = &v
}

func (o GRPCClientSettings) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Endpoint != nil {
		toSerialize["endpoint"] = o.Endpoint
	}
	if o.ReadBufferSize != nil {
		toSerialize["readBufferSize"] = o.ReadBufferSize
	}
	if o.WriteBufferSize != nil {
		toSerialize["writeBufferSize"] = o.WriteBufferSize
	}
	if o.WaitForReady != nil {
		toSerialize["waitForReady"] = o.WaitForReady
	}
	if o.Headers != nil {
		toSerialize["headers"] = o.Headers
	}
	if o.BalancerName != nil {
		toSerialize["balancerName"] = o.BalancerName
	}
	if o.Compression != nil {
		toSerialize["compression"] = o.Compression
	}
	if o.Tls != nil {
		toSerialize["tls"] = o.Tls
	}
	if o.Auth != nil {
		toSerialize["auth"] = o.Auth
	}
	return json.Marshal(toSerialize)
}

type NullableGRPCClientSettings struct {
	value *GRPCClientSettings
	isSet bool
}

func (v NullableGRPCClientSettings) Get() *GRPCClientSettings {
	return v.value
}

func (v *NullableGRPCClientSettings) Set(val *GRPCClientSettings) {
	v.value = val
	v.isSet = true
}

func (v NullableGRPCClientSettings) IsSet() bool {
	return v.isSet
}

func (v *NullableGRPCClientSettings) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGRPCClientSettings(val *GRPCClientSettings) *NullableGRPCClientSettings {
	return &NullableGRPCClientSettings{value: val, isSet: true}
}

func (v NullableGRPCClientSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGRPCClientSettings) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
