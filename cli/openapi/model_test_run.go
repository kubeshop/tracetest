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

// TestRun struct for TestRun
type TestRun struct {
	Id      *string `json:"id,omitempty"`
	TraceId *string `json:"traceId,omitempty"`
	SpanId  *string `json:"spanId,omitempty"`
	// Test version used when running this test run
	TestVersion *int32 `json:"testVersion,omitempty"`
	// Current execution state
	State *string `json:"state,omitempty"`
	// Details of the cause for the last `FAILED` state
	LastErrorState *string `json:"lastErrorState,omitempty"`
	// time it took for the test to complete, either success or fail. If the test is still running, it will show the time up to the time of the request
	ExecutionTime             *int32             `json:"executionTime,omitempty"`
	CreatedAt                 *time.Time         `json:"createdAt,omitempty"`
	ServiceTriggeredAt        *time.Time         `json:"serviceTriggeredAt,omitempty"`
	ServiceTriggerCompletedAt *time.Time         `json:"serviceTriggerCompletedAt,omitempty"`
	ObtainedTraceAt           *time.Time         `json:"obtainedTraceAt,omitempty"`
	CompletedAt               *time.Time         `json:"completedAt,omitempty"`
	TriggerResult             *TriggerResult     `json:"triggerResult,omitempty"`
	Trace                     *Trace             `json:"trace,omitempty"`
	Result                    *AssertionResults  `json:"result,omitempty"`
	Metadata                  *map[string]string `json:"metadata,omitempty"`
}

// NewTestRun instantiates a new TestRun object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTestRun() *TestRun {
	this := TestRun{}
	return &this
}

// NewTestRunWithDefaults instantiates a new TestRun object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTestRunWithDefaults() *TestRun {
	this := TestRun{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *TestRun) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *TestRun) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *TestRun) SetId(v string) {
	o.Id = &v
}

// GetTraceId returns the TraceId field value if set, zero value otherwise.
func (o *TestRun) GetTraceId() string {
	if o == nil || o.TraceId == nil {
		var ret string
		return ret
	}
	return *o.TraceId
}

// GetTraceIdOk returns a tuple with the TraceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetTraceIdOk() (*string, bool) {
	if o == nil || o.TraceId == nil {
		return nil, false
	}
	return o.TraceId, true
}

// HasTraceId returns a boolean if a field has been set.
func (o *TestRun) HasTraceId() bool {
	if o != nil && o.TraceId != nil {
		return true
	}

	return false
}

// SetTraceId gets a reference to the given string and assigns it to the TraceId field.
func (o *TestRun) SetTraceId(v string) {
	o.TraceId = &v
}

// GetSpanId returns the SpanId field value if set, zero value otherwise.
func (o *TestRun) GetSpanId() string {
	if o == nil || o.SpanId == nil {
		var ret string
		return ret
	}
	return *o.SpanId
}

// GetSpanIdOk returns a tuple with the SpanId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetSpanIdOk() (*string, bool) {
	if o == nil || o.SpanId == nil {
		return nil, false
	}
	return o.SpanId, true
}

// HasSpanId returns a boolean if a field has been set.
func (o *TestRun) HasSpanId() bool {
	if o != nil && o.SpanId != nil {
		return true
	}

	return false
}

// SetSpanId gets a reference to the given string and assigns it to the SpanId field.
func (o *TestRun) SetSpanId(v string) {
	o.SpanId = &v
}

// GetTestVersion returns the TestVersion field value if set, zero value otherwise.
func (o *TestRun) GetTestVersion() int32 {
	if o == nil || o.TestVersion == nil {
		var ret int32
		return ret
	}
	return *o.TestVersion
}

// GetTestVersionOk returns a tuple with the TestVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetTestVersionOk() (*int32, bool) {
	if o == nil || o.TestVersion == nil {
		return nil, false
	}
	return o.TestVersion, true
}

