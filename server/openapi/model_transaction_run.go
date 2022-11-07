/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

type TransactionRun struct {
	Id string `json:"id,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`

	State string `json:"state,omitempty"`

	Environment Environment `json:"environment,omitempty"`
}

// AssertTransactionRunRequired checks if the required fields are not zero-ed
func AssertTransactionRunRequired(obj TransactionRun) error {
	if err := AssertEnvironmentRequired(obj.Environment); err != nil {
		return err
	}
	return nil
}

// AssertRecurseTransactionRunRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of TransactionRun (e.g. [][]TransactionRun), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTransactionRunRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTransactionRun, ok := obj.(TransactionRun)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTransactionRunRequired(aTransactionRun)
	})
}
