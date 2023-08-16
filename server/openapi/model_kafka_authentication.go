/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type KafkaAuthentication struct {
	Type string `json:"type,omitempty"`

	Plain HttpAuthBasic `json:"plain,omitempty"`
}

// AssertKafkaAuthenticationRequired checks if the required fields are not zero-ed
func AssertKafkaAuthenticationRequired(obj KafkaAuthentication) error {
	if err := AssertHttpAuthBasicRequired(obj.Plain); err != nil {
		return err
	}
	return nil
}

// AssertRecurseKafkaAuthenticationRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of KafkaAuthentication (e.g. [][]KafkaAuthentication), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseKafkaAuthenticationRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aKafkaAuthentication, ok := obj.(KafkaAuthentication)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertKafkaAuthenticationRequired(aKafkaAuthentication)
	})
}
