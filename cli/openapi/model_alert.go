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

// checks if the Alert type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Alert{}

// Alert struct for Alert
type Alert struct {
	Id      *string  `json:"id,omitempty"`
	Type    *string  `json:"type,omitempty"`
	Webhook *Webhook `json:"webhook,omitempty"`
	Events  []string `json:"events,omitempty"`
}

// NewAlert instantiates a new Alert object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAlert() *Alert {
	this := Alert{}
	return &this
}

// NewAlertWithDefaults instantiates a new Alert object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAlertWithDefaults() *Alert {
	this := Alert{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Alert) GetId() string {
	if o == nil || isNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Alert) GetIdOk() (*string, bool) {
	if o == nil || isNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Alert) HasId() bool {
	if o != nil && !isNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Alert) SetId(v string) {
	o.Id = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *Alert) GetType() string {
	if o == nil || isNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Alert) GetTypeOk() (*string, bool) {
	if o == nil || isNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *Alert) HasType() bool {
	if o != nil && !isNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *Alert) SetType(v string) {
	o.Type = &v
}

// GetWebhook returns the Webhook field value if set, zero value otherwise.
func (o *Alert) GetWebhook() Webhook {
	if o == nil || isNil(o.Webhook) {
		var ret Webhook
		return ret
	}
	return *o.Webhook
}

// GetWebhookOk returns a tuple with the Webhook field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Alert) GetWebhookOk() (*Webhook, bool) {
	if o == nil || isNil(o.Webhook) {
		return nil, false
	}
	return o.Webhook, true
}

// HasWebhook returns a boolean if a field has been set.
func (o *Alert) HasWebhook() bool {
	if o != nil && !isNil(o.Webhook) {
		return true
	}

	return false
}

// SetWebhook gets a reference to the given Webhook and assigns it to the Webhook field.
func (o *Alert) SetWebhook(v Webhook) {
	o.Webhook = &v
}

// GetEvents returns the Events field value if set, zero value otherwise.
func (o *Alert) GetEvents() []string {
	if o == nil || isNil(o.Events) {
		var ret []string
		return ret
	}
	return o.Events
}

// GetEventsOk returns a tuple with the Events field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Alert) GetEventsOk() ([]string, bool) {
	if o == nil || isNil(o.Events) {
		return nil, false
	}
	return o.Events, true
}

// HasEvents returns a boolean if a field has been set.
func (o *Alert) HasEvents() bool {
	if o != nil && !isNil(o.Events) {
		return true
	}

	return false
}

// SetEvents gets a reference to the given []string and assigns it to the Events field.
func (o *Alert) SetEvents(v []string) {
	o.Events = v
}

func (o Alert) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Alert) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !isNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !isNil(o.Webhook) {
		toSerialize["webhook"] = o.Webhook
	}
	if !isNil(o.Events) {
		toSerialize["events"] = o.Events
	}
	return toSerialize, nil
}

type NullableAlert struct {
	value *Alert
	isSet bool
}

func (v NullableAlert) Get() *Alert {
	return v.value
}

func (v *NullableAlert) Set(val *Alert) {
	v.value = val
	v.isSet = true
}

func (v NullableAlert) IsSet() bool {
	return v.isSet
}

func (v *NullableAlert) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAlert(val *Alert) *NullableAlert {
	return &NullableAlert{value: val, isSet: true}
}

func (v NullableAlert) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAlert) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
