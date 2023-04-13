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
			"/api/demo",
			c.CreateDemo,
		},
		{
			"DeleteDemo",
			strings.ToUpper("Delete"),
			"/api/demo/{demoId}",
			c.DeleteDemo,
		},
		{
			"GetConfiguration",
			strings.ToUpper("Get"),
			"/api/configs/{configId}",
			c.GetConfiguration,
		},
		{
			"GetDemo",
			strings.ToUpper("Get"),
			"/api/demo/{demoId}",
			c.GetDemo,
		},
		{
			"GetPollingProfile",
			strings.ToUpper("Get"),
			"/api/pollingprofile/{pollingProfileId}",
			c.GetPollingProfile,
		},
		{
			"ListDemos",
			strings.ToUpper("Get"),
			"/api/demo",
			c.ListDemos,
		},
		{
			"UpdateConfiguration",
			strings.ToUpper("Put"),
			"/api/configs/{configId}",
			c.UpdateConfiguration,
		},
		{
			"UpdateDemo",
			strings.ToUpper("Put"),
			"/api/demo/{demoId}",
			c.UpdateDemo,
		},
		{
			"UpdatePollingProfile",
			strings.ToUpper("Put"),
			"/api/pollingprofile/{pollingProfileId}",
			c.UpdatePollingProfile,
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