// HasTestVersion returns a boolean if a field has been set.
func (o *TestRun) HasTestVersion() bool {
	if o != nil && o.TestVersion != nil {
		return true
	}

	return false
}

// SetTestVersion gets a reference to the given int32 and assigns it to the TestVersion field.
func (o *TestRun) SetTestVersion(v int32) {
	o.TestVersion = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TestRun) GetState() string {
	if o == nil || o.State == nil {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetStateOk() (*string, bool) {
	if o == nil || o.State == nil {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TestRun) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *TestRun) SetState(v string) {
	o.State = &v
}

// GetLastErrorState returns the LastErrorState field value if set, zero value otherwise.
func (o *TestRun) GetLastErrorState() string {
	if o == nil || o.LastErrorState == nil {
		var ret string
		return ret
	}
	return *o.LastErrorState
}

// GetLastErrorStateOk returns a tuple with the LastErrorState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetLastErrorStateOk() (*string, bool) {
	if o == nil || o.LastErrorState == nil {
		return nil, false
	}
	return o.LastErrorState, true
}

// HasLastErrorState returns a boolean if a field has been set.
func (o *TestRun) HasLastErrorState() bool {
	if o != nil && o.LastErrorState != nil {
		return true
	}

	return false
}

// SetLastErrorState gets a reference to the given string and assigns it to the LastErrorState field.
func (o *TestRun) SetLastErrorState(v string) {
	o.LastErrorState = &v
}

// GetExecutionTime returns the ExecutionTime field value if set, zero value otherwise.
func (o *TestRun) GetExecutionTime() int32 {
	if o == nil || o.ExecutionTime == nil {
		var ret int32
		return ret
	}
	return *o.ExecutionTime
}

// GetExecutionTimeOk returns a tuple with the ExecutionTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetExecutionTimeOk() (*int32, bool) {
	if o == nil || o.ExecutionTime == nil {
		return nil, false
	}
	return o.ExecutionTime, true
}

// HasExecutionTime returns a boolean if a field has been set.
func (o *TestRun) HasExecutionTime() bool {
	if o != nil && o.ExecutionTime != nil {
		return true
	}

	return false
}

