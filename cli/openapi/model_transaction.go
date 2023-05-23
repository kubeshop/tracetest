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

// checks if the Transaction type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Transaction{}

// Transaction struct for Transaction
type Transaction struct {
	Id          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	// version number of the test
	Version *int32 `json:"version,omitempty"`
	// list of steps of the transaction containing just each test id
	Steps []string `json:"steps,omitempty"`
	// list of steps of the transaction containing the whole test object
	FullSteps []Test       `json:"fullSteps,omitempty"`
	CreatedAt *time.Time   `json:"createdAt,omitempty"`
	Summary   *TestSummary `json:"summary,omitempty"`
}

// NewTransaction instantiates a new Transaction object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTransaction() *Transaction {
	this := Transaction{}
	return &this
}

// NewTransactionWithDefaults instantiates a new Transaction object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTransactionWithDefaults() *Transaction {
	this := Transaction{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Transaction) GetId() string {
	if o == nil || isNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetIdOk() (*string, bool) {
	if o == nil || isNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Transaction) HasId() bool {
	if o != nil && !isNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Transaction) SetId(v string) {
	o.Id = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Transaction) GetName() string {
	if o == nil || isNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetNameOk() (*string, bool) {
	if o == nil || isNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Transaction) HasName() bool {
	if o != nil && !isNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Transaction) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Transaction) GetDescription() string {
	if o == nil || isNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetDescriptionOk() (*string, bool) {
	if o == nil || isNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Transaction) HasDescription() bool {
	if o != nil && !isNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Transaction) SetDescription(v string) {
	o.Description = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *Transaction) GetVersion() int32 {
	if o == nil || isNil(o.Version) {
		var ret int32
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetVersionOk() (*int32, bool) {
	if o == nil || isNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *Transaction) HasVersion() bool {
	if o != nil && !isNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given int32 and assigns it to the Version field.
func (o *Transaction) SetVersion(v int32) {
	o.Version = &v
}

// GetSteps returns the Steps field value if set, zero value otherwise.
func (o *Transaction) GetSteps() []string {
	if o == nil || isNil(o.Steps) {
		var ret []string
		return ret
	}
	return o.Steps
}

// GetStepsOk returns a tuple with the Steps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetStepsOk() ([]string, bool) {
	if o == nil || isNil(o.Steps) {
		return nil, false
	}
	return o.Steps, true
}

// HasSteps returns a boolean if a field has been set.
func (o *Transaction) HasSteps() bool {
	if o != nil && !isNil(o.Steps) {
		return true
	}

	return false
}

// SetSteps gets a reference to the given []string and assigns it to the Steps field.
func (o *Transaction) SetSteps(v []string) {
	o.Steps = v
}

// GetFullSteps returns the FullSteps field value if set, zero value otherwise.
func (o *Transaction) GetFullSteps() []Test {
	if o == nil || isNil(o.FullSteps) {
		var ret []Test
		return ret
	}
	return o.FullSteps
}

// GetFullStepsOk returns a tuple with the FullSteps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetFullStepsOk() ([]Test, bool) {
	if o == nil || isNil(o.FullSteps) {
		return nil, false
	}
	return o.FullSteps, true
}

// HasFullSteps returns a boolean if a field has been set.
func (o *Transaction) HasFullSteps() bool {
	if o != nil && !isNil(o.FullSteps) {
		return true
	}

	return false
}

// SetFullSteps gets a reference to the given []Test and assigns it to the FullSteps field.
func (o *Transaction) SetFullSteps(v []Test) {
	o.FullSteps = v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *Transaction) GetCreatedAt() time.Time {
	if o == nil || isNil(o.CreatedAt) {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || isNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *Transaction) HasCreatedAt() bool {
	if o != nil && !isNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *Transaction) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetSummary returns the Summary field value if set, zero value otherwise.
func (o *Transaction) GetSummary() TestSummary {
	if o == nil || isNil(o.Summary) {
		var ret TestSummary
		return ret
	}
	return *o.Summary
}

// GetSummaryOk returns a tuple with the Summary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Transaction) GetSummaryOk() (*TestSummary, bool) {
	if o == nil || isNil(o.Summary) {
		return nil, false
	}
	return o.Summary, true
}

// HasSummary returns a boolean if a field has been set.
func (o *Transaction) HasSummary() bool {
	if o != nil && !isNil(o.Summary) {
		return true
	}

	return false
}

// SetSummary gets a reference to the given TestSummary and assigns it to the Summary field.
func (o *Transaction) SetSummary(v TestSummary) {
	o.Summary = &v
}

func (o Transaction) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Transaction) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	// skip: id is readOnly
	if !isNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !isNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !isNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !isNil(o.Steps) {
		toSerialize["steps"] = o.Steps
	}
	if !isNil(o.FullSteps) {
		toSerialize["fullSteps"] = o.FullSteps
	}
	if !isNil(o.CreatedAt) {
		toSerialize["createdAt"] = o.CreatedAt
	}
	if !isNil(o.Summary) {
		toSerialize["summary"] = o.Summary
	}
	return toSerialize, nil
}

type NullableTransaction struct {
	value *Transaction
	isSet bool
}

func (v NullableTransaction) Get() *Transaction {
	return v.value
}

func (v *NullableTransaction) Set(val *Transaction) {
	v.value = val
	v.isSet = true
}

func (v NullableTransaction) IsSet() bool {
	return v.isSet
}

func (v *NullableTransaction) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTransaction(val *Transaction) *NullableTransaction {
	return &NullableTransaction{value: val, isSet: true}
}

func (v NullableTransaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTransaction) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
