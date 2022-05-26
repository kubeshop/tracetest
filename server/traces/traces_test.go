package traces_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestJSONEncoding(t *testing.T) {

	rootSpan := createSpan("root")
	subSpan1 := createSpan("subSpan1")
	subSubSpan1 := createSpan("subSubSpan1")
	subSpan2 := createSpan("subSpan2")

	rootSpan.Children = []*traces.Span{subSpan1, subSpan2}
	subSpan1.Parent = rootSpan
	subSpan2.Parent = rootSpan

	subSpan1.Children = []*traces.Span{subSubSpan1}
	subSubSpan1.Parent = subSpan1

	tid := id.NewRandGenerator().TraceID()
	trace := traces.Trace{
		ID:       tid,
		RootSpan: *rootSpan,
		Flat: map[trace.SpanID]*traces.Span{
			rootSpan.ID:    rootSpan,
			subSpan1.ID:    subSpan1,
			subSubSpan1.ID: subSubSpan1,
			subSpan2.ID:    subSpan2,
		},
	}

	jsonEncoded := fmt.Sprintf(`{
		"ID": "%s",
		"RootSpan": {
			"ID": "%s",
			"Name":"root",
			"StartTime":"%s",
			"EndTime":"%s",
			"Attributes": {
				"service.name": "root"
			},
			"Children": [
				{
					"ID": "%s",
					"Name":"subSpan1",
					"StartTime":"%s",
					"EndTime":"%s",
					"Attributes": {
						"service.name": "subSpan1"
					},
					"Children": [
						{
							"ID": "%s",
							"StartTime":"%s",
							"EndTime":"%s",
							"Name":"subSubSpan1",
							"Attributes": {
								"service.name": "subSubSpan1"
							},
							"Children": []
						}
					]
				},
				{
					"ID": "%s",
					"Name":"subSpan2",
					"StartTime":"%s",
					"EndTime":"%s",
					"Attributes": {
						"service.name": "subSpan2"
					},
					"Children": []
				}
			]
		}
	}`,
		tid.String(),

		rootSpan.ID.String(),
		rootSpan.StartTime.Format(time.RFC3339),
		rootSpan.EndTime.Format(time.RFC3339),

		subSpan1.ID.String(),
		subSpan1.StartTime.Format(time.RFC3339),
		subSpan1.EndTime.Format(time.RFC3339),

		subSubSpan1.ID.String(),
		subSubSpan1.StartTime.Format(time.RFC3339),
		subSubSpan1.EndTime.Format(time.RFC3339),

		subSpan2.ID.String(),
		subSpan2.StartTime.Format(time.RFC3339),
		subSpan2.EndTime.Format(time.RFC3339),
	)

	t.Run("encode", func(t *testing.T) {
		actual, err := json.Marshal(&trace)
		require.NoError(t, err)

		assert.JSONEq(t, jsonEncoded, string(actual))
	})

	t.Run("decode", func(t *testing.T) {
		var actual traces.Trace
		err := json.Unmarshal([]byte(jsonEncoded), &actual)
		require.NoError(t, err)

		fmt.Printf("%+v\n", actual.RootSpan.Children[0].Children[0])

		assert.Equal(t, trace, actual)
	})
}

func createSpan(name string) *traces.Span {
	return &traces.Span{
		ID:        id.NewRandGenerator().SpanID(),
		Name:      name,
		StartTime: time.Date(2021, 11, 24, 14, 05, 12, 0, time.UTC),
		EndTime:   time.Date(2021, 11, 24, 14, 05, 17, 0, time.UTC),
		Attributes: traces.Attributes{
			"service.name": name,
		},
	}
}
