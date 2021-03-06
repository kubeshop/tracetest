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

// Trigger struct for Trigger
type Trigger struct {
	TriggerType     *string                 `json:"triggerType,omitempty"`
	TriggerSettings *TriggerTriggerSettings `json:"triggerSettings,omitempty"`
}

// NewTrigger instantiates a new Trigger object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTrigger() *Trigger {
	this := Trigger{}
	return &this
}

// NewTriggerWithDefaults instantiates a new Trigger object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTriggerWithDefaults() *Trigger {
	this := Trigger{}
	return &this
}

// GetTriggerType returns the TriggerType field value if set, zero value otherwise.
func (o *Trigger) GetTriggerType() string {
	if o == nil || o.TriggerType == nil {
		var ret string
		return ret
	}
	return *o.TriggerType
}

// GetTriggerTypeOk returns a tuple with the TriggerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Trigger) GetTriggerTypeOk() (*string, bool) {
	if o == nil || o.TriggerType == nil {
		return nil, false
	}
	return o.TriggerType, true
}

// HasTriggerType returns a boolean if a field has been set.
func (o *Trigger) HasTriggerType() bool {
	if o != nil && o.TriggerType != nil {
		return true
	}

	return false
}

// SetTriggerType gets a reference to the given string and assigns it to the TriggerType field.
func (o *Trigger) SetTriggerType(v string) {
	o.TriggerType = &v
}

// GetTriggerSettings returns the TriggerSettings field value if set, zero value otherwise.
func (o *Trigger) GetTriggerSettings() TriggerTriggerSettings {
	if o == nil || o.TriggerSettings == nil {
		var ret TriggerTriggerSettings
		return ret
	}
	return *o.TriggerSettings
}

// GetTriggerSettingsOk returns a tuple with the TriggerSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Trigger) GetTriggerSettingsOk() (*TriggerTriggerSettings, bool) {
	if o == nil || o.TriggerSettings == nil {
		return nil, false
	}
	return o.TriggerSettings, true
}

// HasTriggerSettings returns a boolean if a field has been set.
func (o *Trigger) HasTriggerSettings() bool {
	if o != nil && o.TriggerSettings != nil {
		return true
	}

	return false
}

// SetTriggerSettings gets a reference to the given TriggerTriggerSettings and assigns it to the TriggerSettings field.
func (o *Trigger) SetTriggerSettings(v TriggerTriggerSettings) {
	o.TriggerSettings = &v
}

func (o Trigger) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.TriggerType != nil {
		toSerialize["triggerType"] = o.TriggerType
	}
	if o.TriggerSettings != nil {
		toSerialize["triggerSettings"] = o.TriggerSettings
	}
	return json.Marshal(toSerialize)
}

type NullableTrigger struct {
	value *Trigger
	isSet bool
}

func (v NullableTrigger) Get() *Trigger {
	return v.value
}

func (v *NullableTrigger) Set(val *Trigger) {
	v.value = val
	v.isSet = true
}

func (v NullableTrigger) IsSet() bool {
	return v.isSet
}

func (v *NullableTrigger) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTrigger(val *Trigger) *NullableTrigger {
	return &NullableTrigger{value: val, isSet: true}
}

func (v NullableTrigger) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTrigger) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
