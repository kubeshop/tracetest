package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/stretchr/testify/assert"
)

type dummyHandler struct {
	t         *testing.T
	OnRequest func(r *http.Request)
}

func (d *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.OnRequest(r)
}

func TestMiddleware(t *testing.T) {
	dummyHandler := &dummyHandler{t: t}
	nextHandler := http.HandlerFunc(dummyHandler.ServeHTTP)
	handlerToTest := middleware.TenantMiddleware(nextHandler)

	t.Run("should set the tenant id in the context", func(t *testing.T) {
		uuid := "16700d36-8e0a-4169-9bb8-30a249281841"
		onRequest := func(r *http.Request) {
			assert.Equal(t, uuid, r.Context().Value(middleware.TenantIDKey).(string))
		}

		dummyHandler.OnRequest = onRequest

		req := httptest.NewRequest("GET", "http://testing", nil)
		req.Header.Set(middleware.HeaderTenantID, uuid)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})

	t.Run("should set the tenant id as empty string", func(t *testing.T) {
		onRequest := func(r *http.Request) {
			assert.Equal(t, "", r.Context().Value(middleware.TenantIDKey).(string))
		}

		dummyHandler.OnRequest = onRequest

		req := httptest.NewRequest("GET", "http://testing", nil)
		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	})
}
