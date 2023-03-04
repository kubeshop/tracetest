package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Operation string

type OperationTester interface {
	buildRequest(*testing.T, *httptest.Server, contentTypeConverter, ResourceTypeTest) *http.Request
	assertResponse(*testing.T, *http.Response, contentTypeConverter, ResourceTypeTest)
	name() Operation
	postAssert(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server)
}

var (
	defaultOperations = []OperationTester{
		createNoIDOperation{},
		createSuccessOperation{},

		updateNotFoundOperation{},
		updateSuccessOperation{},

		getNotFoundOperation{},
		getSuccessOperation{},

		deleteNotFoundOperation{},
		deleteSuccessOperation{},

		listNoResultsOperation{},
		listSuccessOperation{},
		// TODO: add tests for pagination etc
	}

	errorOperations = []OperationTester{
		createInternalErrorOperation{},
		updateInternalErrorOperation{},
		getInternalErrorOperation{},
		deleteInternalErrorOperation{},
		listInternalErrorOperation{},
	}
)
