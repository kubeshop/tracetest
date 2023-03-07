package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Operation string

type operationTester struct {
	name     Operation
	getSteps func(*testing.T, ResourceTypeTest) []operationTesterStep
}

type operationTesterStep struct {
	buildRequest   func(*testing.T, *httptest.Server, contentTypeConverter, ResourceTypeTest) *http.Request
	assertResponse func(*testing.T, *http.Response, contentTypeConverter, ResourceTypeTest)
	postAssert     func(*testing.T, contentTypeConverter, ResourceTypeTest, *httptest.Server)
}

type singleStepOperationTester struct {
	name           Operation
	buildRequest   func(*testing.T, *httptest.Server, contentTypeConverter, ResourceTypeTest) *http.Request
	assertResponse func(*testing.T, *http.Response, contentTypeConverter, ResourceTypeTest)
	postAssert     func(*testing.T, contentTypeConverter, ResourceTypeTest, *httptest.Server)
}

func buildSingleStepOperation(operation singleStepOperationTester) operationTester {
	return operationTester{
		name: operation.name,
		getSteps: func(t *testing.T, rt ResourceTypeTest) []operationTesterStep {
			return []operationTesterStep{
				{
					buildRequest:   operation.buildRequest,
					assertResponse: operation.assertResponse,
					postAssert:     operation.postAssert,
				},
			}
		},
	}
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
		listPaginatedSuccessOperation,
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
