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

// checks if the LinterResourceRule type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &LinterResourceRule{}

// LinterResourceRule struct for LinterResourceRule
type LinterResourceRule struct {
	Id               *string  `json:"id,omitempty"`
	Weight           *int32   `json:"weight,omitempty"`
	Name             *string  `json:"name,omitempty"`
	Description      *string  `json:"description,omitempty"`
	ErrorDescription *string  `json:"errorDescription,omitempty"`
	Documentation    *string  `json:"documentation,omitempty"`
	Tips             []string `json:"tips,omitempty"`
	ErrorLevel       *string  `json:"errorLevel,omitempty"`
}

// NewLinterResourceRule instantiates a new LinterResourceRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLinterResourceRule() *LinterResourceRule {
	this := LinterResourceRule{}
	return &this
}

// NewLinterResourceRuleWithDefaults instantiates a new LinterResourceRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLinterResourceRuleWithDefaults() *LinterResourceRule {
	this := LinterResourceRule{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *LinterResourceRule) GetId() string {
	if o == nil || isNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetIdOk() (*string, bool) {
	if o == nil || isNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *LinterResourceRule) HasId() bool {
	if o != nil && !isNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *LinterResourceRule) SetId(v string) {
	o.Id = &v
}

// GetWeight returns the Weight field value if set, zero value otherwise.
func (o *LinterResourceRule) GetWeight() int32 {
	if o == nil || isNil(o.Weight) {
		var ret int32
		return ret
	}
	return *o.Weight
}

// GetWeightOk returns a tuple with the Weight field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetWeightOk() (*int32, bool) {
	if o == nil || isNil(o.Weight) {
		return nil, false
	}
	return o.Weight, true
}

// HasWeight returns a boolean if a field has been set.
func (o *LinterResourceRule) HasWeight() bool {
	if o != nil && !isNil(o.Weight) {
		return true
	}

	return false
}

// SetWeight gets a reference to the given int32 and assigns it to the Weight field.
func (o *LinterResourceRule) SetWeight(v int32) {
	o.Weight = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *LinterResourceRule) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *LinterResourceRule) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *LinterResourceRule) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *LinterResourceRule) GetDescription() string {
	if o == nil || isNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetDescriptionOk() (*string, bool) {
	if o == nil || isNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *LinterResourceRule) HasDescription() bool {
	if o != nil && !isNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *LinterResourceRule) SetDescription(v string) {
	o.Description = &v
}

// GetErrorDescription returns the ErrorDescription field value if set, zero value otherwise.
func (o *LinterResourceRule) GetErrorDescription() string {
	if o == nil || isNil(o.ErrorDescription) {
		var ret string
		return ret
	}
	return *o.ErrorDescription
}

// GetErrorDescriptionOk returns a tuple with the ErrorDescription field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetErrorDescriptionOk() (*string, bool) {
	if o == nil || isNil(o.ErrorDescription) {
		return nil, false
	}
	return o.ErrorDescription, true
}

// HasErrorDescription returns a boolean if a field has been set.
func (o *LinterResourceRule) HasErrorDescription() bool {
	if o != nil && !isNil(o.ErrorDescription) {
		return true
	}

	return false
}

// SetErrorDescription gets a reference to the given string and assigns it to the ErrorDescription field.
func (o *LinterResourceRule) SetErrorDescription(v string) {
	o.ErrorDescription = &v
}

// GetDocumentation returns the Documentation field value if set, zero value otherwise.
func (o *LinterResourceRule) GetDocumentation() string {
	if o == nil || isNil(o.Documentation) {
		var ret string
		return ret
	}
	return *o.Documentation
}

// GetDocumentationOk returns a tuple with the Documentation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetDocumentationOk() (*string, bool) {
	if o == nil || isNil(o.Documentation) {
		return nil, false
	}
	return o.Documentation, true
}

// HasDocumentation returns a boolean if a field has been set.
func (o *LinterResourceRule) HasDocumentation() bool {
	if o != nil && !isNil(o.Documentation) {
		return true
	}

	return false
}

// SetDocumentation gets a reference to the given string and assigns it to the Documentation field.
func (o *LinterResourceRule) SetDocumentation(v string) {
	o.Documentation = &v
}

// GetTips returns the Tips field value if set, zero value otherwise.
func (o *LinterResourceRule) GetTips() []string {
	if o == nil || isNil(o.Tips) {
		var ret []string
		return ret
	}
	return o.Tips
}

// GetTipsOk returns a tuple with the Tips field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetTipsOk() ([]string, bool) {
	if o == nil || isNil(o.Tips) {
		return nil, false
	}
	return o.Tips, true
}

// HasTips returns a boolean if a field has been set.
func (o *LinterResourceRule) HasTips() bool {
	if o != nil && !isNil(o.Tips) {
		return true
	}

	return false
}

// SetTips gets a reference to the given []string and assigns it to the Tips field.
func (o *LinterResourceRule) SetTips(v []string) {
	o.Tips = v
}

// GetErrorLevel returns the ErrorLevel field value if set, zero value otherwise.
func (o *LinterResourceRule) GetErrorLevel() string {
	if o == nil || isNil(o.ErrorLevel) {
		var ret string
		return ret
	}
	return *o.ErrorLevel
}

// GetErrorLevelOk returns a tuple with the ErrorLevel field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetErrorLevelOk() (*string, bool) {
	if o == nil || isNil(o.ErrorLevel) {
		return nil, false
	}
	return o.ErrorLevel, true
}

// HasErrorLevel returns a boolean if a field has been set.
func (o *LinterResourceRule) HasErrorLevel() bool {
	if o != nil && !isNil(o.ErrorLevel) {
		return true
	}

	return false
}

// SetErrorLevel gets a reference to the given string and assigns it to the ErrorLevel field.
func (o *LinterResourceRule) SetErrorLevel(v string) {
	o.ErrorLevel = &v
}

func (o LinterResourceRule) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o LinterResourceRule) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !isNil(o.Weight) {
		toSerialize["weight"] = o.Weight
	}
	// skip: name is readOnly
	// skip: description is readOnly
	// skip: errorDescription is readOnly
	// skip: documentation is readOnly
	// skip: tips is readOnly
	if !isNil(o.ErrorLevel) {
		toSerialize["errorLevel"] = o.ErrorLevel
	}
	return toSerialize, nil
}

type NullableLinterResourceRule struct {
	value *LinterResourceRule
	isSet bool
}

func (v NullableLinterResourceRule) Get() *LinterResourceRule {
	return v.value
}

func (v *NullableLinterResourceRule) Set(val *LinterResourceRule) {
	v.value = val
	v.isSet = true
}

func (v NullableLinterResourceRule) IsSet() bool {
	return v.isSet
}

func (v *NullableLinterResourceRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLinterResourceRule(val *LinterResourceRule) *NullableLinterResourceRule {
	return &NullableLinterResourceRule{value: val, isSet: true}
}

func (v NullableLinterResourceRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLinterResourceRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
