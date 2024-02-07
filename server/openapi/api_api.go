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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ApiApiController binds http requests to an api service and writes the service results to the http response
type ApiApiController struct {
	service      ApiApiServicer
	errorHandler ErrorHandler
}

// ApiApiOption for how the controller is set up.
type ApiApiOption func(*ApiApiController)

// WithApiApiErrorHandler inject ErrorHandler into controller
func WithApiApiErrorHandler(h ErrorHandler) ApiApiOption {
	return func(c *ApiApiController) {
		c.errorHandler = h
	}
}

// NewApiApiController creates a default api controller
func NewApiApiController(s ApiApiServicer, opts ...ApiApiOption) Router {
	controller := &ApiApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ApiApiController
func (c *ApiApiController) Routes() Routes {
	return Routes{
		{
			"DeleteTestRun",
			strings.ToUpper("Delete"),
			"/api/tests/{testId}/run/{runId}",
			c.DeleteTestRun,
		},
		{
			"DeleteTestSuiteRun",
			strings.ToUpper("Delete"),
			"/api/testsuites/{testSuiteId}/run/{runId}",
			c.DeleteTestSuiteRun,
		},
		{
			"DryRunAssertion",
			strings.ToUpper("Put"),
			"/api/tests/{testId}/run/{runId}/dry-run",
			c.DryRunAssertion,
		},
		{
			"ExportTestRun",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run/{runId}/export",
			c.ExportTestRun,
		},
		{
			"ExpressionResolve",
			strings.ToUpper("Post"),
			"/api/expressions/resolve",
			c.ExpressionResolve,
		},
		{
			"GetOTLPConnectionInformation",
			strings.ToUpper("Get"),
			"/api/config/connection/otlp",
			c.GetOTLPConnectionInformation,
		},
		{
			"GetResources",
			strings.ToUpper("Get"),
			"/api/resources",
			c.GetResources,
		},
		{
			"GetRunResultJUnit",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run/{runId}/junit.xml",
			c.GetRunResultJUnit,
		},
		{
			"GetTestResultSelectedSpans",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run/{runId}/select",
			c.GetTestResultSelectedSpans,
		},
		{
			"GetTestRun",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run/{runId}",
			c.GetTestRun,
		},
		{
			"GetTestRunEvents",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run/{runId}/events",
			c.GetTestRunEvents,
		},
		{
			"GetTestRuns",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run",
			c.GetTestRuns,
		},
		{
			"GetTestSpecs",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/definition",
			c.GetTestSpecs,
		},
		{
			"GetTestSuiteRun",
			strings.ToUpper("Get"),
			"/api/testsuites/{testSuiteId}/run/{runId}",
			c.GetTestSuiteRun,
		},
		{
			"GetTestSuiteRuns",
			strings.ToUpper("Get"),
			"/api/testsuites/{testSuiteId}/run",
			c.GetTestSuiteRuns,
		},
		{
			"GetTestSuiteVersion",
			strings.ToUpper("Get"),
			"/api/testsuites/{testSuiteId}/version/{version}",
			c.GetTestSuiteVersion,
		},
		{
			"GetTestVersion",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/version/{version}",
			c.GetTestVersion,
		},
		{
			"GetVersion",
			strings.ToUpper("Get"),
			"/api/version.{fileExtension}",
			c.GetVersion,
		},
		{
			"GetWizard",
			strings.ToUpper("Get"),
			"/api/wizard",
			c.GetWizard,
		},
		{
			"ImportTestRun",
			strings.ToUpper("Post"),
			"/api/tests/import",
			c.ImportTestRun,
		},
		{
			"RerunTestRun",
			strings.ToUpper("Post"),
			"/api/tests/{testId}/run/{runId}/rerun",
			c.RerunTestRun,
		},
		{
			"ResetOTLPConnectionInformation",
			strings.ToUpper("Post"),
			"/api/config/connection/otlp/reset",
			c.ResetOTLPConnectionInformation,
		},
		{
			"RunTest",
			strings.ToUpper("Post"),
			"/api/tests/{testId}/run",
			c.RunTest,
		},
		{
			"RunTestSuite",
			strings.ToUpper("Post"),
			"/api/testsuites/{testSuiteId}/run",
			c.RunTestSuite,
		},
		{
			"SearchSpans",
			strings.ToUpper("Post"),
			"/api/tests/{testId}/run/{runId}/search-spans",
			c.SearchSpans,
		},
		{
			"SkipTraceCollection",
			strings.ToUpper("Post"),
			"/api/tests/{testId}/run/{runId}/skipPolling",
			c.SkipTraceCollection,
		},
		{
			"StopTestRun",
			strings.ToUpper("Post"),
			"/api/tests/{testId}/run/{runId}/stop",
			c.StopTestRun,
		},
		{
			"TestConnection",
			strings.ToUpper("Post"),
			"/api/config/connection",
			c.TestConnection,
		},
		{
			"UpdateTestRun",
			strings.ToUpper("Patch"),
			"/api/tests/{testId}/run/{runId}",
			c.UpdateTestRun,
		},
		{
			"UpdateWizard",
			strings.ToUpper("Put"),
			"/api/wizard",
			c.UpdateWizard,
		},
	}
}

// DeleteTestRun - delete a test run
func (c *ApiApiController) DeleteTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.DeleteTestRun(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteTestSuiteRun - Delete a specific run from a particular TestSuite
func (c *ApiApiController) DeleteTestSuiteRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.DeleteTestSuiteRun(r.Context(), testSuiteIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DryRunAssertion - run given assertions against the traces from the given run without persisting anything
func (c *ApiApiController) DryRunAssertion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	testSpecsParam := TestSpecs{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testSpecsParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestSpecsRequired(testSpecsParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.DryRunAssertion(r.Context(), testIdParam, runIdParam, testSpecsParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ExportTestRun - export test and test run information
func (c *ApiApiController) ExportTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.ExportTestRun(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ExpressionResolve - resolves an expression and returns the result string
func (c *ApiApiController) ExpressionResolve(w http.ResponseWriter, r *http.Request) {
	resolveRequestInfoParam := ResolveRequestInfo{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&resolveRequestInfoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertResolveRequestInfoRequired(resolveRequestInfoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ExpressionResolve(r.Context(), resolveRequestInfoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetOTLPConnectionInformation - get information about the OTLP connection
func (c *ApiApiController) GetOTLPConnectionInformation(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetOTLPConnectionInformation(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetResources - Get resources
func (c *ApiApiController) GetResources(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	takeParam, err := parseInt32Parameter(query.Get("take"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	skipParam, err := parseInt32Parameter(query.Get("skip"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	queryParam := query.Get("query")
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.GetResources(r.Context(), takeParam, skipParam, queryParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetRunResultJUnit - get test run results in JUnit xml format
func (c *ApiApiController) GetRunResultJUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetRunResultJUnit(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestResultSelectedSpans - retrieve spans that will be selected by selector
func (c *ApiApiController) GetTestResultSelectedSpans(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	queryParam := query.Get("query")
	result, err := c.service.GetTestResultSelectedSpans(r.Context(), testIdParam, runIdParam, queryParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestRun - get test Run
func (c *ApiApiController) GetTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestRun(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestRunEvents - get events from a test run
func (c *ApiApiController) GetTestRunEvents(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestRunEvents(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestRuns - get the runs for a test
func (c *ApiApiController) GetTestRuns(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	testIdParam := params["testId"]

	takeParam, err := parseInt32Parameter(query.Get("take"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	skipParam, err := parseInt32Parameter(query.Get("skip"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.GetTestRuns(r.Context(), testIdParam, takeParam, skipParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestSpecs - Get definition for a test
func (c *ApiApiController) GetTestSpecs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	result, err := c.service.GetTestSpecs(r.Context(), testIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestSuiteRun - Get a specific run from a particular TestSuite
func (c *ApiApiController) GetTestSuiteRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestSuiteRun(r.Context(), testSuiteIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestSuiteRuns - Get all runs from a particular TestSuite
func (c *ApiApiController) GetTestSuiteRuns(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	testSuiteIdParam := params["testSuiteId"]

	takeParam, err := parseInt32Parameter(query.Get("take"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	skipParam, err := parseInt32Parameter(query.Get("skip"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.GetTestSuiteRuns(r.Context(), testSuiteIdParam, takeParam, skipParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestSuiteVersion - get a TestSuite specific version
func (c *ApiApiController) GetTestSuiteVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	versionParam, err := parseInt32Parameter(params["version"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestSuiteVersion(r.Context(), testSuiteIdParam, versionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestVersion - get a test specific version
func (c *ApiApiController) GetTestVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	versionParam, err := parseInt32Parameter(params["version"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestVersion(r.Context(), testIdParam, versionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetVersion - Get the version of the API
func (c *ApiApiController) GetVersion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileExtensionParam := params["fileExtension"]

	result, err := c.service.GetVersion(r.Context(), fileExtensionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetWizard - Get a specific wizard
func (c *ApiApiController) GetWizard(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetWizard(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ImportTestRun - import test and test run information
func (c *ApiApiController) ImportTestRun(w http.ResponseWriter, r *http.Request) {
	exportedTestInformationParam := ExportedTestInformation{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&exportedTestInformationParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertExportedTestInformationRequired(exportedTestInformationParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.ImportTestRun(r.Context(), exportedTestInformationParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// RerunTestRun - rerun a test run
func (c *ApiApiController) RerunTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.RerunTestRun(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ResetOTLPConnectionInformation - reset the OTLP connection span count
func (c *ApiApiController) ResetOTLPConnectionInformation(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.ResetOTLPConnectionInformation(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// RunTest - run test
func (c *ApiApiController) RunTest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runInformationParam := RunInformation{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&runInformationParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRunInformationRequired(runInformationParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.RunTest(r.Context(), testIdParam, runInformationParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// RunTestSuite - run TestSuite
func (c *ApiApiController) RunTestSuite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	runInformationParam := RunInformation{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&runInformationParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertRunInformationRequired(runInformationParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.RunTestSuite(r.Context(), testSuiteIdParam, runInformationParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SearchSpans - get spans fileter by query
func (c *ApiApiController) SearchSpans(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	searchSpansRequestParam := SearchSpansRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&searchSpansRequestParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertSearchSpansRequestRequired(searchSpansRequestParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.SearchSpans(r.Context(), testIdParam, runIdParam, searchSpansRequestParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SkipTraceCollection - skips the trace collection of a test run
func (c *ApiApiController) SkipTraceCollection(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.SkipTraceCollection(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// StopTestRun - stops the execution of a test run
func (c *ApiApiController) StopTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.StopTestRun(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// TestConnection - Tests the config data store/exporter connection
func (c *ApiApiController) TestConnection(w http.ResponseWriter, r *http.Request) {
	dataStoreParam := DataStore{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&dataStoreParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertDataStoreRequired(dataStoreParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.TestConnection(r.Context(), dataStoreParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateTestRun - update a test run
func (c *ApiApiController) UpdateTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam, err := parseInt32Parameter(params["runId"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	testRunParam := TestRun{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testRunParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestRunRequired(testRunParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateTestRun(r.Context(), testIdParam, runIdParam, testRunParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateWizard - Update a Wizard
func (c *ApiApiController) UpdateWizard(w http.ResponseWriter, r *http.Request) {
	wizardParam := Wizard{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&wizardParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertWizardRequired(wizardParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateWizard(r.Context(), wizardParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
