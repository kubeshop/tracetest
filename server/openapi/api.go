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
	CreateEnvironment(http.ResponseWriter, *http.Request)
	CreateTest(http.ResponseWriter, *http.Request)
	CreateTransaction(http.ResponseWriter, *http.Request)
	DeleteEnvironment(http.ResponseWriter, *http.Request)
	DeleteTest(http.ResponseWriter, *http.Request)
	DeleteTestRun(http.ResponseWriter, *http.Request)
	DeleteTransaction(http.ResponseWriter, *http.Request)
	DryRunAssertion(http.ResponseWriter, *http.Request)
	ExecuteDefinition(http.ResponseWriter, *http.Request)
	ExportTestRun(http.ResponseWriter, *http.Request)
	ExpressionResolve(http.ResponseWriter, *http.Request)
	GetEnvironment(http.ResponseWriter, *http.Request)
	GetEnvironments(http.ResponseWriter, *http.Request)
	GetResources(http.ResponseWriter, *http.Request)
	GetRunResultJUnit(http.ResponseWriter, *http.Request)
	GetTest(http.ResponseWriter, *http.Request)
	GetTestResultSelectedSpans(http.ResponseWriter, *http.Request)
	GetTestRun(http.ResponseWriter, *http.Request)
	GetTestRuns(http.ResponseWriter, *http.Request)
	GetTestSpecs(http.ResponseWriter, *http.Request)
	GetTestVersion(http.ResponseWriter, *http.Request)
	GetTestVersionDefinitionFile(http.ResponseWriter, *http.Request)
	GetTests(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	GetTransactions(http.ResponseWriter, *http.Request)
	ImportTestRun(http.ResponseWriter, *http.Request)
	RerunTestRun(http.ResponseWriter, *http.Request)
	RunTest(http.ResponseWriter, *http.Request)
	SetTestOutputs(http.ResponseWriter, *http.Request)
	SetTestSpecs(http.ResponseWriter, *http.Request)
	UpdateEnvironment(http.ResponseWriter, *http.Request)
	UpdateTest(http.ResponseWriter, *http.Request)
	UpdateTransaction(http.ResponseWriter, *http.Request)
}

// ApiApiServicer defines the api actions for the ApiApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ApiApiServicer interface {
	CreateEnvironment(context.Context, Environment) (ImplResponse, error)
	CreateTest(context.Context, Test) (ImplResponse, error)
	CreateTransaction(context.Context, Transaction) (ImplResponse, error)
	DeleteEnvironment(context.Context, string) (ImplResponse, error)
	DeleteTest(context.Context, string) (ImplResponse, error)
	DeleteTestRun(context.Context, string, string) (ImplResponse, error)
	DeleteTransaction(context.Context, string) (ImplResponse, error)
	DryRunAssertion(context.Context, string, string, TestSpecs) (ImplResponse, error)
	ExecuteDefinition(context.Context, TextDefinition) (ImplResponse, error)
	ExportTestRun(context.Context, string, string) (ImplResponse, error)
	ExpressionResolve(context.Context, ResolveRequestInfo) (ImplResponse, error)
	GetEnvironment(context.Context, string) (ImplResponse, error)
	GetEnvironments(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	GetResources(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	GetRunResultJUnit(context.Context, string, string) (ImplResponse, error)
	GetTest(context.Context, string) (ImplResponse, error)
	GetTestResultSelectedSpans(context.Context, string, string, string) (ImplResponse, error)
	GetTestRun(context.Context, string, string) (ImplResponse, error)
	GetTestRuns(context.Context, string, int32, int32) (ImplResponse, error)
	GetTestSpecs(context.Context, string) (ImplResponse, error)
	GetTestVersion(context.Context, string, int32) (ImplResponse, error)
	GetTestVersionDefinitionFile(context.Context, string, int32) (ImplResponse, error)
	GetTests(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	GetTransaction(context.Context, string) (ImplResponse, error)
	GetTransactions(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	ImportTestRun(context.Context, ExportedTestInformation) (ImplResponse, error)
	RerunTestRun(context.Context, string, string) (ImplResponse, error)
	RunTest(context.Context, string, TestRunInformation) (ImplResponse, error)
	SetTestOutputs(context.Context, string, []TestOutput) (ImplResponse, error)
	SetTestSpecs(context.Context, string, TestSpecs) (ImplResponse, error)
	UpdateEnvironment(context.Context, string, Environment) (ImplResponse, error)
	UpdateTest(context.Context, string, Test) (ImplResponse, error)
	UpdateTransaction(context.Context, string, Transaction) (ImplResponse, error)
}
