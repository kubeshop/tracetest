package http

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/openapi"
)

type PrefixedRouter struct {
	prefix string
	router openapi.Router
}

var _ openapi.Router = PrefixedRouter{}

func NewPrefixedRouter(router openapi.Router, prefix string) openapi.Router {
	prefixRouter := PrefixedRouter{
		prefix: prefix,
		router: router,
	}

	return prefixRouter
}

func (r PrefixedRouter) Routes() openapi.Routes {
	routes := r.router.Routes()
	prefixedRoutes := make(openapi.Routes, 0, len(routes))
	for _, route := range routes {
		prefixedRoute := route
		prefixedRoute.Pattern = fmt.Sprintf("%s%s", r.prefix, route.Pattern)
		prefixedRoutes = append(prefixedRoutes, prefixedRoute)
	}

	return prefixedRoutes
}
