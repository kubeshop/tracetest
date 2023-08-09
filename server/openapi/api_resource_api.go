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

// ResourceApiApiController binds http requests to an api service and writes the service results to the http response
type ResourceApiApiController struct {
	service      ResourceApiApiServicer
	errorHandler ErrorHandler
}

// ResourceApiApiOption for how the controller is set up.
type ResourceApiApiOption func(*ResourceApiApiController)

// WithResourceApiApiErrorHandler inject ErrorHandler into controller
func WithResourceApiApiErrorHandler(h ErrorHandler) ResourceApiApiOption {
	return func(c *ResourceApiApiController) {
		c.errorHandler = h
	}
}

// NewResourceApiApiController creates a default api controller
func NewResourceApiApiController(s ResourceApiApiServicer, opts ...ResourceApiApiOption) Router {
	controller := &ResourceApiApiController{
		service:      s,
		errorHandler: DefaultErrorHandler,
	}

	for _, opt := range opts {
		opt(controller)
	}

	return controller
}

// Routes returns all the api routes for the ResourceApiApiController
func (c *ResourceApiApiController) Routes() Routes {
	return Routes{
		{
			"CreateDemo",
			strings.ToUpper("Post"),
			"/api/demos",
			c.CreateDemo,
		},
		{
			"CreateLinter",
			strings.ToUpper("Post"),
			"/api/linters",
			c.CreateLinter,
		},
		{
			"CreateTest",
			strings.ToUpper("Post"),
			"/api/tests",
			c.CreateTest,
		},
		{
			"CreateTestSuite",
			strings.ToUpper("Post"),
			"/api/testsuites",
			c.CreateTestSuite,
		},
		{
			"CreateVariableSet",
			strings.ToUpper("Post"),
			"/api/variableSets",
			c.CreateVariableSet,
		},
		{
			"DeleteDataStore",
			strings.ToUpper("Delete"),
			"/api/datastores/{dataStoreId}",
			c.DeleteDataStore,
		},
		{
			"DeleteDemo",
			strings.ToUpper("Delete"),
			"/api/demos/{demoId}",
			c.DeleteDemo,
		},
		{
			"DeleteLinter",
			strings.ToUpper("Delete"),
			"/api/linters/{LinterId}",
			c.DeleteLinter,
		},
		{
			"DeleteTest",
			strings.ToUpper("Delete"),
			"/api/tests/{testId}",
			c.DeleteTest,
		},
		{
			"DeleteTestSuite",
			strings.ToUpper("Delete"),
			"/api/testsuites/{testSuiteId}",
			c.DeleteTestSuite,
		},
		{
			"DeleteVariableSet",
			strings.ToUpper("Delete"),
			"/api/variableSets/{variableSetId}",
			c.DeleteVariableSet,
		},
		{
			"GetConfiguration",
			strings.ToUpper("Get"),
			"/api/configs/{configId}",
			c.GetConfiguration,
		},
		{
			"GetDataStore",
			strings.ToUpper("Get"),
			"/api/datastores/{dataStoreId}",
			c.GetDataStore,
		},
		{
			"GetDemo",
			strings.ToUpper("Get"),
			"/api/demos/{demoId}",
			c.GetDemo,
		},
		{
			"GetLinter",
			strings.ToUpper("Get"),
			"/api/linters/{LinterId}",
			c.GetLinter,
		},
		{
			"GetPollingProfile",
			strings.ToUpper("Get"),
			"/api/pollingprofiles/{pollingProfileId}",
			c.GetPollingProfile,
		},
		{
			"GetTestSuite",
			strings.ToUpper("Get"),
			"/api/testsuites/{testSuiteId}",
			c.GetTestSuite,
		},
		{
			"GetTestSuites",
			strings.ToUpper("Get"),
			"/api/testsuites",
			c.GetTestSuites,
		},
		{
			"GetTests",
			strings.ToUpper("Get"),
			"/api/tests",
			c.GetTests,
		},
		{
			"GetVariableSet",
			strings.ToUpper("Get"),
			"/api/variableSets/{variableSetId}",
			c.GetVariableSet,
		},
		{
			"ListConfiguration",
			strings.ToUpper("Get"),
			"/api/configs",
			c.ListConfiguration,
		},
		{
			"ListDataStore",
			strings.ToUpper("Get"),
			"/api/datastores",
			c.ListDataStore,
		},
		{
			"ListDemos",
			strings.ToUpper("Get"),
			"/api/demos",
			c.ListDemos,
		},
		{
			"ListLinters",
			strings.ToUpper("Get"),
			"/api/linters",
			c.ListLinters,
		},
		{
			"ListPollingProfile",
			strings.ToUpper("Get"),
			"/api/pollingprofiles",
			c.ListPollingProfile,
		},
		{
			"ListVariableSets",
			strings.ToUpper("Get"),
			"/api/variableSets",
			c.ListVariableSets,
		},
		{
			"TestsTestIdGet",
			strings.ToUpper("Get"),
			"/api/tests/{testId}",
			c.TestsTestIdGet,
		},
		{
			"UpdateConfiguration",
			strings.ToUpper("Put"),
			"/api/configs/{configId}",
			c.UpdateConfiguration,
		},
		{
			"UpdateDataStore",
			strings.ToUpper("Put"),
			"/api/datastores/{dataStoreId}",
			c.UpdateDataStore,
		},
		{
			"UpdateDemo",
			strings.ToUpper("Put"),
			"/api/demos/{demoId}",
			c.UpdateDemo,
		},
		{
			"UpdateLinter",
			strings.ToUpper("Put"),
			"/api/linters/{LinterId}",
			c.UpdateLinter,
		},
		{
			"UpdatePollingProfile",
			strings.ToUpper("Put"),
			"/api/pollingprofiles/{pollingProfileId}",
			c.UpdatePollingProfile,
		},
		{
			"UpdateTest",
			strings.ToUpper("Put"),
			"/api/tests/{testId}",
			c.UpdateTest,
		},
		{
			"UpdateTestSuite",
			strings.ToUpper("Put"),
			"/api/testsuites/{testSuiteId}",
			c.UpdateTestSuite,
		},
		{
			"UpdateVariableSet",
			strings.ToUpper("Put"),
			"/api/variableSets/{variableSetId}",
			c.UpdateVariableSet,
		},
	}
}

