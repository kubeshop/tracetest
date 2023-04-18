package traces_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestJSONEncoding(t *testing.T) {

	rootSpan := createSpan("root")
	subSpan1 := createSpan("subSpan1")
	subSubSpan1 := createSpan("subSubSpan1")
	subSpan2 := createSpan("subSpan2")

	subSpan1.Parent = rootSpan
	subSpan2.Parent = rootSpan
	subSubSpan1.Parent = subSpan1

	flat := map[trace.SpanID]*model.Span{
		// We copy those spans so they don't have children injected into them
		// the flat structure shouldn't have children.
		rootSpan.ID:    copySpan(rootSpan),
		subSpan1.ID:    copySpan(subSpan1),
		subSubSpan1.ID: copySpan(subSubSpan1),
		subSpan2.ID:    copySpan(subSpan2),
	}

	rootSpan.Children = []*model.Span{subSpan1, subSpan2}
	subSpan1.Children = []*model.Span{subSubSpan1}

	tid := id.NewRandGenerator().TraceID()
	trace := model.Trace{
		ID:       tid,
		RootSpan: *rootSpan,
		Flat:     flat,
	}

	jsonEncoded := fmt.Sprintf(`{
		"ID": "%s",
		"RootSpan": {
			"ID": "%s",
			"Name":"root",
			"StartTime":"%d",
			"EndTime":"%d",
			"Attributes": {
				"service.name": "root"
			},
			"Children": [
				{
					"ID": "%s",
					"Name":"subSpan1",
					"StartTime":"%d",
					"EndTime":"%d",
					"Attributes": {
						"service.name": "subSpan1"
					},
					"Children": [
						{
							"ID": "%s",
							"StartTime":"%d",
							"EndTime":"%d",
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
					"StartTime":"%d",
					"EndTime":"%d",
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
		rootSpan.StartTime.UnixMilli(),
		rootSpan.EndTime.UnixMilli(),

		subSpan1.ID.String(),
		subSpan1.StartTime.UnixMilli(),
		subSpan1.EndTime.UnixMilli(),

		subSubSpan1.ID.String(),
		subSubSpan1.StartTime.UnixMilli(),
		subSubSpan1.EndTime.UnixMilli(),

		subSpan2.ID.String(),
		subSpan2.StartTime.UnixMilli(),
		subSpan2.EndTime.UnixMilli(),
	)

	t.Run("encode", func(t *testing.T) {
		actual, err := json.Marshal(&trace)
		require.NoError(t, err)

		assert.JSONEq(t, jsonEncoded, string(actual))
	})

	t.Run("decode", func(t *testing.T) {
		var actual model.Trace
		err := json.Unmarshal([]byte(jsonEncoded), &actual)
		require.NoError(t, err)

		// I've added more specific validations to be easier to find where the problem is
		require.Equal(t, trace.ID, actual.ID)
		require.Equal(t, trace.RootSpan, actual.RootSpan)
		require.Equal(t, trace.Flat, actual.Flat)

		// I left this as a guarantee we won't forget to change this test in case we add
		// another attribute to our traces.
		assert.Equal(t, trace, actual)
	})
}

func copySpan(span *model.Span) *model.Span {
	newSpan := *span
	return &newSpan
}

func createSpan(name string) *model.Span {
	return &model.Span{
		ID:        id.NewRandGenerator().SpanID(),
		Name:      name,
		StartTime: time.Date(2021, 11, 24, 14, 05, 12, 0, time.UTC),
		EndTime:   time.Date(2021, 11, 24, 14, 05, 17, 0, time.UTC),
		Attributes: model.Attributes{
			"service.name": name,
		},
		Children: nil,
	}
}

var now = time.Now()

func getTime(n int) time.Time {
	return now.Add(time.Duration(n) * time.Second)
}

func TestSort(t *testing.T) {
	randGenerator := id.NewRandGenerator()
	trace := model.Trace{
		ID: randGenerator.TraceID(),
		RootSpan: model.Span{
			Name:       "root",
			StartTime:  getTime(0),
			Attributes: model.Attributes{},
			Children: []*model.Span{
				{
					Name:      "child 2",
					StartTime: getTime(2),
				},
				{
					Name:      "child 3",
					StartTime: getTime(3),
				},
				{
					Name:      "child 1",
					StartTime: getTime(1),
					Children: []*model.Span{
						{
							Name:      "grandchild 1",
							StartTime: getTime(2),
						},
						{
							Name:      "grandchild 2",
							StartTime: getTime(3),
						},
					},
				},
			},
		},
	}

	sortedTrace := trace.Sort()

	expectedTrace := model.Trace{
		ID: randGenerator.TraceID(),
		RootSpan: model.Span{
			Name:       "root",
			StartTime:  getTime(0),
			Attributes: model.Attributes{},
			Children: []*model.Span{
				{
					Name:      "child 1",
					StartTime: getTime(1),
					Children: []*model.Span{
						{

							Name:      "grandchild 1",
							StartTime: getTime(2),
							Children:  make([]*model.Span, 0),
						},
						{

							Name:      "grandchild 2",
							StartTime: getTime(3),
							Children:  make([]*model.Span, 0),
						},
					},
				},
				{
					Name:      "child 2",
					StartTime: getTime(2),
					Children:  make([]*model.Span, 0),
				},
				{
					Name:      "child 3",
					StartTime: getTime(3),
					Children:  make([]*model.Span, 0),
				},
			},
		},
	}

	assert.Equal(t, expectedTrace.RootSpan, sortedTrace.RootSpan)
}
