package testutil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kubeshop/tracetest/server/resourcemanager"
	"golang.org/x/exp/slices"
)

type Operation string

type operationTester struct {
	name               Operation
	neededForOperation resourcemanager.Operation
	getSteps           func(*testing.T, ResourceTypeTest) []operationTesterStep
}

func (ot operationTester) needsToRun(enabledOperations []resourcemanager.Operation) bool {
	if ot.neededForOperation == "" {
		panic(fmt.Errorf("operation %s does not define neededForOperation", ot.name))
	}
	return slices.Contains(enabledOperations, ot.neededForOperation)
}

type operationTesterStep struct {
	buildRequest   func(*testing.T, *httptest.Server, contentTypeConverter, ResourceTypeTest) *http.Request
	assertResponse func(*testing.T, *http.Response, contentTypeConverter, ResourceTypeTest)
	postAssert     func(*testing.T, contentTypeConverter, ResourceTypeTest, *httptest.Server)
}

type singleStepOperationTester struct {
	name               Operation
	neededForOperation resourcemanager.Operation
	buildRequest       func(*testing.T, *httptest.Server, contentTypeConverter, ResourceTypeTest) *http.Request
	assertResponse     func(*testing.T, *http.Response, contentTypeConverter, ResourceTypeTest)
	postAssert         func(*testing.T, contentTypeConverter, ResourceTypeTest, *httptest.Server)
}

func buildSingleStepOperation(operation singleStepOperationTester) operationTester {
	return operationTester{
		name:               operation.name,
		neededForOperation: operation.neededForOperation,
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
