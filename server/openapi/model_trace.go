/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Trace struct {
	TraceId string `json:"traceId,omitempty"`

	Tree Span `json:"tree,omitempty"`

	// falttened version, mapped as spanId -> span{}
	Flat map[string]Span `json:"flat,omitempty"`
}

// AssertTraceRequired checks if the required fields are not zero-ed
func AssertTraceRequired(obj Trace) error {
	if err := AssertSpanRequired(obj.Tree); err != nil {
		return err
	}
	return nil
}

// AssertRecurseTraceRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Trace (e.g. [][]Trace), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTraceRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTrace, ok := obj.(Trace)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTraceRequired(aTrace)
	})
}
