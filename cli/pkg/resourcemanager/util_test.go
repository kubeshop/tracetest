package resourcemanager_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"gotest.tools/v3/assert"
)

func TestGetResourceType(t *testing.T) {
	type resource struct {
		Type string
		Spec any
	}
	testCases := []struct {
		Name          string
		Input         any
		ExpectedType  string
		ExpectedError string
	}{
		{
			Name:         "should get correct `Test` type",
			Input:        resource{Type: "Test", Spec: "anything"},
			ExpectedType: "Test",
		},
		{
			Name:         "should get correct `TestSuite` type",
			Input:        resource{Type: "TestSuite", Spec: "anything"},
			ExpectedType: "TestSuite",
		},
		{
			Name:          "should fail when input is not a struct",
			Input:         "anything",
			ExpectedError: "input must be a struct",
		},
		{
			Name:          "should fail when input is an struct without Type field",
			Input:         struct{ Value string }{Value: "anything"},
			ExpectedError: `struct has no "Type" field`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			resourceType, err := resourcemanager.GetResourceType(testCase.Input)
			if testCase.ExpectedError != "" {
				assert.Error(t, err, testCase.ExpectedError)
			}

			if testCase.ExpectedType != "" {
				assert.Equal(t, testCase.ExpectedType, resourceType)
			}
		})
	}
}
