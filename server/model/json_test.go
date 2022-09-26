package model_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/modeltest"
	"github.com/kubeshop/tracetest/server/traces"
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
		Specs: (model.OrderedMap[model.SpanQuery, []model.Assertion]{}).
			MustAdd(model.SpanQuery(`span[name="test"]`), []model.Assertion{
				{"name", comparator.Eq, createExpressionFromString("test")},
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

	rootSpan := traces.Span{
		ID:        sid,
		Name:      "Root Span",
		StartTime: t1,
		EndTime:   t2,
		Attributes: traces.Attributes{
			"tracetest.span.duration": "200",
			"tracetest.span.type":     "http",
		},
	}
	exampleTrace := &traces.Trace{
		ID:       tid,
		RootSpan: rootSpan,
		Flat: map[trace.SpanID]*traces.Span{
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
				State:              model.RunStateFailed,
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

func TestOldAssertionSerialization(t *testing.T) {
	var newAssertion model.Assertion
	oldAssertion := struct {
		Attribute  string
		Comparator string
		Value      string
	}{
		Attribute:  "tracetest.span.type",
		Comparator: "=",
		Value:      "http",
	}

	oldAssertionJson, err := json.Marshal(oldAssertion)
	require.NoError(t, err)

	err = json.Unmarshal(oldAssertionJson, &newAssertion)
	require.NoError(t, err)

	assert.Equal(t, oldAssertion.Attribute, newAssertion.Attribute.String())
	assert.Equal(t, oldAssertion.Comparator, newAssertion.Comparator.String())
	assert.Equal(t, oldAssertion.Value, newAssertion.Value.String())
}
