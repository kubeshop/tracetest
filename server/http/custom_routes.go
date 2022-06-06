package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/openapi"
	"go.opentelemetry.io/otel/attribute"
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
		ctx, span := c.tracer.Start(r.Context(), name)
		defer span.End()
		span.SetAttributes(
			attribute.String(string(semconv.HTTPMethodKey), r.Method),
			attribute.String(string(semconv.HTTPRouteKey), route),
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
