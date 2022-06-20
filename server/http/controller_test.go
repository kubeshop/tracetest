package http_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	exampleRun = model.Run{
		ID:      http.IDGen.UUID(),
		TraceID: http.IDGen.TraceID(),
		Trace: &traces.Trace{
			ID: http.IDGen.TraceID(),
			RootSpan: traces.Span{
				ID:   http.IDGen.SpanID(),
				Name: "POST /pokemon/import",
				Attributes: traces.Attributes{
					"tracetest.span.type": "http",
					"service.name":        "pokeshop",
					"http.response.body":  `{"id":52}`,
				},
			},
		},
	}
)

// https://github.com/kubeshop/tracetest/issues/617
func TestContains_Issue617(t *testing.T) {

	definition := openapi.TestDefinition{
		Definitions: []openapi.TestDefinitionDefinitions{
			{
				Selector: openapi.Selector{
					Query: `span[tracetest.span.type = "http" service.name = "pokeshop"  name = "POST /pokemon/import"]`,
				},
				Assertions: []openapi.Assertion{
					{
						Attribute:  "http.response.body",
						Comparator: "contains",
						Expected:   "52",
					},
				},
			},
		},
	}

	expected := openapi.AssertionResults{
		AllPassed: true,
		Results: []openapi.AssertionResultsResults{
			{
				Selector: openapi.Selector{
					Query: `span[tracetest.span.type = "http" service.name = "pokeshop"  name = "POST /pokemon/import"]`,
					Structure: []openapi.SpanSelector{
						{
							Filters: []openapi.SelectorFilter{
								{
									Property: "tracetest.span.type",
									Operator: "=",
									Value:    "http",
								},
								{
									Property: "service.name",
									Operator: "=",
									Value:    "pokeshop",
								},
								{
									Property: "name",
									Operator: "=",
									Value:    "POST /pokemon/import",
								},
							},
						},
					},
				},
				Results: []openapi.AssertionResult{
					{
						AllPassed: true,
						Assertion: openapi.Assertion{
							Attribute:  "http.response.body",
							Comparator: "contains",
							Expected:   "52",
						},
						SpanResults: []openapi.AssertionSpanResult{
							{
								SpanId:        exampleRun.Trace.RootSpan.ID.String(),
								ObservedValue: `{"id":52}`,
								Passed:        true,
								Error:         "",
							},
						},
					},
				},
			},
		},
	}

	f := setupController(t)
	f.expectGetRun(exampleRun)

	actual, err := f.c.DryRunAssertion(context.TODO(), "", exampleRun.ID.String(), definition)
	require.NoError(t, err)

	assert.Equal(t, 200, actual.Code)
	assert.Equal(t, expected, actual.Body)
}

func setupController(t *testing.T) controllerFixture {
	mdb := new(testdb.MockRepository)
	mdb.Test(t)
	return controllerFixture{
		db: mdb,
		c:  http.NewController(mdb, nil, nil),
	}
}

type controllerFixture struct {
	db *testdb.MockRepository
	c  openapi.ApiApiServicer
}

func (f controllerFixture) expectGetRun(r model.Run) {
	f.db.
		On("GetRun", r.ID).
		Return(r, nil)
}
