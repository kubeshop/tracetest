/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type LinterResultPluginRuleResult struct {
	SpanId string `json:"spanId,omitempty"`

	Errors []LinterResultPluginRuleResultError `json:"errors,omitempty"`

	Passed bool `json:"passed,omitempty"`

	Severity string `json:"severity,omitempty"`
}

// AssertLinterResultPluginRuleResultRequired checks if the required fields are not zero-ed
func AssertLinterResultPluginRuleResultRequired(obj LinterResultPluginRuleResult) error {
	for _, el := range obj.Errors {
		if err := AssertLinterResultPluginRuleResultErrorRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseLinterResultPluginRuleResultRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of LinterResultPluginRuleResult (e.g. [][]LinterResultPluginRuleResult), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseLinterResultPluginRuleResultRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aLinterResultPluginRuleResult, ok := obj.(LinterResultPluginRuleResult)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertLinterResultPluginRuleResultRequired(aLinterResultPluginRuleResult)
	})
}
