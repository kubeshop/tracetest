/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ParseResponseInfo struct {
	Expression string `json:"expression,omitempty"`
}

// AssertParseResponseInfoRequired checks if the required fields are not zero-ed
func AssertParseResponseInfoRequired(obj ParseResponseInfo) error {
	return nil
}

// AssertRecurseParseResponseInfoRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ParseResponseInfo (e.g. [][]ParseResponseInfo), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseParseResponseInfoRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aParseResponseInfo, ok := obj.(ParseResponseInfo)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertParseResponseInfoRequired(aParseResponseInfo)
	})
}
