/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type TransactionVariables struct {
	Items []TestVariables
}

// AssertTransactionVariablesRequired checks if the required fields are not zero-ed
func AssertTransactionVariablesRequired(obj TransactionVariables) error {
	return nil
}

// AssertRecurseTransactionVariablesRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of TransactionVariables (e.g. [][]TransactionVariables), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTransactionVariablesRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTransactionVariables, ok := obj.(TransactionVariables)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTransactionVariablesRequired(aTransactionVariables)
	})
}
