/*
TraceTest

OpenAPI definition for TraceTest endpoint and resources

API version: 0.2.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"time"
)

// checks if the TestSummaryLastRun type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TestSummaryLastRun{}

// TestSummaryLastRun struct for TestSummaryLastRun
type TestSummaryLastRun struct {
	Time   NullableTime `json:"time,omitempty"`
	Passes *int32       `json:"passes,omitempty"`
	Fails  *int32       `json:"fails,omitempty"`
}

// NewTestSummaryLastRun instantiates a new TestSummaryLastRun object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTestSummaryLastRun() *TestSummaryLastRun {
	this := TestSummaryLastRun{}
	return &this
}

// NewTestSummaryLastRunWithDefaults instantiates a new TestSummaryLastRun object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTestSummaryLastRunWithDefaults() *TestSummaryLastRun {
	this := TestSummaryLastRun{}
	return &this
}

// GetTime returns the Time field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TestSummaryLastRun) GetTime() time.Time {
	if o == nil || isNil(o.Time.Get()) {
		var ret time.Time
		return ret
	}
	return *o.Time.Get()
}

// GetTimeOk returns a tuple with the Time field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TestSummaryLastRun) GetTimeOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return o.Time.Get(), o.Time.IsSet()
}

// HasTime returns a boolean if a field has been set.
func (o *TestSummaryLastRun) HasTime() bool {
	if o != nil && o.Time.IsSet() {
		return true
	}

	return false
}

// SetTime gets a reference to the given NullableTime and assigns it to the Time field.
func (o *TestSummaryLastRun) SetTime(v time.Time) {
	o.Time.Set(&v)
}

// SetTimeNil sets the value for Time to be an explicit nil
func (o *TestSummaryLastRun) SetTimeNil() {
	o.Time.Set(nil)
}

// UnsetTime ensures that no value is present for Time, not even an explicit nil
func (o *TestSummaryLastRun) UnsetTime() {
	o.Time.Unset()
}

// GetPasses returns the Passes field value if set, zero value otherwise.
func (o *TestSummaryLastRun) GetPasses() int32 {
	if o == nil || isNil(o.Passes) {
		var ret int32
		return ret
	}
	return *o.Passes
}

// GetPassesOk returns a tuple with the Passes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestSummaryLastRun) GetPassesOk() (*int32, bool) {
	if o == nil || isNil(o.Passes) {
		return nil, false
	}
	return o.Passes, true
}

// HasPasses returns a boolean if a field has been set.
func (o *TestSummaryLastRun) HasPasses() bool {
	if o != nil && !isNil(o.Passes) {
		return true
	}

	return false
}

// SetPasses gets a reference to the given int32 and assigns it to the Passes field.
func (o *TestSummaryLastRun) SetPasses(v int32) {
	o.Passes = &v
}

// GetFails returns the Fails field value if set, zero value otherwise.
func (o *TestSummaryLastRun) GetFails() int32 {
	if o == nil || isNil(o.Fails) {
		var ret int32
		return ret
	}
	return *o.Fails
}

// GetFailsOk returns a tuple with the Fails field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestSummaryLastRun) GetFailsOk() (*int32, bool) {
	if o == nil || isNil(o.Fails) {
		return nil, false
	}
	return o.Fails, true
}

// HasFails returns a boolean if a field has been set.
func (o *TestSummaryLastRun) HasFails() bool {
	if o != nil && !isNil(o.Fails) {
		return true
	}

	return false
}

// SetFails gets a reference to the given int32 and assigns it to the Fails field.
func (o *TestSummaryLastRun) SetFails(v int32) {
	o.Fails = &v
}

func (o TestSummaryLastRun) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TestSummaryLastRun) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Time.IsSet() {
		toSerialize["time"] = o.Time.Get()
	}
	// skip: passes is readOnly
	// skip: fails is readOnly
	return toSerialize, nil
}

type NullableTestSummaryLastRun struct {
	value *TestSummaryLastRun
	isSet bool
}

func (v NullableTestSummaryLastRun) Get() *TestSummaryLastRun {
	return v.value
}

func (v *NullableTestSummaryLastRun) Set(val *TestSummaryLastRun) {
	v.value = val
	v.isSet = true
}

func (v NullableTestSummaryLastRun) IsSet() bool {
	return v.isSet
}

func (v *NullableTestSummaryLastRun) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTestSummaryLastRun(val *TestSummaryLastRun) *NullableTestSummaryLastRun {
	return &NullableTestSummaryLastRun{value: val, isSet: true}
}

func (v NullableTestSummaryLastRun) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTestSummaryLastRun) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
