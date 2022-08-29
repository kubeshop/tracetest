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

// Routes returns all of the api route for the ApiApiController
func (c *ApiApiController) Routes() Routes {
	return Routes{
		{
			"CreateTest",
			strings.ToUpper("Post"),
			"/api/tests",
			c.CreateTest,
		},
		{
			"CreateTestFromDefinition",
			strings.ToUpper("Post"),
			"/api/tests/definition.yaml",
			c.CreateTestFromDefinition,
		},
		{
			"DeleteTest",
			strings.ToUpper("Delete"),
			"/api/tests/{testId}",
			c.DeleteTest,
		},
		{
			"DeleteTestRun",
			strings.ToUpper("Delete"),
			"/api/tests/{testId}/run/{runId}",
			c.DeleteTestRun,
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
			"GetRunResultJUnit",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/run/{runId}/junit.xml",
			c.GetRunResultJUnit,
		},
		{
			"GetTest",
			strings.ToUpper("Get"),
			"/api/tests/{testId}",
			c.GetTest,
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
			"GetTestVersionDefinitionFile",
			strings.ToUpper("Get"),
			"/api/tests/{testId}/version/{version}/definition.yaml",
			c.GetTestVersionDefinitionFile,
		},
		{
			"GetTests",
			strings.ToUpper("Get"),
			"/api/tests",
			c.GetTests,
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
			"RunTest",
			strings.ToUpper("Post"),
			"/api/tests/{testId}/run",
			c.RunTest,
		},
		{
			"SetTestSpecs",
			strings.ToUpper("Put"),
			"/api/tests/{testId}/definition",
			c.SetTestSpecs,
		},
		{
			"UpdateTest",
			strings.ToUpper("Put"),
			"/api/tests/{testId}",
			c.UpdateTest,
		},
		{
			"UpdateTestFromDefinition",
			strings.ToUpper("Put"),
			"/api/tests/{testId}/definition.yaml",
			c.UpdateTestFromDefinition,
		},
	}
}

// CreateTest - Create new test
func (c *ApiApiController) CreateTest(w http.ResponseWriter, r *http.Request) {
	testParam := Test{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestRequired(testParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateTest(r.Context(), testParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CreateTestFromDefinition - Create new test using the yaml definition
func (c *ApiApiController) CreateTestFromDefinition(w http.ResponseWriter, r *http.Request) {
	textDefinitionParam := TextDefinition{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&textDefinitionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTextDefinitionRequired(textDefinitionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateTestFromDefinition(r.Context(), textDefinitionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteTest - delete a test
func (c *ApiApiController) DeleteTest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	result, err := c.service.DeleteTest(r.Context(), testIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteTestRun - delete a test run
func (c *ApiApiController) DeleteTestRun(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam := params["runId"]

	result, err := c.service.DeleteTestRun(r.Context(), testIdParam, runIdParam)
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

	runIdParam := params["runId"]

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

	runIdParam := params["runId"]

	result, err := c.service.ExportTestRun(r.Context(), testIdParam, runIdParam)
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

	runIdParam := params["runId"]

	result, err := c.service.GetRunResultJUnit(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTest - get test
func (c *ApiApiController) GetTest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	result, err := c.service.GetTest(r.Context(), testIdParam)
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

	runIdParam := params["runId"]

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

	runIdParam := params["runId"]

	result, err := c.service.GetTestRun(r.Context(), testIdParam, runIdParam)
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

// GetTestVersionDefinitionFile - Get the test definition as an YAML file
func (c *ApiApiController) GetTestVersionDefinitionFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	versionParam, err := parseInt32Parameter(params["version"], true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestVersionDefinitionFile(r.Context(), testIdParam, versionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTests - Get tests
func (c *ApiApiController) GetTests(w http.ResponseWriter, r *http.Request) {
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
	result, err := c.service.GetTests(r.Context(), takeParam, skipParam, queryParam)
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

	runIdParam := params["runId"]

	result, err := c.service.RerunTestRun(r.Context(), testIdParam, runIdParam)
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

	result, err := c.service.RunTest(r.Context(), testIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// SetTestSpecs - Set spec for a test
func (c *ApiApiController) SetTestSpecs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

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
	result, err := c.service.SetTestSpecs(r.Context(), testIdParam, testSpecsParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateTest - update test
func (c *ApiApiController) UpdateTest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	testParam := Test{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestRequired(testParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateTest(r.Context(), testIdParam, testParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateTestFromDefinition - update test from definition file
func (c *ApiApiController) UpdateTestFromDefinition(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	textDefinitionParam := TextDefinition{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&textDefinitionParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTextDefinitionRequired(textDefinitionParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateTestFromDefinition(r.Context(), testIdParam, textDefinitionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
