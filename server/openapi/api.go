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
	"context"
	"net/http"
)

// ApiApiRouter defines the required methods for binding the api requests to a responses for the ApiApi
// The ApiApiRouter implementation should parse necessary information from the http request,
// pass the data to a ApiApiServicer to perform the required actions, then write the service results to the http response.
type ApiApiRouter interface {
	CreateTest(http.ResponseWriter, *http.Request)
	CreateTestFromDefinition(http.ResponseWriter, *http.Request)
	DeleteTest(http.ResponseWriter, *http.Request)
	DeleteTestRun(http.ResponseWriter, *http.Request)
	DryRunAssertion(http.ResponseWriter, *http.Request)
	ExportTestRun(http.ResponseWriter, *http.Request)
	GetRunResultJUnit(http.ResponseWriter, *http.Request)
	GetTest(http.ResponseWriter, *http.Request)
	GetTestResultSelectedSpans(http.ResponseWriter, *http.Request)
	GetTestRun(http.ResponseWriter, *http.Request)
	GetTestRuns(http.ResponseWriter, *http.Request)
	GetTestSpecs(http.ResponseWriter, *http.Request)
	GetTestVersion(http.ResponseWriter, *http.Request)
	GetTestVersionDefinitionFile(http.ResponseWriter, *http.Request)
	GetTests(http.ResponseWriter, *http.Request)
	ImportTestRun(http.ResponseWriter, *http.Request)
	RerunTestRun(http.ResponseWriter, *http.Request)
	RunTest(http.ResponseWriter, *http.Request)
	SetTestSpecs(http.ResponseWriter, *http.Request)
	UpdateTest(http.ResponseWriter, *http.Request)
	UpdateTestFromDefinition(http.ResponseWriter, *http.Request)
}

// ApiApiServicer defines the api actions for the ApiApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ApiApiServicer interface {
	CreateTest(context.Context, Test) (ImplResponse, error)
	CreateTestFromDefinition(context.Context, TextDefinition) (ImplResponse, error)
	DeleteTest(context.Context, string) (ImplResponse, error)
	DeleteTestRun(context.Context, string, string) (ImplResponse, error)
	DryRunAssertion(context.Context, string, string, TestSpecs) (ImplResponse, error)
	ExportTestRun(context.Context, string, string) (ImplResponse, error)
	GetRunResultJUnit(context.Context, string, string) (ImplResponse, error)
	GetTest(context.Context, string) (ImplResponse, error)
	GetTestResultSelectedSpans(context.Context, string, string, string) (ImplResponse, error)
	GetTestRun(context.Context, string, string) (ImplResponse, error)
	GetTestRuns(context.Context, string, int32, int32) (ImplResponse, error)
	GetTestSpecs(context.Context, string) (ImplResponse, error)
	GetTestVersion(context.Context, string, int32) (ImplResponse, error)
	GetTestVersionDefinitionFile(context.Context, string, int32) (ImplResponse, error)
	GetTests(context.Context, int32, int32, string) (ImplResponse, error)
	ImportTestRun(context.Context, ExportedTestInformation) (ImplResponse, error)
	RerunTestRun(context.Context, string, string) (ImplResponse, error)
	RunTest(context.Context, string, TestRunInformation) (ImplResponse, error)
	SetTestSpecs(context.Context, string, TestSpecs) (ImplResponse, error)
	UpdateTest(context.Context, string, Test) (ImplResponse, error)
	UpdateTestFromDefinition(context.Context, string, TextDefinition) (ImplResponse, error)
}
