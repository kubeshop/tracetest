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

// checks if the LinternResultPluginRule type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LinternResultPluginRule{}

// LinternResultPluginRule struct for LinternResultPluginRule
type LinternResultPluginRule struct {
	Name        *string                         `json:"name,omitempty"`
	Description *string                         `json:"description,omitempty"`
	Passed      *bool                           `json:"passed,omitempty"`
	Weight      *int32                          `json:"weight,omitempty"`
	Tips        *string                         `json:"tips,omitempty"`
	Results     []LinternResultPluginRuleResult `json:"results,omitempty"`
}

// NewLinternResultPluginRule instantiates a new LinternResultPluginRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinternResultPluginRule() *LinternResultPluginRule {
	this := LinternResultPluginRule{}
	return &this
}

// NewLinternResultPluginRuleWithDefaults instantiates a new LinternResultPluginRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinternResultPluginRuleWithDefaults() *LinternResultPluginRule {
	this := LinternResultPluginRule{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *LinternResultPluginRule) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinternResultPluginRule) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *LinternResultPluginRule) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *LinternResultPluginRule) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *LinternResultPluginRule) GetDescription() string {
	if o == nil || isNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinternResultPluginRule) GetDescriptionOk() (*string, bool) {
	if o == nil || isNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *LinternResultPluginRule) HasDescription() bool {
	if o != nil && !isNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *LinternResultPluginRule) SetDescription(v string) {
	o.Description = &v
}

// GetPassed returns the Passed field value if set, zero value otherwise.
func (o *LinternResultPluginRule) GetPassed() bool {
	if o == nil || isNil(o.Passed) {
		var ret bool
		return ret
	}
	return *o.Passed
}

// GetPassedOk returns a tuple with the Passed field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinternResultPluginRule) GetPassedOk() (*bool, bool) {
	if o == nil || isNil(o.Passed) {
		return nil, false
	}
	return o.Passed, true
}

// HasPassed returns a boolean if a field has been set.
func (o *LinternResultPluginRule) HasPassed() bool {
	if o != nil && !isNil(o.Passed) {
		return true
	}

	return false
}

// SetPassed gets a reference to the given bool and assigns it to the Passed field.
func (o *LinternResultPluginRule) SetPassed(v bool) {
	o.Passed = &v
}

// GetWeight returns the Weight field value if set, zero value otherwise.
func (o *LinternResultPluginRule) GetWeight() int32 {
	if o == nil || isNil(o.Weight) {
		var ret int32
		return ret
	}
	return *o.Weight
}

// GetWeightOk returns a tuple with the Weight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinternResultPluginRule) GetWeightOk() (*int32, bool) {
	if o == nil || isNil(o.Weight) {
		return nil, false
	}
	return o.Weight, true
}

// HasWeight returns a boolean if a field has been set.
func (o *LinternResultPluginRule) HasWeight() bool {
	if o != nil && !isNil(o.Weight) {
		return true
	}

	return false
}

// SetWeight gets a reference to the given int32 and assigns it to the Weight field.
func (o *LinternResultPluginRule) SetWeight(v int32) {
	o.Weight = &v
}

// GetTips returns the Tips field value if set, zero value otherwise.
func (o *LinternResultPluginRule) GetTips() string {
	if o == nil || isNil(o.Tips) {
		var ret string
		return ret
	}
	return *o.Tips
}

// GetTipsOk returns a tuple with the Tips field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinternResultPluginRule) GetTipsOk() (*string, bool) {
	if o == nil || isNil(o.Tips) {
		return nil, false
	}
	return o.Tips, true
}

// HasTips returns a boolean if a field has been set.
func (o *LinternResultPluginRule) HasTips() bool {
	if o != nil && !isNil(o.Tips) {
		return true
	}

	return false
}

// SetTips gets a reference to the given string and assigns it to the Tips field.
func (o *LinternResultPluginRule) SetTips(v string) {
	o.Tips = &v
}

// GetResults returns the Results field value if set, zero value otherwise.
func (o *LinternResultPluginRule) GetResults() []LinternResultPluginRuleResult {
	if o == nil || isNil(o.Results) {
		var ret []LinternResultPluginRuleResult
		return ret
	}
	return o.Results
}

// GetResultsOk returns a tuple with the Results field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinternResultPluginRule) GetResultsOk() ([]LinternResultPluginRuleResult, bool) {
	if o == nil || isNil(o.Results) {
		return nil, false
	}
	return o.Results, true
}

// HasResults returns a boolean if a field has been set.
func (o *LinternResultPluginRule) HasResults() bool {
	if o != nil && !isNil(o.Results) {
		return true
	}

	return false
}

// SetResults gets a reference to the given []LinternResultPluginRuleResult and assigns it to the Results field.
func (o *LinternResultPluginRule) SetResults(v []LinternResultPluginRuleResult) {
	o.Results = v
}

func (o LinternResultPluginRule) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LinternResultPluginRule) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !isNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !isNil(o.Passed) {
		toSerialize["passed"] = o.Passed
	}
	if !isNil(o.Weight) {
		toSerialize["weight"] = o.Weight
	}
	if !isNil(o.Tips) {
		toSerialize["tips"] = o.Tips
	}
	if !isNil(o.Results) {
		toSerialize["results"] = o.Results
	}
	return toSerialize, nil
}

type NullableLinternResultPluginRule struct {
	value *LinternResultPluginRule
	isSet bool
}

func (v NullableLinternResultPluginRule) Get() *LinternResultPluginRule {
	return v.value
}

func (v *NullableLinternResultPluginRule) Set(val *LinternResultPluginRule) {
	v.value = val
	v.isSet = true
}

func (v NullableLinternResultPluginRule) IsSet() bool {
	return v.isSet
}

func (v *NullableLinternResultPluginRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinternResultPluginRule(val *LinternResultPluginRule) *NullableLinternResultPluginRule {
	return &NullableLinternResultPluginRule{value: val, isSet: true}
}

func (v NullableLinternResultPluginRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinternResultPluginRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
