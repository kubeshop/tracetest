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

// checks if the Version type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Version{}

// Version struct for Version
type Version struct {
	Version       *string `json:"version,omitempty"`
	Type          *string `json:"type,omitempty"`
	UiEndpoint    *string `json:"uiEndpoint,omitempty"`
	AgentEndpoint *string `json:"agentEndpoint,omitempty"`
}

// NewVersion instantiates a new Version object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewVersion() *Version {
	this := Version{}
	return &this
}

// NewVersionWithDefaults instantiates a new Version object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewVersionWithDefaults() *Version {
	this := Version{}
	return &this
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *Version) GetVersion() string {
	if o == nil || isNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Version) GetVersionOk() (*string, bool) {
	if o == nil || isNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *Version) HasVersion() bool {
	if o != nil && !isNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *Version) SetVersion(v string) {
	o.Version = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *Version) GetType() string {
	if o == nil || isNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Version) GetTypeOk() (*string, bool) {
	if o == nil || isNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *Version) HasType() bool {
	if o != nil && !isNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *Version) SetType(v string) {
	o.Type = &v
}

// GetUiEndpoint returns the UiEndpoint field value if set, zero value otherwise.
func (o *Version) GetUiEndpoint() string {
	if o == nil || isNil(o.UiEndpoint) {
		var ret string
		return ret
	}
	return *o.UiEndpoint
}

// GetUiEndpointOk returns a tuple with the UiEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Version) GetUiEndpointOk() (*string, bool) {
	if o == nil || isNil(o.UiEndpoint) {
		return nil, false
	}
	return o.UiEndpoint, true
}

// HasUiEndpoint returns a boolean if a field has been set.
func (o *Version) HasUiEndpoint() bool {
	if o != nil && !isNil(o.UiEndpoint) {
		return true
	}

	return false
}

// SetUiEndpoint gets a reference to the given string and assigns it to the UiEndpoint field.
func (o *Version) SetUiEndpoint(v string) {
	o.UiEndpoint = &v
}

// GetAgentEndpoint returns the AgentEndpoint field value if set, zero value otherwise.
func (o *Version) GetAgentEndpoint() string {
	if o == nil || isNil(o.AgentEndpoint) {
		var ret string
		return ret
	}
	return *o.AgentEndpoint
}

// GetAgentEndpointOk returns a tuple with the AgentEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Version) GetAgentEndpointOk() (*string, bool) {
	if o == nil || isNil(o.AgentEndpoint) {
		return nil, false
	}
	return o.AgentEndpoint, true
}

// HasAgentEndpoint returns a boolean if a field has been set.
func (o *Version) HasAgentEndpoint() bool {
	if o != nil && !isNil(o.AgentEndpoint) {
		return true
	}

	return false
}

// SetAgentEndpoint gets a reference to the given string and assigns it to the AgentEndpoint field.
func (o *Version) SetAgentEndpoint(v string) {
	o.AgentEndpoint = &v
}

func (o Version) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Version) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !isNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !isNil(o.UiEndpoint) {
		toSerialize["uiEndpoint"] = o.UiEndpoint
	}
	if !isNil(o.AgentEndpoint) {
		toSerialize["agentEndpoint"] = o.AgentEndpoint
	}
	return toSerialize, nil
}

type NullableVersion struct {
	value *Version
	isSet bool
}

func (v NullableVersion) Get() *Version {
	return v.value
}

func (v *NullableVersion) Set(val *Version) {
	v.value = val
	v.isSet = true
}

func (v NullableVersion) IsSet() bool {
	return v.isSet
}

func (v *NullableVersion) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableVersion(val *Version) *NullableVersion {
	return &NullableVersion{value: val, isSet: true}
}

func (v NullableVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableVersion) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
