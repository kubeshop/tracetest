/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type LinterResourcePlugin struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Enabled bool `json:"enabled,omitempty"`

	Rules []LinterResourceRule `json:"rules,omitempty"`
}

// AssertLinterResourcePluginRequired checks if the required fields are not zero-ed
func AssertLinterResourcePluginRequired(obj LinterResourcePlugin) error {
	for _, el := range obj.Rules {
		if err := AssertLinterResourceRuleRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseLinterResourcePluginRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of LinterResourcePlugin (e.g. [][]LinterResourcePlugin), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseLinterResourcePluginRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aLinterResourcePlugin, ok := obj.(LinterResourcePlugin)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertLinterResourcePluginRequired(aLinterResourcePlugin)
	})
}
