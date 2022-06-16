package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

	for index, route := range routes {
		routeName := fmt.Sprintf("%s %s", route.Method, route.Pattern)
		newRouteHandlerFunc := c.instrumentRoute(routeName, route.Pattern, route.HandlerFunc)
		route.HandlerFunc = newRouteHandlerFunc

		routes[index] = route
	}

	return routes
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

// GetRunResultJUnit - get test run results in JUnit xml format
func (c *customController) GetRunResultJUnit(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testIdParam := params["testId"]

	runIdParam := params["runId"]

	result, err := c.service.GetRunResultJUnit(r.Context(), testIdParam, runIdParam)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, &result)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=UTF-8")
	w.Write(result.Body.([]byte))
}