// SetExecutionTime gets a reference to the given int32 and assigns it to the ExecutionTime field.
func (o *TestRun) SetExecutionTime(v int32) {
	o.ExecutionTime = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *TestRun) GetCreatedAt() time.Time {
	if o == nil || o.CreatedAt == nil {
		var ret time.Time
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetCreatedAtOk() (*time.Time, bool) {
	if o == nil || o.CreatedAt == nil {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *TestRun) HasCreatedAt() bool {
	if o != nil && o.CreatedAt != nil {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given time.Time and assigns it to the CreatedAt field.
func (o *TestRun) SetCreatedAt(v time.Time) {
	o.CreatedAt = &v
}

// GetServiceTriggeredAt returns the ServiceTriggeredAt field value if set, zero value otherwise.
func (o *TestRun) GetServiceTriggeredAt() time.Time {
	if o == nil || o.ServiceTriggeredAt == nil {
		var ret time.Time
		return ret
	}
	return *o.ServiceTriggeredAt
}

// GetServiceTriggeredAtOk returns a tuple with the ServiceTriggeredAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetServiceTriggeredAtOk() (*time.Time, bool) {
	if o == nil || o.ServiceTriggeredAt == nil {
		return nil, false
	}
	return o.ServiceTriggeredAt, true
}

// HasServiceTriggeredAt returns a boolean if a field has been set.
func (o *TestRun) HasServiceTriggeredAt() bool {
	if o != nil && o.ServiceTriggeredAt != nil {
		return true
	}

	return false
}

// SetServiceTriggeredAt gets a reference to the given time.Time and assigns it to the ServiceTriggeredAt field.
func (o *TestRun) SetServiceTriggeredAt(v time.Time) {
	o.ServiceTriggeredAt = &v
}

// GetServiceTriggerCompletedAt returns the ServiceTriggerCompletedAt field value if set, zero value otherwise.
func (o *TestRun) GetServiceTriggerCompletedAt() time.Time {
	if o == nil || o.ServiceTriggerCompletedAt == nil {
		var ret time.Time
		return ret
	}
	return *o.ServiceTriggerCompletedAt
}

// GetServiceTriggerCompletedAtOk returns a tuple with the ServiceTriggerCompletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetServiceTriggerCompletedAtOk() (*time.Time, bool) {
	if o == nil || o.ServiceTriggerCompletedAt == nil {
		return nil, false
	}
	return o.ServiceTriggerCompletedAt, true
}

// HasServiceTriggerCompletedAt returns a boolean if a field has been set.
func (o *TestRun) HasServiceTriggerCompletedAt() bool {
	if o != nil && o.ServiceTriggerCompletedAt != nil {
		return true
	}

	return false
}

// SetServiceTriggerCompletedAt gets a reference to the given time.Time and assigns it to the ServiceTriggerCompletedAt field.
func (o *TestRun) SetServiceTriggerCompletedAt(v time.Time) {
	o.ServiceTriggerCompletedAt = &v
}

// GetObtainedTraceAt returns the ObtainedTraceAt field value if set, zero value otherwise.
func (o *TestRun) GetObtainedTraceAt() time.Time {
	if o == nil || o.ObtainedTraceAt == nil {
		var ret time.Time
		return ret
	}
	return *o.ObtainedTraceAt
}

// GetObtainedTraceAtOk returns a tuple with the ObtainedTraceAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetObtainedTraceAtOk() (*time.Time, bool) {
	if o == nil || o.ObtainedTraceAt == nil {
		return nil, false
	}
	return o.ObtainedTraceAt, true
}

// HasObtainedTraceAt returns a boolean if a field has been set.
func (o *TestRun) HasObtainedTraceAt() bool {
	if o != nil && o.ObtainedTraceAt != nil {
		return true
	}

	return false
}

// SetObtainedTraceAt gets a reference to the given time.Time and assigns it to the ObtainedTraceAt field.
func (o *TestRun) SetObtainedTraceAt(v time.Time) {
	o.ObtainedTraceAt = &v
}

// GetCompletedAt returns the CompletedAt field value if set, zero value otherwise.
func (o *TestRun) GetCompletedAt() time.Time {
	if o == nil || o.CompletedAt == nil {
		var ret time.Time
		return ret
	}
	return *o.CompletedAt
}

// GetCompletedAtOk returns a tuple with the CompletedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetCompletedAtOk() (*time.Time, bool) {
	if o == nil || o.CompletedAt == nil {
		return nil, false
	}
	return o.CompletedAt, true
}

// HasCompletedAt returns a boolean if a field has been set.
func (o *TestRun) HasCompletedAt() bool {
	if o != nil && o.CompletedAt != nil {
		return true
	}

	return false
}

// SetCompletedAt gets a reference to the given time.Time and assigns it to the CompletedAt field.
func (o *TestRun) SetCompletedAt(v time.Time) {
	o.CompletedAt = &v
}

// GetTriggerResult returns the TriggerResult field value if set, zero value otherwise.
func (o *TestRun) GetTriggerResult() TriggerResult {
	if o == nil || o.TriggerResult == nil {
		var ret TriggerResult
		return ret
	}
	return *o.TriggerResult
}

// GetTriggerResultOk returns a tuple with the TriggerResult field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetTriggerResultOk() (*TriggerResult, bool) {
	if o == nil || o.TriggerResult == nil {
		return nil, false
	}
	return o.TriggerResult, true
}

// HasTriggerResult returns a boolean if a field has been set.
func (o *TestRun) HasTriggerResult() bool {
	if o != nil && o.TriggerResult != nil {
		return true
	}

	return false
}

// SetTriggerResult gets a reference to the given TriggerResult and assigns it to the TriggerResult field.
func (o *TestRun) SetTriggerResult(v TriggerResult) {
	o.TriggerResult = &v
}

