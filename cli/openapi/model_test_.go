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

// Test struct for Test
type Test struct {
	Id          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	// version number of the test
	Version          *int32                `json:"version,omitempty"`
	ServiceUnderTest *TestServiceUnderTest `json:"serviceUnderTest,omitempty"`
	Definition       *TestDefinition       `json:"definition,omitempty"`
}

// NewTest instantiates a new Test object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTest() *Test {
	this := Test{}
	return &this
}

// NewTestWithDefaults instantiates a new Test object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTestWithDefaults() *Test {
	this := Test{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Test) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Test) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Test) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Test) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Test) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Test) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Test) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Test) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Test) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Test) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Test) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Test) SetDescription(v string) {
	o.Description = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *Test) GetVersion() int32 {
	if o == nil || o.Version == nil {
		var ret int32
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Test) GetVersionOk() (*int32, bool) {
	if o == nil || o.Version == nil {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *Test) HasVersion() bool {
	if o != nil && o.Version != nil {
		return true
	}

	return false
}

// SetVersion gets a reference to the given int32 and assigns it to the Version field.
func (o *Test) SetVersion(v int32) {
	o.Version = &v
}

// GetServiceUnderTest returns the ServiceUnderTest field value if set, zero value otherwise.
func (o *Test) GetServiceUnderTest() TestServiceUnderTest {
	if o == nil || o.ServiceUnderTest == nil {
		var ret TestServiceUnderTest
		return ret
	}
	return *o.ServiceUnderTest
}

// GetServiceUnderTestOk returns a tuple with the ServiceUnderTest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Test) GetServiceUnderTestOk() (*TestServiceUnderTest, bool) {
	if o == nil || o.ServiceUnderTest == nil {
		return nil, false
	}
	return o.ServiceUnderTest, true
}

// HasServiceUnderTest returns a boolean if a field has been set.
func (o *Test) HasServiceUnderTest() bool {
	if o != nil && o.ServiceUnderTest != nil {
		return true
	}

	return false
}

// SetServiceUnderTest gets a reference to the given TestServiceUnderTest and assigns it to the ServiceUnderTest field.
func (o *Test) SetServiceUnderTest(v TestServiceUnderTest) {
	o.ServiceUnderTest = &v
}

// GetDefinition returns the Definition field value if set, zero value otherwise.
func (o *Test) GetDefinition() TestDefinition {
	if o == nil || o.Definition == nil {
		var ret TestDefinition
		return ret
	}
	return *o.Definition
}

// GetDefinitionOk returns a tuple with the Definition field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Test) GetDefinitionOk() (*TestDefinition, bool) {
	if o == nil || o.Definition == nil {
		return nil, false
	}
	return o.Definition, true
}

// HasDefinition returns a boolean if a field has been set.
func (o *Test) HasDefinition() bool {
	if o != nil && o.Definition != nil {
		return true
	}

	return false
}

// SetDefinition gets a reference to the given TestDefinition and assigns it to the Definition field.
func (o *Test) SetDefinition(v TestDefinition) {
	o.Definition = &v
}

func (o Test) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Version != nil {
		toSerialize["version"] = o.Version
	}
	if o.ServiceUnderTest != nil {
		toSerialize["serviceUnderTest"] = o.ServiceUnderTest
	}
	if o.Definition != nil {
		toSerialize["definition"] = o.Definition
	}
	return json.Marshal(toSerialize)
}

type NullableTest struct {
	value *Test
	isSet bool
}

func (v NullableTest) Get() *Test {
	return v.value
}

func (v *NullableTest) Set(val *Test) {
	v.value = val
	v.isSet = true
}

func (v NullableTest) IsSet() bool {
	return v.isSet
}

func (v *NullableTest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTest(val *Test) *NullableTest {
	return &NullableTest{value: val, isSet: true}
}

func (v NullableTest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
