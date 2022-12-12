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

// DataStore struct for DataStore
type DataStore struct {
	Type       *SupportedDataStores `json:"type,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Jaeger     *GRPCClientSettings  `json:"jaeger,omitempty"`
	Tempo      *GRPCClientSettings  `json:"tempo,omitempty"`
	OpenSearch *OpenSearch          `json:"openSearch,omitempty"`
	SignalFx   *SignalFX            `json:"signalFx,omitempty"`
}

// NewDataStore instantiates a new DataStore object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDataStore() *DataStore {
	this := DataStore{}
	return &this
}

// NewDataStoreWithDefaults instantiates a new DataStore object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDataStoreWithDefaults() *DataStore {
	this := DataStore{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *DataStore) GetType() SupportedDataStores {
	if o == nil || o.Type == nil {
		var ret SupportedDataStores
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataStore) GetTypeOk() (*SupportedDataStores, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *DataStore) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given SupportedDataStores and assigns it to the Type field.
func (o *DataStore) SetType(v SupportedDataStores) {
	o.Type = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *DataStore) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataStore) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *DataStore) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *DataStore) SetName(v string) {
	o.Name = &v
}

// GetJaeger returns the Jaeger field value if set, zero value otherwise.
func (o *DataStore) GetJaeger() GRPCClientSettings {
	if o == nil || o.Jaeger == nil {
		var ret GRPCClientSettings
		return ret
	}
	return *o.Jaeger
}

// GetJaegerOk returns a tuple with the Jaeger field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataStore) GetJaegerOk() (*GRPCClientSettings, bool) {
	if o == nil || o.Jaeger == nil {
		return nil, false
	}
	return o.Jaeger, true
}

// HasJaeger returns a boolean if a field has been set.
func (o *DataStore) HasJaeger() bool {
	if o != nil && o.Jaeger != nil {
		return true
	}

	return false
}

// SetJaeger gets a reference to the given GRPCClientSettings and assigns it to the Jaeger field.
func (o *DataStore) SetJaeger(v GRPCClientSettings) {
	o.Jaeger = &v
}

// GetTempo returns the Tempo field value if set, zero value otherwise.
func (o *DataStore) GetTempo() GRPCClientSettings {
	if o == nil || o.Tempo == nil {
		var ret GRPCClientSettings
		return ret
	}
	return *o.Tempo
}

// GetTempoOk returns a tuple with the Tempo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataStore) GetTempoOk() (*GRPCClientSettings, bool) {
	if o == nil || o.Tempo == nil {
		return nil, false
	}
	return o.Tempo, true
}

// HasTempo returns a boolean if a field has been set.
func (o *DataStore) HasTempo() bool {
	if o != nil && o.Tempo != nil {
		return true
	}

	return false
}

// SetTempo gets a reference to the given GRPCClientSettings and assigns it to the Tempo field.
func (o *DataStore) SetTempo(v GRPCClientSettings) {
	o.Tempo = &v
}

// GetOpenSearch returns the OpenSearch field value if set, zero value otherwise.
func (o *DataStore) GetOpenSearch() OpenSearch {
	if o == nil || o.OpenSearch == nil {
		var ret OpenSearch
		return ret
	}
	return *o.OpenSearch
}

// GetOpenSearchOk returns a tuple with the OpenSearch field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataStore) GetOpenSearchOk() (*OpenSearch, bool) {
	if o == nil || o.OpenSearch == nil {
		return nil, false
	}
	return o.OpenSearch, true
}

// HasOpenSearch returns a boolean if a field has been set.
func (o *DataStore) HasOpenSearch() bool {
	if o != nil && o.OpenSearch != nil {
		return true
	}

	return false
}

// SetOpenSearch gets a reference to the given OpenSearch and assigns it to the OpenSearch field.
func (o *DataStore) SetOpenSearch(v OpenSearch) {
	o.OpenSearch = &v
}

// GetSignalFx returns the SignalFx field value if set, zero value otherwise.
func (o *DataStore) GetSignalFx() SignalFX {
	if o == nil || o.SignalFx == nil {
		var ret SignalFX
		return ret
	}
	return *o.SignalFx
}

// GetSignalFxOk returns a tuple with the SignalFx field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DataStore) GetSignalFxOk() (*SignalFX, bool) {
	if o == nil || o.SignalFx == nil {
		return nil, false
	}
	return o.SignalFx, true
}

// HasSignalFx returns a boolean if a field has been set.
func (o *DataStore) HasSignalFx() bool {
	if o != nil && o.SignalFx != nil {
		return true
	}

	return false
}

// SetSignalFx gets a reference to the given SignalFX and assigns it to the SignalFx field.
func (o *DataStore) SetSignalFx(v SignalFX) {
	o.SignalFx = &v
}

func (o DataStore) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Jaeger != nil {
		toSerialize["jaeger"] = o.Jaeger
	}
	if o.Tempo != nil {
		toSerialize["tempo"] = o.Tempo
	}
	if o.OpenSearch != nil {
		toSerialize["openSearch"] = o.OpenSearch
	}
	if o.SignalFx != nil {
		toSerialize["signalFx"] = o.SignalFx
	}
	return json.Marshal(toSerialize)
}

type NullableDataStore struct {
	value *DataStore
	isSet bool
}

func (v NullableDataStore) Get() *DataStore {
	return v.value
}

func (v *NullableDataStore) Set(val *DataStore) {
	v.value = val
	v.isSet = true
}

func (v NullableDataStore) IsSet() bool {
	return v.isSet
}

func (v *NullableDataStore) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDataStore(val *DataStore) *NullableDataStore {
	return &NullableDataStore{value: val, isSet: true}
}

func (v NullableDataStore) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDataStore) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