// GetTrace returns the Trace field value if set, zero value otherwise.
func (o *TestRun) GetTrace() Trace {
	if o == nil || o.Trace == nil {
		var ret Trace
		return ret
	}
	return *o.Trace
}

// GetTraceOk returns a tuple with the Trace field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetTraceOk() (*Trace, bool) {
	if o == nil || o.Trace == nil {
		return nil, false
	}
	return o.Trace, true
}

// HasTrace returns a boolean if a field has been set.
func (o *TestRun) HasTrace() bool {
	if o != nil && o.Trace != nil {
		return true
	}

	return false
}

// SetTrace gets a reference to the given Trace and assigns it to the Trace field.
func (o *TestRun) SetTrace(v Trace) {
	o.Trace = &v
}

// GetResult returns the Result field value if set, zero value otherwise.
func (o *TestRun) GetResult() AssertionResults {
	if o == nil || o.Result == nil {
		var ret AssertionResults
		return ret
	}
	return *o.Result
}

// GetResultOk returns a tuple with the Result field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetResultOk() (*AssertionResults, bool) {
	if o == nil || o.Result == nil {
		return nil, false
	}
	return o.Result, true
}

// HasResult returns a boolean if a field has been set.
func (o *TestRun) HasResult() bool {
	if o != nil && o.Result != nil {
		return true
	}

	return false
}

// SetResult gets a reference to the given AssertionResults and assigns it to the Result field.
func (o *TestRun) SetResult(v AssertionResults) {
	o.Result = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *TestRun) GetMetadata() map[string]string {
	if o == nil || o.Metadata == nil {
		var ret map[string]string
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TestRun) GetMetadataOk() (*map[string]string, bool) {
	if o == nil || o.Metadata == nil {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *TestRun) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]string and assigns it to the Metadata field.
func (o *TestRun) SetMetadata(v map[string]string) {
	o.Metadata = &v
}

func (o TestRun) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.TraceId != nil {
		toSerialize["traceId"] = o.TraceId
	}
	if o.SpanId != nil {
		toSerialize["spanId"] = o.SpanId
	}
	if o.TestVersion != nil {
		toSerialize["testVersion"] = o.TestVersion
	}
	if o.State != nil {
		toSerialize["state"] = o.State
	}
	if o.LastErrorState != nil {
		toSerialize["lastErrorState"] = o.LastErrorState
	}
	if o.ExecutionTime != nil {
		toSerialize["executionTime"] = o.ExecutionTime
	}
	if o.CreatedAt != nil {
		toSerialize["createdAt"] = o.CreatedAt
	}
	if o.ServiceTriggeredAt != nil {
		toSerialize["serviceTriggeredAt"] = o.ServiceTriggeredAt
	}
	if o.ServiceTriggerCompletedAt != nil {
		toSerialize["serviceTriggerCompletedAt"] = o.ServiceTriggerCompletedAt
	}
	if o.ObtainedTraceAt != nil {
		toSerialize["obtainedTraceAt"] = o.ObtainedTraceAt
	}
	if o.CompletedAt != nil {
		toSerialize["completedAt"] = o.CompletedAt
	}
	if o.TriggerResult != nil {
		toSerialize["triggerResult"] = o.TriggerResult
	}
	if o.Trace != nil {
		toSerialize["trace"] = o.Trace
	}
	if o.Result != nil {
		toSerialize["result"] = o.Result
	}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}
	return json.Marshal(toSerialize)
}

type NullableTestRun struct {
	value *TestRun
	isSet bool
}

func (v NullableTestRun) Get() *TestRun {
	return v.value
}

func (v *NullableTestRun) Set(val *TestRun) {
	v.value = val
	v.isSet = true
}

func (v NullableTestRun) IsSet() bool {
	return v.isSet
}

func (v *NullableTestRun) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTestRun(val *TestRun) *NullableTestRun {
	return &NullableTestRun{value: val, isSet: true}
}

func (v NullableTestRun) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTestRun) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
