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
	DeleteTransactionRun(http.ResponseWriter, *http.Request)
	DryRunAssertion(http.ResponseWriter, *http.Request)
	ExecuteDefinition(http.ResponseWriter, *http.Request)
	ExportTestRun(http.ResponseWriter, *http.Request)
	ExpressionResolve(http.ResponseWriter, *http.Request)
	GetEnvironment(http.ResponseWriter, *http.Request)
	GetEnvironmentDefinitionFile(http.ResponseWriter, *http.Request)
	GetEnvironments(http.ResponseWriter, *http.Request)
	GetResources(http.ResponseWriter, *http.Request)
	GetRunResultJUnit(http.ResponseWriter, *http.Request)
	GetTest(http.ResponseWriter, *http.Request)
	GetTestResultSelectedSpans(http.ResponseWriter, *http.Request)
	GetTestRun(http.ResponseWriter, *http.Request)
	GetTestRunEvents(http.ResponseWriter, *http.Request)
	GetTestRuns(http.ResponseWriter, *http.Request)
	GetTestSpecs(http.ResponseWriter, *http.Request)
	GetTestVersion(http.ResponseWriter, *http.Request)
	GetTestVersionDefinitionFile(http.ResponseWriter, *http.Request)
	GetTests(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	GetTransactionRun(http.ResponseWriter, *http.Request)
	GetTransactionRuns(http.ResponseWriter, *http.Request)
	GetTransactionVersion(http.ResponseWriter, *http.Request)
	GetTransactionVersionDefinitionFile(http.ResponseWriter, *http.Request)
	GetTransactions(http.ResponseWriter, *http.Request)
	ImportTestRun(http.ResponseWriter, *http.Request)
	RerunTestRun(http.ResponseWriter, *http.Request)
	RunTest(http.ResponseWriter, *http.Request)
	RunTransaction(http.ResponseWriter, *http.Request)
	StopTestRun(http.ResponseWriter, *http.Request)
	TestConnection(http.ResponseWriter, *http.Request)
	UpdateEnvironment(http.ResponseWriter, *http.Request)
	UpdateTest(http.ResponseWriter, *http.Request)
	UpdateTransaction(http.ResponseWriter, *http.Request)
	UpsertDefinition(http.ResponseWriter, *http.Request)
}

// ResourceApiApiRouter defines the required methods for binding the api requests to a responses for the ResourceApiApi
// The ResourceApiApiRouter implementation should parse necessary information from the http request,
// pass the data to a ResourceApiApiServicer to perform the required actions, then write the service results to the http response.
type ResourceApiApiRouter interface {
	CreateDemo(http.ResponseWriter, *http.Request)
	DeleteDataStore(http.ResponseWriter, *http.Request)
	DeleteDemo(http.ResponseWriter, *http.Request)
	GetConfiguration(http.ResponseWriter, *http.Request)
	GetDataStore(http.ResponseWriter, *http.Request)
	GetDemo(http.ResponseWriter, *http.Request)
	GetPollingProfile(http.ResponseWriter, *http.Request)
	ListDemos(http.ResponseWriter, *http.Request)
	UpdateConfiguration(http.ResponseWriter, *http.Request)
	UpdateDataStore(http.ResponseWriter, *http.Request)
	UpdateDemo(http.ResponseWriter, *http.Request)
	UpdatePollingProfile(http.ResponseWriter, *http.Request)
}

// ApiApiServicer defines the api actions for the ApiApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ApiApiServicer interface {
	CreateEnvironment(context.Context, Environment) (ImplResponse, error)
	CreateTest(context.Context, Test) (ImplResponse, error)
	CreateTransaction(context.Context, Transaction) (ImplResponse, error)
	DeleteEnvironment(context.Context, string) (ImplResponse, error)
	DeleteTest(context.Context, string) (ImplResponse, error)
	DeleteTestRun(context.Context, string, int32) (ImplResponse, error)
	DeleteTransaction(context.Context, string) (ImplResponse, error)
	DeleteTransactionRun(context.Context, string, int32) (ImplResponse, error)
	DryRunAssertion(context.Context, string, int32, TestSpecs) (ImplResponse, error)
	ExecuteDefinition(context.Context, TextDefinition) (ImplResponse, error)
	ExportTestRun(context.Context, string, int32) (ImplResponse, error)
	ExpressionResolve(context.Context, ResolveRequestInfo) (ImplResponse, error)
	GetEnvironment(context.Context, string) (ImplResponse, error)
	GetEnvironmentDefinitionFile(context.Context, string) (ImplResponse, error)
	GetEnvironments(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	GetResources(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	GetRunResultJUnit(context.Context, string, int32) (ImplResponse, error)
	GetTest(context.Context, string) (ImplResponse, error)
	GetTestResultSelectedSpans(context.Context, string, int32, string) (ImplResponse, error)
	GetTestRun(context.Context, string, int32) (ImplResponse, error)
	GetTestRunEvents(context.Context, string, int32) (ImplResponse, error)
	GetTestRuns(context.Context, string, int32, int32) (ImplResponse, error)
	GetTestSpecs(context.Context, string) (ImplResponse, error)
	GetTestVersion(context.Context, string, int32) (ImplResponse, error)
	GetTestVersionDefinitionFile(context.Context, string, int32) (ImplResponse, error)
	GetTests(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	GetTransaction(context.Context, string) (ImplResponse, error)
	GetTransactionRun(context.Context, string, int32) (ImplResponse, error)
	GetTransactionRuns(context.Context, string, int32, int32) (ImplResponse, error)
	GetTransactionVersion(context.Context, string, int32) (ImplResponse, error)
	GetTransactionVersionDefinitionFile(context.Context, string, int32) (ImplResponse, error)
	GetTransactions(context.Context, int32, int32, string, string, string) (ImplResponse, error)
	ImportTestRun(context.Context, ExportedTestInformation) (ImplResponse, error)
	RerunTestRun(context.Context, string, int32) (ImplResponse, error)
	RunTest(context.Context, string, RunInformation) (ImplResponse, error)
	RunTransaction(context.Context, string, RunInformation) (ImplResponse, error)
	StopTestRun(context.Context, string, int32) (ImplResponse, error)
	TestConnection(context.Context, DataStore) (ImplResponse, error)
	UpdateEnvironment(context.Context, string, Environment) (ImplResponse, error)
	UpdateTest(context.Context, string, Test) (ImplResponse, error)
	UpdateTransaction(context.Context, string, Transaction) (ImplResponse, error)
	UpsertDefinition(context.Context, TextDefinition) (ImplResponse, error)
}

// ResourceApiApiServicer defines the api actions for the ResourceApiApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type ResourceApiApiServicer interface {
	CreateDemo(context.Context, Demo) (ImplResponse, error)
	DeleteDataStore(context.Context, string) (ImplResponse, error)
	DeleteDemo(context.Context, string) (ImplResponse, error)
	GetConfiguration(context.Context, string) (ImplResponse, error)
	GetDataStore(context.Context, string) (ImplResponse, error)
	GetDemo(context.Context, string) (ImplResponse, error)
	GetPollingProfile(context.Context, string) (ImplResponse, error)
	ListDemos(context.Context, int32, int32, string, string) (ImplResponse, error)
	UpdateConfiguration(context.Context, string, ConfigurationResource) (ImplResponse, error)
	UpdateDataStore(context.Context, string, DataStore) (ImplResponse, error)
	UpdateDemo(context.Context, string, Demo) (ImplResponse, error)
	UpdatePollingProfile(context.Context, string, PollingProfile) (ImplResponse, error)
}
