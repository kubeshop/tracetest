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
			"CreateEnvironment",
			strings.ToUpper("Post"),
			"/api/environments",
			c.CreateEnvironment,
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
			"CreateTransaction",
			strings.ToUpper("Post"),
			"/api/transactions",
			c.CreateTransaction,
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
			"DeleteEnvironment",
			strings.ToUpper("Delete"),
			"/api/environments/{environmentId}",
			c.DeleteEnvironment,
		},
		{
			"DeleteLinter",
			strings.ToUpper("Delete"),
			"/api/linters/{LinterId}",
			c.DeleteLinter,
		},
		{
			"DeleteTransaction",
			strings.ToUpper("Delete"),
			"/api/transactions/{transactionId}",
			c.DeleteTransaction,
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
			"GetEnvironment",
			strings.ToUpper("Get"),
			"/api/environments/{environmentId}",
			c.GetEnvironment,
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
			"GetTests",
			strings.ToUpper("Get"),
			"/api/tests",
			c.GetTests,
		},
		{
			"GetTransaction",
			strings.ToUpper("Get"),
			"/api/transactions/{transactionId}",
			c.GetTransaction,
		},
		{
			"GetTransactions",
			strings.ToUpper("Get"),
			"/api/transactions",
			c.GetTransactions,
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
			"ListEnvironments",
			strings.ToUpper("Get"),
			"/api/environments",
			c.ListEnvironments,
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
			"UpdateEnvironment",
			strings.ToUpper("Put"),
			"/api/environments/{environmentId}",
			c.UpdateEnvironment,
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
			"UpdateTransaction",
			strings.ToUpper("Put"),
			"/api/transactions/{transactionId}",
			c.UpdateTransaction,
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

// CreateEnvironment - Create an environment
func (c *ResourceApiApiController) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	environmentResourceParam := EnvironmentResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&environmentResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertEnvironmentResourceRequired(environmentResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateEnvironment(r.Context(), environmentResourceParam)
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
	testResourceParam := TestResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&testResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTestResourceRequired(testResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateTest(r.Context(), testResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// CreateTransaction - Create new transaction
func (c *ResourceApiApiController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	transactionResourceParam := TransactionResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&transactionResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTransactionResourceRequired(transactionResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.CreateTransaction(r.Context(), transactionResourceParam)
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

// DeleteEnvironment - Delete an environment
func (c *ResourceApiApiController) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	environmentIdParam := params["environmentId"]

	result, err := c.service.DeleteEnvironment(r.Context(), environmentIdParam)
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

// DeleteTransaction - delete a transaction
func (c *ResourceApiApiController) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionIdParam := params["transactionId"]

	result, err := c.service.DeleteTransaction(r.Context(), transactionIdParam)
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

// GetEnvironment - Get a specific environment
func (c *ResourceApiApiController) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	environmentIdParam := params["environmentId"]

	result, err := c.service.GetEnvironment(r.Context(), environmentIdParam)
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

// GetTransaction - get transaction
func (c *ResourceApiApiController) GetTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionIdParam := params["transactionId"]

	result, err := c.service.GetTransaction(r.Context(), transactionIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}

// GetTransactions - Get transactions
func (c *ResourceApiApiController) GetTransactions(w http.ResponseWriter, r *http.Request) {
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
	result, err := c.service.GetTransactions(r.Context(), takeParam, skipParam, queryParam, sortByParam, sortDirectionParam)
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

// ListEnvironments - List environments
func (c *ResourceApiApiController) ListEnvironments(w http.ResponseWriter, r *http.Request) {
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
	result, err := c.service.ListEnvironments(r.Context(), takeParam, skipParam, sortByParam, sortDirectionParam)
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

// UpdateEnvironment - Update an environment
func (c *ResourceApiApiController) UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	environmentIdParam := params["environmentId"]

	environmentResourceParam := EnvironmentResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&environmentResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertEnvironmentResourceRequired(environmentResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateEnvironment(r.Context(), environmentIdParam, environmentResourceParam)
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

// UpdateTransaction - update transaction
func (c *ResourceApiApiController) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	transactionIdParam := params["transactionId"]

	transactionResourceParam := TransactionResource{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&transactionResourceParam); err != nil {
		c.errorHandler(w, r, &ParsingError{Err: err}, nil)
		return
	}
	if err := AssertTransactionResourceRequired(transactionResourceParam); err != nil {
		c.errorHandler(w, r, err, nil)
		return
	}
	result, err := c.service.UpdateTransaction(r.Context(), transactionIdParam, transactionResourceParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