// CreateDemo - Create a Demonstration setting
func (c *ResourceApiApiController) CreateDemo(w http.ResponseWriter, r *http.Request) {
	demoParam := Demo{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&demoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertDemoRequired(demoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateDemo(r.Context(), demoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CreateLinter - Create an Linter
func (c *ResourceApiApiController) CreateLinter(w http.ResponseWriter, r *http.Request) {
	linterResourceParam := LinterResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&linterResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertLinterResourceRequired(linterResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateLinter(r.Context(), linterResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CreateTest - Create new test
func (c *ResourceApiApiController) CreateTest(w http.ResponseWriter, r *http.Request) {
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

// CreateTestSuite - Create new TestSuite
func (c *ResourceApiApiController) CreateTestSuite(w http.ResponseWriter, r *http.Request) {
	testSuiteResourceParam := TestSuiteResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testSuiteResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestSuiteResourceRequired(testSuiteResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateTestSuite(r.Context(), testSuiteResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CreateVariableSet - Create a VariableSet
func (c *ResourceApiApiController) CreateVariableSet(w http.ResponseWriter, r *http.Request) {
	variableSetResourceParam := VariableSetResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&variableSetResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertVariableSetResourceRequired(variableSetResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateVariableSet(r.Context(), variableSetResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteDataStore - Delete a Data Store
func (c *ResourceApiApiController) DeleteDataStore(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dataStoreIdParam := params["dataStoreId"]

	result, err := c.service.DeleteDataStore(r.Context(), dataStoreIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteDemo - Delete a Demonstration setting
func (c *ResourceApiApiController) DeleteDemo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	demoIdParam := params["demoId"]

	result, err := c.service.DeleteDemo(r.Context(), demoIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteLinter - Delete an Linter
func (c *ResourceApiApiController) DeleteLinter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	linterIdParam := params["LinterId"]

	result, err := c.service.DeleteLinter(r.Context(), linterIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteTest - delete a test
func (c *ResourceApiApiController) DeleteTest(w http.ResponseWriter, r *http.Request) {
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

// DeleteTestSuite - delete a TestSuite
func (c *ResourceApiApiController) DeleteTestSuite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	result, err := c.service.DeleteTestSuite(r.Context(), testSuiteIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// DeleteVariableSet - Delete an variable set
func (c *ResourceApiApiController) DeleteVariableSet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	variableSetIdParam := params["variableSetId"]

	result, err := c.service.DeleteVariableSet(r.Context(), variableSetIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetConfiguration - Get Tracetest configuration
func (c *ResourceApiApiController) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	configIdParam := params["configId"]

	result, err := c.service.GetConfiguration(r.Context(), configIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetDataStore - Get a Data Store
func (c *ResourceApiApiController) GetDataStore(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dataStoreIdParam := params["dataStoreId"]

	result, err := c.service.GetDataStore(r.Context(), dataStoreIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetDemo - Get Demonstration setting
func (c *ResourceApiApiController) GetDemo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	demoIdParam := params["demoId"]

	result, err := c.service.GetDemo(r.Context(), demoIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetLinter - Get a specific Linter
func (c *ResourceApiApiController) GetLinter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	linterIdParam := params["LinterId"]

	result, err := c.service.GetLinter(r.Context(), linterIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetPollingProfile - Get Polling Profile
func (c *ResourceApiApiController) GetPollingProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pollingProfileIdParam := params["pollingProfileId"]

	result, err := c.service.GetPollingProfile(r.Context(), pollingProfileIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestSuite - get TestSuite
func (c *ResourceApiApiController) GetTestSuite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	result, err := c.service.GetTestSuite(r.Context(), testSuiteIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTestSuites - Get testsuites
func (c *ResourceApiApiController) GetTestSuites(w http.ResponseWriter, r *http.Request) {
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
	result, err := c.service.GetTestSuites(r.Context(), takeParam, skipParam, queryParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTests - Get tests
func (c *ResourceApiApiController) GetTests(w http.ResponseWriter, r *http.Request) {
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
	result, err := c.service.GetTests(r.Context(), takeParam, skipParam, queryParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetVariableSet - Get a specific VariableSet
func (c *ResourceApiApiController) GetVariableSet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	variableSetIdParam := params["variableSetId"]

	result, err := c.service.GetVariableSet(r.Context(), variableSetIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListConfiguration - List Tracetest configuration
func (c *ResourceApiApiController) ListConfiguration(w http.ResponseWriter, r *http.Request) {
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
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.ListConfiguration(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListDataStore - List Data Store
func (c *ResourceApiApiController) ListDataStore(w http.ResponseWriter, r *http.Request) {
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
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.ListDataStore(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListDemos - List Demonstrations
func (c *ResourceApiApiController) ListDemos(w http.ResponseWriter, r *http.Request) {
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
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.ListDemos(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListLinters - List Linters
func (c *ResourceApiApiController) ListLinters(w http.ResponseWriter, r *http.Request) {
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
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.ListLinters(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListPollingProfile - List Polling Profile Configuration
func (c *ResourceApiApiController) ListPollingProfile(w http.ResponseWriter, r *http.Request) {
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
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.ListPollingProfile(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// ListVariableSets - List VariableSets
func (c *ResourceApiApiController) ListVariableSets(w http.ResponseWriter, r *http.Request) {
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
	sortByParam := query.Get("sortBy")
	sortDirectionParam := query.Get("sortDirection")
	result, err := c.service.ListVariableSets(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// TestsTestIdGet - get test
func (c *ResourceApiApiController) TestsTestIdGet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	result, err := c.service.TestsTestIdGet(r.Context(), testIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateConfiguration - Update Tracetest configuration
func (c *ResourceApiApiController) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	configIdParam := params["configId"]

	configurationResourceParam := ConfigurationResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&configurationResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertConfigurationResourceRequired(configurationResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateConfiguration(r.Context(), configIdParam, configurationResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateDataStore - Update a Data Store
func (c *ResourceApiApiController) UpdateDataStore(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	dataStoreIdParam := params["dataStoreId"]

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
	result, err := c.service.UpdateDataStore(r.Context(), dataStoreIdParam, dataStoreParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateDemo - Update a Demonstration setting
func (c *ResourceApiApiController) UpdateDemo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	demoIdParam := params["demoId"]

	demoParam := Demo{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&demoParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertDemoRequired(demoParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateDemo(r.Context(), demoIdParam, demoParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateLinter - Update a Linter
func (c *ResourceApiApiController) UpdateLinter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	linterIdParam := params["LinterId"]

	linterResourceParam := LinterResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&linterResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertLinterResourceRequired(linterResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateLinter(r.Context(), linterIdParam, linterResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdatePollingProfile - Update a Polling Profile
func (c *ResourceApiApiController) UpdatePollingProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pollingProfileIdParam := params["pollingProfileId"]

	pollingProfileParam := PollingProfile{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&pollingProfileParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertPollingProfileRequired(pollingProfileParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdatePollingProfile(r.Context(), pollingProfileIdParam, pollingProfileParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateTest - update test
func (c *ResourceApiApiController) UpdateTest(w http.ResponseWriter, r *http.Request) {
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

// UpdateTestSuite - update TestSuite
func (c *ResourceApiApiController) UpdateTestSuite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testSuiteIdParam := params["testSuiteId"]

	testSuiteResourceParam := TestSuiteResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testSuiteResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestSuiteResourceRequired(testSuiteResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateTestSuite(r.Context(), testSuiteIdParam, testSuiteResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// UpdateVariableSet - Update a VariableSet
func (c *ResourceApiApiController) UpdateVariableSet(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	variableSetIdParam := params["variableSetId"]

	variableSetResourceParam := VariableSetResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&variableSetResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertVariableSetResourceRequired(variableSetResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateVariableSet(r.Context(), variableSetIdParam, variableSetResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
