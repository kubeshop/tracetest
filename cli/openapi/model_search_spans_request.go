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

// checks if the SearchSpansRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SearchSpansRequest{}

// SearchSpansRequest struct for SearchSpansRequest
type SearchSpansRequest struct {
	// query to filter spans, can be either a full text search or a Span Query Language query
	Query *string `json:"query,omitempty"`
}

// NewSearchSpansRequest instantiates a new SearchSpansRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSearchSpansRequest() *SearchSpansRequest {
	this := SearchSpansRequest{}
	return &this
}

// NewSearchSpansRequestWithDefaults instantiates a new SearchSpansRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSearchSpansRequestWithDefaults() *SearchSpansRequest {
	this := SearchSpansRequest{}
	return &this
}

// GetQuery returns the Query field value if set, zero value otherwise.
func (o *SearchSpansRequest) GetQuery() string {
	if o == nil || isNil(o.Query) {
		var ret string
		return ret
	}
	return *o.Query
}

// GetQueryOk returns a tuple with the Query field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SearchSpansRequest) GetQueryOk() (*string, bool) {
	if o == nil || isNil(o.Query) {
		return nil, false
	}
	return o.Query, true
}

// HasQuery returns a boolean if a field has been set.
func (o *SearchSpansRequest) HasQuery() bool {
	if o != nil && !isNil(o.Query) {
		return true
	}

	return false
}

// SetQuery gets a reference to the given string and assigns it to the Query field.
func (o *SearchSpansRequest) SetQuery(v string) {
	o.Query = &v
}

func (o SearchSpansRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SearchSpansRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.Query) {
		toSerialize["query"] = o.Query
	}
	return toSerialize, nil
}

type NullableSearchSpansRequest struct {
	value *SearchSpansRequest
	isSet bool
}

func (v NullableSearchSpansRequest) Get() *SearchSpansRequest {
	return v.value
}

func (v *NullableSearchSpansRequest) Set(val *SearchSpansRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableSearchSpansRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableSearchSpansRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSearchSpansRequest(val *SearchSpansRequest) *NullableSearchSpansRequest {
	return &NullableSearchSpansRequest{value: val, isSet: true}
}

func (v NullableSearchSpansRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSearchSpansRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
