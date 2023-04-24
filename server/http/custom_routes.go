package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

func NewCustomController(s openapi.ApiApiServicer, r openapi.Router, eh openapi.ErrorHandler, t trace.Tracer) openapi.Router {
	return &customController{s, r, eh, t}
}

type customController struct {
	service      openapi.ApiApiServicer
	router       openapi.Router
	errorHandler openapi.ErrorHandler
	tracer       trace.Tracer
}

func (c *customController) Routes() openapi.Routes {

	routes := c.router.Routes()

	routes[c.getRouteIndex("GetRunResultJUnit")].HandlerFunc = c.GetRunResultJUnit
	routes[c.getRouteIndex("GetTestVersionDefinitionFile")].HandlerFunc = c.GetTestVersionDefinitionFile

	routes[c.getRouteIndex("GetTestRuns")].HandlerFunc = c.GetTestRuns

	routes[c.getRouteIndex("GetTests")].HandlerFunc = paginatedEndpoint[openapi.Test](c.service.GetTests, c.errorHandler)
	routes[c.getRouteIndex("GetResources")].HandlerFunc = paginatedEndpoint[openapi.Resource](c.service.GetResources, c.errorHandler)

	for index, route := range routes {
		routeName := fmt.Sprintf("%s %s", route.Method, route.Pattern)
		hf := route.HandlerFunc

		hf = c.instrumentRoute(routeName, route.Pattern, hf)
		hf = c.analytics(route.Name, hf)

		route.HandlerFunc = hf

		routes[index] = route
	}

	return routes
}

// GetTestRuns - get the runs for a test
func (c *customController) GetTestRuns(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	testIdParam := params["testId"]

	takeParam, err := parseInt32Parameter(query.Get("take"), false)
	if err != nil {
		c.errorHandler(w, r, &openapi.ParsingError{Err: err}, nil)
		return
	}
	skipParam, err := parseInt32Parameter(query.Get("skip"), false)
	if err != nil {
		c.errorHandler(w, r, &openapi.ParsingError{Err: err}, nil)
		return
	}
	result, err := c.service.GetTestRuns(r.Context(), testIdParam, takeParam, skipParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	res := result.Body.(paginated[openapi.TestRun])

	w.Header().Set("X-Total-Count", strconv.Itoa(res.count))
	openapi.EncodeJSONResponse(res.items, &result.Code, w)
}

// GetRunResultJUnit - get test run results in JUnit xml format
func (c *customController) GetRunResultJUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]
	runIdParam, err := strconv.Atoi(params["runId"])
	if err != nil {
		c.errorHandler(w, r, fmt.Errorf("could not convert runId to integer: %w", err), nil)
	}

	result, err := c.service.GetRunResultJUnit(r.Context(), testIdParam, int32(runIdParam))
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	w.Write(result.Body.([]byte))
}

// GetTestVersionDefinitionFile - Get the test definition as an YAML file
func (c *customController) GetTestVersionDefinitionFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	versionParam, err := parseInt32Parameter(params["version"], true)
	if err != nil {
		c.errorHandler(w, r, &openapi.ParsingError{Err: err}, nil)
		return
	}

	result, err := c.service.GetTestVersionDefinitionFile(r.Context(), testIdParam, versionParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=UTF-8")
	w.Write(result.Body.([]byte))
}

func paginatedEndpoint[T any](
	f func(c context.Context, take, skip int32, query string, sortBy string, sortDirection string) (openapi.ImplResponse, error),
	errorHandler openapi.ErrorHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		takeParam, err := parseInt32Parameter(query.Get("take"), false)
		if err != nil {
			errorHandler(w, r, &openapi.ParsingError{Err: err}, nil)
			return
		}
		skipParam, err := parseInt32Parameter(query.Get("skip"), false)
		if err != nil {
			errorHandler(w, r, &openapi.ParsingError{Err: err}, nil)
			return
		}
		queryParam := query.Get("query")
		sortByParam := query.Get("sortBy")
		sortDirectionParam := query.Get("sortDirection")
		result, err := f(r.Context(), takeParam, skipParam, queryParam, sortByParam, sortDirectionParam)
		// If an error occurred, encode the error with the status code
		if err != nil {
			errorHandler(w, r, err, &result)
			return
		}
		res := result.Body.(paginated[T])

		w.Header().Set("X-Total-Count", strconv.Itoa(res.count))
		openapi.EncodeJSONResponse(res.items, &result.Code, w)
	}
}

const errMsgRequiredMissing = "required parameter is missing"

func parseInt32Parameter(param string, required bool) (int32, error) {
	if param == "" {
		if required {
			return 0, errors.New(errMsgRequiredMissing)
		}

		return 0, nil
	}

	val, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return -1, err
	}

	return int32(val), nil
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toWords(str string) string {
	words := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	words = matchAllCap.ReplaceAllString(words, "${1} ${2}")
	return words
}

func (c *customController) analytics(name string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		machineID := r.Header.Get("x-client-id")
		analytics.SendEvent(toWords(name), "test", machineID, nil)

		f(w, r)
	}
}

func (c *customController) instrumentRoute(name string, route string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := c.tracer.Start(ctx, name)
		defer span.End()

		params := make(map[string]interface{}, 0)
		for key, value := range mux.Vars(r) {
			params[key] = value
		}

		paramsJson, _ := json.Marshal(params)

		queryString := make(map[string]interface{}, 0)
		for key, value := range r.URL.Query() {
			queryString[key] = value
		}
		queryStringJson, _ := json.Marshal(queryString)

		headers := make(map[string]interface{}, 0)
		for key, value := range r.Header {
			headers[key] = value
		}
		headersJson, _ := json.Marshal(headers)

		span.SetAttributes(
			attribute.String(string(semconv.HTTPMethodKey), r.Method),
			attribute.String(string(semconv.HTTPRouteKey), route),
			attribute.String(string(semconv.HTTPTargetKey), r.URL.String()),
			attribute.String("http.request.params", string(paramsJson)),
			attribute.String("http.request.query", string(queryStringJson)),
			attribute.String("http.request.headers", string(headersJson)),
		)

		newRequest := r.WithContext(ctx)

		f(w, newRequest)
	}
}

func (c *customController) getRouteIndex(key string) int {
	routes := (&openapi.ApiApiController{}).Routes()
	for i, r := range routes {
		if r.Name == key {
			return i
		}
	}

	panic(fmt.Errorf(`route "%s" not found`, key))
}
