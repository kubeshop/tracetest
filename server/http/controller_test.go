package http_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

var (
	exampleRun = model.Run{
		ID:      1,
		TestID:  id.ID("abc123"),
		TraceID: http.IDGen.TraceID(),
		Trace: &model.Trace{
			ID: http.IDGen.TraceID(),
			RootSpan: model.Span{
				ID:   http.IDGen.SpanID(),
				Name: "POST /pokemon/import",
				Attributes: model.Attributes{
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

	spec := openapi.TestSpecs{
		Specs: []openapi.TestSpecsSpecsInner{
			{
				Selector: openapi.Selector{
					Query: `span[tracetest.span.type = "http" service.name = "pokeshop"  name = "POST /pokemon/import"]`,
				},
				Assertions: []string{
					`attr:http.response.body contains 52`,
				},
			},
		},
	}

	expected := openapi.AssertionResults{
		AllPassed: true,
		Results: []openapi.AssertionResultsResultsInner{
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
						Assertion: `attr:http.response.body contains 52`,
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

	actual, err := f.c.DryRunAssertion(context.TODO(), exampleRun.TestID.String(), int32(exampleRun.ID), spec)
	require.NoError(t, err)

	assert.Equal(t, 200, actual.Code)
	assert.Equal(t, expected, actual.Body)
}

func setupController(t *testing.T) controllerFixture {
	mdb := new(testdb.MockRepository)
	mdb.Test(t)
	return controllerFixture{
		db: mdb,
		c: http.NewController(
			mdb,
			nil,
			nil,
			mappings.New(traces.NewConversionConfig(), comparator.DefaultRegistry(), mdb),
			environment.NewRepository(nil),
			&trigger.Registry{},
			trace.NewNoopTracerProvider().Tracer("tracer"),
			"unit-test",
		),
	}
}

type controllerFixture struct {
	db *testdb.MockRepository
	c  openapi.ApiApiServicer
}

func (f controllerFixture) expectGetRun(r model.Run) {
	f.db.
		On("GetRun", r.TestID, r.ID).
		Return(r, nil)
}
