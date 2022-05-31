package conversion_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/conversion"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
)

func TestDefinitionToOpenAPIConversion(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          definition.Test
		ExpectedOutput openapi.Test
	}{
		{
			Name: "Should be able to convert request with no authentication and no body",
			Input: definition.Test{
				Name:        "A test",
				Description: "A test description",
				Trigger: definition.TestTrigger{
					Type: "http",
					HTTPRequest: definition.HttpRequest{
						URL:    "http://localhost:1234",
						Method: "POST",
						Headers: []definition.HTTPHeader{
							{Key: "Content-Type", Value: "application/json"},
						},
						Authentication: definition.HTTPAuthentication{},
						Body:           definition.HTTPBody{},
					},
				},
			},
			ExpectedOutput: openapi.Test{
				Name:        conversion.ConvertStringIntoOpenAPIString("A test"),
				Description: conversion.ConvertStringIntoOpenAPIString("A test description"),
				ServiceUnderTest: &openapi.TestServiceUnderTest{
					Request: &openapi.HTTPRequest{
						Url:    conversion.ConvertStringIntoOpenAPIString("http://localhost:1234"),
						Method: conversion.ConvertStringIntoOpenAPIString("POST"),
						Headers: []openapi.HTTPHeader{
							{Key: conversion.ConvertStringIntoOpenAPIString("Content-Type"), Value: conversion.ConvertStringIntoOpenAPIString("application/json")},
						},
						Body: conversion.ConvertStringIntoOpenAPIString(""),
						Auth: &openapi.HTTPAuth{},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			output := conversion.ConvertTestDefinitionIntoOpenAPIObject(testCase.Input)
			assert.Equal(t, testCase.ExpectedOutput, output)
		})
	}
}
