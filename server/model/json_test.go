package model_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/modeltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestTestEncoding(t *testing.T) {
	t1 := time.Date(2022, 06, 07, 13, 03, 24, 100, time.UTC)

	test := model.Test{
		ID:          id.ID("test1"),
		CreatedAt:   t1,
		Name:        "the name",
		Description: "the description",
		Version:     1,
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    "http://localhost:11633/list",
				Method: model.HTTPMethodGET,
			},
		},
		Specs: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).
			MustAdd(model.SpanQuery(`span[name="test"]`), model.NamedAssertions{
				Name: "test",
				Assertions: []model.Assertion{
					model.Assertion(`attr:name = "test"`),
				},
			}),
	}

	encoded, err := json.Marshal(test)
	require.NoError(t, err)

	var actual model.Test
	err = json.Unmarshal(encoded, &actual)
	require.NoError(t, err)

	assert.Equal(t, test, actual)
}

func TestRunEncoding(t *testing.T) {
	tid, _ := trace.TraceIDFromHex("83c7f2fb8b556416e12e1d18c05a30c3")
	sid, _ := trace.SpanIDFromHex("9ed1382a48be2649")

	t1 := time.Date(2022, 06, 07, 13, 03, 24, 100, time.UTC)
	t2 := time.Date(2022, 06, 07, 13, 03, 25, 100, time.UTC)
	t3 := time.Date(2022, 06, 07, 13, 03, 27, 100, time.UTC)
	t4 := time.Date(2022, 06, 07, 13, 03, 28, 100, time.UTC)

	rootSpan := model.Span{
		ID:        sid,
		Name:      "Root Span",
		StartTime: t1,
		EndTime:   t2,
		Attributes: model.Attributes{
			"tracetest.span.duration": "200",
			"tracetest.span.type":     "http",
		},
	}
	exampleTrace := &model.Trace{
		ID:       tid,
		RootSpan: rootSpan,
		Flat: map[trace.SpanID]*model.Span{
			sid: &rootSpan,
		},
	}

	cases := []struct {
		name string
		run  model.Run
	}{
		{
			name: "Errors",
			run: model.Run{
				ID:                 1,
				TraceID:            tid,
				SpanID:             sid,
				State:              model.RunStateTriggerFailed,
				LastError:          errors.New("some error"),
				CreatedAt:          t1,
				ServiceTriggeredAt: t1,
				CompletedAt:        t1,
				TestVersion:        1,
				Metadata:           map[string]string{"key": "value"},
			},
		},
		{
			name: "Success",
			run: model.Run{
				ID:                        1,
				TraceID:                   tid,
				SpanID:                    sid,
				State:                     model.RunStateFinished,
				LastError:                 nil,
				CreatedAt:                 t1,
				ServiceTriggeredAt:        t1,
				ServiceTriggerCompletedAt: t2,
				ObtainedTraceAt:           t3,
				CompletedAt:               t4,
				TriggerResult: model.TriggerResult{
					Type: model.TriggerTypeHTTP,
					HTTP: &model.HTTPResponse{
						Status:     "OK",
						StatusCode: 200,
						Headers: []model.HTTPHeader{
							{"Content-Type", "application/json"},
							{"Length", "9"},
						},
						Body: `{"id":52}`,
					},
				},
				Trace:       exampleTrace,
				TestVersion: 2,
				Metadata:    map[string]string{"another_key": "value"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			run := cl.run

			encoded, err := json.Marshal(run)
			require.NoError(t, err)

			var actual model.Run
			err = json.Unmarshal(encoded, &actual)
			require.NoError(t, err)

			modeltest.AssertRunEqual(t, cl.run, actual)
		})
	}
}

func TestOldAssertionSpecsFormatWithoutNames(t *testing.T) {
	type OldTest struct {
		ID               id.ID
		CreatedAt        time.Time
		Name             string
		Description      string
		Version          int
		ServiceUnderTest model.Trigger
		Specs            model.OrderedMap[model.SpanQuery, []model.Assertion]
		Summary          model.Summary
	}

	expectedSpecs := model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}
	expectedSpecs = expectedSpecs.MustAdd(model.SpanQuery(`span[tracetest.span.type = "http"]`), model.NamedAssertions{
		Name: "",
		Assertions: []model.Assertion{
			model.Assertion(`attr:http.status = 200`),
		},
	})

	specs := model.OrderedMap[model.SpanQuery, []model.Assertion]{}
	specs = specs.MustAdd(model.SpanQuery(`span[tracetest.span.type = "http"]`), []model.Assertion{
		model.Assertion(`attr:http.status = 200`),
	})
	oldTest := OldTest{
		ID:               id.NewRandGenerator().ID(),
		CreatedAt:        time.Now(),
		Name:             "my test name",
		Description:      "this is an old test using the old test format from version <= 0.7.2",
		Version:          1,
		ServiceUnderTest: model.Trigger{},
		Specs:            specs,
		Summary:          model.Summary{},
	}

	oldTestJson, err := json.Marshal(oldTest)
	require.NoError(t, err)

	var test model.Test
	err = json.Unmarshal(oldTestJson, &test)
	require.NoError(t, err)

	assert.Equal(t, expectedSpecs, test.Specs)
}

func TestNewAssertionSpecFormat(t *testing.T) {
	test := model.Test{
		ID:               id.NewRandGenerator().ID(),
		CreatedAt:        time.Now(),
		Name:             gofakeit.Name(),
		Description:      gofakeit.AdjectiveDescriptive(),
		Version:          1,
		ServiceUnderTest: model.Trigger{},
		Specs: model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}.MustAdd(
			model.SpanQuery(`span[tracetest.span.type = "http"`), model.NamedAssertions{
				Name: "my test",
				Assertions: []model.Assertion{
					model.Assertion(`attr:http.status = 200`),
				},
			},
		),
		Summary: model.Summary{},
	}

	bytes, err := json.Marshal(test)
	require.NoError(t, err)

	var newTest model.Test
	err = json.Unmarshal(bytes, &newTest)

	require.NoError(t, err)
	assert.Equal(t, test.Specs, newTest.Specs)
}
