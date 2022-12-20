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

// SpanSelector struct for SpanSelector
type SpanSelector struct {
	Filters       []SelectorFilter            `json:"filters"`
	PseudoClass   NullableSelectorPseudoClass `json:"pseudoClass,omitempty"`
	ChildSelector NullableSpanSelector        `json:"childSelector,omitempty"`
}

// NewSpanSelector instantiates a new SpanSelector object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSpanSelector(filters []SelectorFilter) *SpanSelector {
	this := SpanSelector{}
	this.Filters = filters
	return &this
}

// NewSpanSelectorWithDefaults instantiates a new SpanSelector object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSpanSelectorWithDefaults() *SpanSelector {
	this := SpanSelector{}
	return &this
}

// GetFilters returns the Filters field value
func (o *SpanSelector) GetFilters() []SelectorFilter {
	if o == nil {
		var ret []SelectorFilter
		return ret
	}

	return o.Filters
}

// GetFiltersOk returns a tuple with the Filters field value
// and a boolean to check if the value has been set.
func (o *SpanSelector) GetFiltersOk() ([]SelectorFilter, bool) {
	if o == nil {
		return nil, false
	}
	return o.Filters, true
}

// SetFilters sets field value
func (o *SpanSelector) SetFilters(v []SelectorFilter) {
	o.Filters = v
}

// GetPseudoClass returns the PseudoClass field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *SpanSelector) GetPseudoClass() SelectorPseudoClass {
	if o == nil || o.PseudoClass.Get() == nil {
		var ret SelectorPseudoClass
		return ret
	}
	return *o.PseudoClass.Get()
}

// GetPseudoClassOk returns a tuple with the PseudoClass field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SpanSelector) GetPseudoClassOk() (*SelectorPseudoClass, bool) {
	if o == nil {
		return nil, false
	}
	return o.PseudoClass.Get(), o.PseudoClass.IsSet()
}

// HasPseudoClass returns a boolean if a field has been set.
func (o *SpanSelector) HasPseudoClass() bool {
	if o != nil && o.PseudoClass.IsSet() {
		return true
	}

	return false
}

// SetPseudoClass gets a reference to the given NullableSelectorPseudoClass and assigns it to the PseudoClass field.
func (o *SpanSelector) SetPseudoClass(v SelectorPseudoClass) {
	o.PseudoClass.Set(&v)
}

// SetPseudoClassNil sets the value for PseudoClass to be an explicit nil
func (o *SpanSelector) SetPseudoClassNil() {
	o.PseudoClass.Set(nil)
}

// UnsetPseudoClass ensures that no value is present for PseudoClass, not even an explicit nil
func (o *SpanSelector) UnsetPseudoClass() {
	o.PseudoClass.Unset()
}

// GetChildSelector returns the ChildSelector field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *SpanSelector) GetChildSelector() SpanSelector {
	if o == nil || o.ChildSelector.Get() == nil {
		var ret SpanSelector
		return ret
	}
	return *o.ChildSelector.Get()
}

// GetChildSelectorOk returns a tuple with the ChildSelector field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SpanSelector) GetChildSelectorOk() (*SpanSelector, bool) {
	if o == nil {
		return nil, false
	}
	return o.ChildSelector.Get(), o.ChildSelector.IsSet()
}

// HasChildSelector returns a boolean if a field has been set.
func (o *SpanSelector) HasChildSelector() bool {
	if o != nil && o.ChildSelector.IsSet() {
		return true
	}

	return false
}

// SetChildSelector gets a reference to the given NullableSpanSelector and assigns it to the ChildSelector field.
func (o *SpanSelector) SetChildSelector(v SpanSelector) {
	o.ChildSelector.Set(&v)
}

// SetChildSelectorNil sets the value for ChildSelector to be an explicit nil
func (o *SpanSelector) SetChildSelectorNil() {
	o.ChildSelector.Set(nil)
}

// UnsetChildSelector ensures that no value is present for ChildSelector, not even an explicit nil
func (o *SpanSelector) UnsetChildSelector() {
	o.ChildSelector.Unset()
}

func (o SpanSelector) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["filters"] = o.Filters
	}
	if o.PseudoClass.IsSet() {
		toSerialize["pseudoClass"] = o.PseudoClass.Get()
	}
	if o.ChildSelector.IsSet() {
		toSerialize["childSelector"] = o.ChildSelector.Get()
	}
	return json.Marshal(toSerialize)
}

type NullableSpanSelector struct {
	value *SpanSelector
	isSet bool
}

func (v NullableSpanSelector) Get() *SpanSelector {
	return v.value
}

func (v *NullableSpanSelector) Set(val *SpanSelector) {
	v.value = val
	v.isSet = true
}

func (v NullableSpanSelector) IsSet() bool {
	return v.isSet
}

func (v *NullableSpanSelector) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSpanSelector(val *SpanSelector) *NullableSpanSelector {
	return &NullableSpanSelector{value: val, isSet: true}
}

func (v NullableSpanSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSpanSelector) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
