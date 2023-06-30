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
	Slug       *string `json:"slug,omitempty"`
	Weight     *int32  `json:"weight,omitempty"`
	ErrorLevel *string `json:"errorLevel,omitempty"`
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

// GetSlug returns the Slug field value if set, zero value otherwise.
func (o *LinterResourceRule) GetSlug() string {
	if o == nil || isNil(o.Slug) {
		var ret string
		return ret
	}
	return *o.Slug
}

// GetSlugOk returns a tuple with the Slug field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *LinterResourceRule) GetSlugOk() (*string, bool) {
	if o == nil || isNil(o.Slug) {
		return nil, false
	}
	return o.Slug, true
}

// HasSlug returns a boolean if a field has been set.
func (o *LinterResourceRule) HasSlug() bool {
	if o != nil && !isNil(o.Slug) {
		return true
	}

	return false
}

// SetSlug gets a reference to the given string and assigns it to the Slug field.
func (o *LinterResourceRule) SetSlug(v string) {
	o.Slug = &v
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
	if !isNil(o.Slug) {
		toSerialize["slug"] = o.Slug
	}
	if !isNil(o.Weight) {
		toSerialize["weight"] = o.Weight
	}
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
