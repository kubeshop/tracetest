package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Operation string

type operationTester struct {
	buildRequest   func(*testing.T, *httptest.Server, contentTypeConverter, ResourceTypeTest) *http.Request
	assertResponse func(*testing.T, *http.Response, contentTypeConverter, ResourceTypeTest)
	name           Operation
	postAssert     func(t *testing.T, ct contentTypeConverter, rt ResourceTypeTest, testServer *httptest.Server)
}

var (
	defaultOperations = []operationTester{
		createNoIDOperation,
		createSuccessOperation,

		updateNotFoundOperation,
		updateSuccessOperation,

		getNotFoundOperation,
		getSuccessOperation,

		deleteNotFoundOperation,
		deleteSuccessOperation,

		listNoResultsOperation,
		listSuccessOperation,
		listPaginatedAscendingSuccessOperation,
		listPaginatedDescendingSuccessOperation,
		// TODO: add tests for other operations
	}

	errorOperations = []operationTester{
		createInternalErrorOperation,
		updateInternalErrorOperation,
		getInternalErrorOperation,
		deleteInternalErrorOperation,
		listInternalErrorOperation,
		listWithInvalidSortFieldOperation,
	}
)
