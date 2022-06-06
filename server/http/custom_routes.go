package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/openapi"
)

func NewCustomController(s openapi.ApiApiServicer, r openapi.Router, eh openapi.ErrorHandler) openapi.Router {
	return &customController{s, r, eh}
}

type customController struct {
	service      openapi.ApiApiServicer
	router       openapi.Router
	errorHandler openapi.ErrorHandler
}

func (c *customController) Routes() openapi.Routes {

	routes := c.router.Routes()

	routes[c.getRouteIndex("GetRunResultJUnit")].HandlerFunc = c.GetRunResultJUnit

	return routes
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
