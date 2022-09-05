package model_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestTestEncoding(t *testing.T) {
	id := uuid.MustParse("ccf94a15-e33e-4c75-ae94-d0b401c53da1")
	t1 := time.Date(2022, 06, 07, 13, 03, 24, 100, time.UTC)

	test := model.Test{
		ID:          id,
		CreatedAt:   t1,
		Name:        "the name",
		Description: "the description",
		Version:     1,
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL:    "http://localhost:8080/list",
				Method: model.HTTPMethodGET,
			},
		},
		Specs: (model.OrderedMap[model.SpanQuery, []model.Assertion]{}).
			MustAdd(model.SpanQuery(`span[name="test"]`), []model.Assertion{
				{"name", comparator.Eq, "test"},
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
	id := uuid.MustParse("ccf94a15-e33e-4c75-ae94-d0b401c53da1")
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
				ID:                 id,
				TraceID:            tid,
				SpanID:             sid,
				State:              model.RunStateFailed,
				LastError:          errors.New("some error"),
				CreatedAt:          t1,
				ServiceTriggeredAt: t1,
				CompletedAt:        t1,
				Trigger: model.Trigger{
					Type: model.TriggerTypeHTTP,
					HTTP: &model.HTTPRequest{
						Method: model.HTTPMethodPOST,
						URL:    "http://google.com",
						Headers: []model.HTTPHeader{
							{"Content-Type", "application/json"},
						},
						Body: `{"id":52}`,
					},
				},
				TestVersion: 1,
				Metadata:    map[string]string{"key": "value"},
			},
		},
		{
			name: "Success",
			run: model.Run{
				ID:                        id,
				TraceID:                   tid,
				SpanID:                    sid,
				State:                     model.RunStateFinished,
				LastError:                 nil,
				CreatedAt:                 t1,
				ServiceTriggeredAt:        t1,
				ServiceTriggerCompletedAt: t2,
				ObtainedTraceAt:           t3,
				CompletedAt:               t4,
				Trigger: model.Trigger{
					Type: model.TriggerTypeHTTP,
					HTTP: &model.HTTPRequest{
						Method: model.HTTPMethodPOST,
						URL:    "http://google.com",
						Headers: []model.HTTPHeader{
							{"Content-Type", "application/json"},
						},
						Body: `{"name":"Larry"}`,
					},
				},
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

			// assert.Equal doesn't work on time.Time vars (see https://stackoverflow.com/a/69362528)
			// We could manually compare each field, but if we add a new field to the struct but not the test,
			// the test would still pass, even if the json encoding is incorrect.
			// So instead, we can reset all date fields and compare them separately.
			// If we add a new date field, the `assert.Equal(t, run, actual)` will catch it
			expectedCreatedAt := cl.run.CreatedAt
			expectedServiceTriggeredAt := cl.run.ServiceTriggeredAt
			expectedServiceTriggerCompletedAt := cl.run.ServiceTriggerCompletedAt
			expectedObtainedTraceAt := cl.run.ObtainedTraceAt
			expectedCompletedAt := cl.run.CompletedAt

			run.CreatedAt = t1
			run.ServiceTriggeredAt = t1
			run.ServiceTriggerCompletedAt = t1
			run.ObtainedTraceAt = t1
			run.CompletedAt = t1
			if run.Trace != nil {
				run.Trace.RootSpan.StartTime = t1
				run.Trace.RootSpan.EndTime = t1
				run.Trace.Flat[sid].StartTime = t1
				run.Trace.Flat[sid].EndTime = t1
			}

			actualCreatedAt := actual.CreatedAt
			actualServiceTriggeredAt := actual.ServiceTriggeredAt
			actualServiceTriggerCompletedAt := actual.ServiceTriggerCompletedAt
			actualObtainedTraceAt := actual.ObtainedTraceAt
			actualCompletedAt := actual.CompletedAt

			actual.CreatedAt = t1
			actual.ServiceTriggeredAt = t1
			actual.ServiceTriggerCompletedAt = t1
			actual.ObtainedTraceAt = t1
			actual.CompletedAt = t1
			if actual.Trace != nil {
				actual.Trace.RootSpan.StartTime = t1
				actual.Trace.RootSpan.EndTime = t1
				actual.Trace.Flat[sid].StartTime = t1
				actual.Trace.Flat[sid].EndTime = t1
			}

			assert.Equal(t, run, actual)

			assert.WithinDuration(t, expectedCreatedAt, actualCreatedAt, 0)
			assert.WithinDuration(t, expectedServiceTriggeredAt, actualServiceTriggeredAt, 0)
			assert.WithinDuration(t, expectedServiceTriggerCompletedAt, actualServiceTriggerCompletedAt, 0)
			assert.WithinDuration(t, expectedObtainedTraceAt, actualObtainedTraceAt, 0)
			assert.WithinDuration(t, expectedCompletedAt, actualCompletedAt, 0)
		})
	}
}
