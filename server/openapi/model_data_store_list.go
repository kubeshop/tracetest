/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type DataStoreList struct {
	Count int32 `json:"count,omitempty"`

	Items []DataStoreResource `json:"items,omitempty"`
}

// AssertDataStoreListRequired checks if the required fields are not zero-ed
func AssertDataStoreListRequired(obj DataStoreList) error {
	for _, el := range obj.Items {
		if err := AssertDataStoreResourceRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseDataStoreListRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of DataStoreList (e.g. [][]DataStoreList), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseDataStoreListRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aDataStoreList, ok := obj.(DataStoreList)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertDataStoreListRequired(aDataStoreList)
	})
}
