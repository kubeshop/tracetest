/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type AwsXRay struct {
	Region string `json:"region,omitempty"`

	AccessKeyId string `json:"accessKeyId,omitempty"`

	SecretAccessKey string `json:"secretAccessKey,omitempty"`
}

// AssertAwsXRayRequired checks if the required fields are not zero-ed
func AssertAwsXRayRequired(obj AwsXRay) error {
	return nil
}

// AssertRecurseAwsXRayRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of AwsXRay (e.g. [][]AwsXRay), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAwsXRayRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAwsXRay, ok := obj.(AwsXRay)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAwsXRayRequired(aAwsXRay)
	})
}
