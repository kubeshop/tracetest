package traces_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestTraces(t *testing.T) {
	rootSpan := newSpan("Root")
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []traces.Span{rootSpan, childSpan1, childSpan2, grandchildSpan}
	trace := traces.NewTrace("trace", spans)

	require.Len(t, trace.Flat, 4)
	assert.Equal(t, "Root", trace.RootSpan.Name)

	assert.Equal(t, "child 1", child(t, &trace.RootSpan, 0).Name)
	assert.Equal(t, "child 2", child(t, &trace.RootSpan, 1).Name)
	assert.Equal(t, "grandchild", grandchild(t, &trace.RootSpan, 1, 0).Name)
}

func TestTraceWithMultipleRoots(t *testing.T) {
	root1 := newSpan("Root 1")
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan("Root 2")
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3")
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []traces.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := traces.NewTrace("trace", spans)

	// agreggate root + 3 roots + 3 child
	assert.Len(t, trace.Flat, 7)
	assert.Equal(t, traces.TemporaryRootSpanName, trace.RootSpan.Name)
	assert.Equal(t, "Root 1", child(t, &trace.RootSpan, 0).Name)
	assert.Equal(t, "Root 2", child(t, &trace.RootSpan, 1).Name)
	assert.Equal(t, "Root 3", child(t, &trace.RootSpan, 2).Name)
	assert.Equal(t, "Child from root 1", grandchild(t, &trace.RootSpan, 0, 0).Name)
	assert.Equal(t, "Child from root 2", grandchild(t, &trace.RootSpan, 1, 0).Name)
	assert.Equal(t, "Child from root 3", grandchild(t, &trace.RootSpan, 2, 0).Name)
}

func TestTraceWithMultipleTemporaryRoots(t *testing.T) {
	root1 := newSpan("Temporary Tracetest root span")
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan("Temporary Tracetest root span")
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Temporary Tracetest root span")
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []traces.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := traces.NewTrace("trace", spans)

	require.Len(t, trace.Flat, 4)
	assert.Equal(t, traces.TemporaryRootSpanName, trace.RootSpan.Name)
	assert.Equal(t, "Child from root 1", child(t, &trace.RootSpan, 0).Name)
	assert.Equal(t, "Child from root 2", child(t, &trace.RootSpan, 1).Name)
	assert.Equal(t, "Child from root 3", child(t, &trace.RootSpan, 2).Name)
}

func TestTraceAssemble(t *testing.T) {
	rootSpan := newSpan("Root")
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []traces.Span{rootSpan, childSpan1, grandchildSpan}
	trace := traces.NewTrace("trace", spans)

	assert.Len(t, trace.Flat, 4)
	assert.Equal(t, "Temporary Tracetest root span", trace.RootSpan.Name)
	assert.Equal(t, "Root", child(t, &trace.RootSpan, 0).Name)
	assert.Equal(t, "child 1", grandchild(t, &trace.RootSpan, 0, 0).Name)
	assert.Equal(t, "grandchild", child(t, &trace.RootSpan, 1).Name)

	trace = traces.NewTrace(trace.ID.String(), append(trace.Spans(), childSpan2))
	assert.Len(t, trace.Flat, 4)
	assert.Equal(t, "Root", trace.RootSpan.Name)
	assert.Equal(t, "child 1", child(t, &trace.RootSpan, 0).Name)
	assert.Equal(t, "child 2", child(t, &trace.RootSpan, 1).Name)
	assert.Equal(t, "grandchild", grandchild(t, &trace.RootSpan, 1, 0).Name)
}

func TestTraceWithMultipleRootsFromOtel(t *testing.T) {
	root1 := newOtelSpan("Root 1", nil)
	root1Child := newOtelSpan("Child from root 1", root1)
	root2 := newOtelSpan("Root 2", nil)
	root2Child := newOtelSpan("Child from root 2", root2)
	root3 := newOtelSpan("Root 3", nil)
	root3Child := newOtelSpan("Child from root 3", root3)

	tracesData := &v1.TracesData{
		ResourceSpans: []*v1.ResourceSpans{
			{
				ScopeSpans: []*v1.ScopeSpans{
					{
						Spans: []*v1.Span{root1, root1Child, root2, root2Child, root3, root3Child},
					},
				},
			},
		},
	}

	trace := traces.FromOtel(tracesData)

	// agreggate root + 3 roots + 3 child
	assert.Len(t, trace.Flat, 7)
	assert.Equal(t, traces.TemporaryRootSpanName, trace.RootSpan.Name)
	assert.Equal(t, "Root 1", child(t, &trace.RootSpan, 0).Name)
	assert.Equal(t, "Root 2", child(t, &trace.RootSpan, 1).Name)
	assert.Equal(t, "Root 3", child(t, &trace.RootSpan, 2).Name)
	assert.Equal(t, "Child from root 1", grandchild(t, &trace.RootSpan, 0, 0).Name)
	assert.Equal(t, "Child from root 2", grandchild(t, &trace.RootSpan, 1, 0).Name)
	assert.Equal(t, "Child from root 3", grandchild(t, &trace.RootSpan, 2, 0).Name)
}

func TestInjectingNewRootWhenSingleRoot(t *testing.T) {
	rootSpan := newSpan("Root")
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []traces.Span{rootSpan, childSpan1, childSpan2, grandchildSpan}
	trace := traces.NewTrace("trace", spans)

	newRoot := newSpan("new Root")
	newTrace := trace.InsertRootSpan(newRoot)

	assert.Len(t, newTrace.Flat, 5)
	assert.Equal(t, "new Root", newTrace.RootSpan.Name)
	assert.Len(t, newTrace.RootSpan.Children, 1)
	assert.Equal(t, "Root", child(t, &newTrace.RootSpan, 0).Name)
}

func TestInjectingNewRootWhenMultipleRoots(t *testing.T) {
	root1 := newSpan("Root 1")
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan("Root 2")
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3")
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []traces.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := traces.NewTrace("trace", spans)

	for _, oldRoot := range trace.RootSpan.Children {
		require.NotNil(t, oldRoot.Parent)
	}

	newRoot := newSpan("new Root")
	newTrace := trace.InsertRootSpan(newRoot)

	assert.Len(t, newTrace.Flat, 7)
	assert.Equal(t, "new Root", newTrace.RootSpan.Name)
	assert.Len(t, newTrace.RootSpan.Children, 3)
	assert.Equal(t, "Root 1", child(t, &newTrace.RootSpan, 0).Name)
	assert.Equal(t, "Root 2", child(t, &newTrace.RootSpan, 1).Name)
	assert.Equal(t, "Root 3", child(t, &newTrace.RootSpan, 2).Name)

	for _, oldRoot := range trace.RootSpan.Children {
		require.NotNil(t, oldRoot.Parent)
		assert.Equal(t, newRoot.ID.String(), oldRoot.Parent.ID.String())
	}
}

func TestNoTemporaryRootIfTracetestRootExists(t *testing.T) {
	root1 := newSpan("Root 1")
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan(traces.TriggerSpanName)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3")
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []traces.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := traces.NewTrace("trace", spans)

	assert.Equal(t, root2.ID, trace.RootSpan.ID)
	assert.Equal(t, root2.Name, trace.RootSpan.Name)
}

func TestNewTemporaryRootIfATemporaryRootExists(t *testing.T) {
	root1 := newSpan("Root 1")
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan(traces.TemporaryRootSpanName)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3")
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []traces.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := traces.NewTrace("trace", spans)

	assert.NotEqual(t, root2.ID, trace.RootSpan.ID)
	assert.Equal(t, traces.TemporaryRootSpanName, root2.Name)
}

func TestTriggerSpanShouldBeRootWhenTemporaryRootExistsToo(t *testing.T) {
	root1 := newSpan(traces.TriggerSpanName)
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan(traces.TemporaryRootSpanName)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3")
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []traces.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := traces.NewTrace("trace", spans)

	assert.Equal(t, root1.ID, trace.RootSpan.ID)
	assert.Equal(t, root1.Name, trace.RootSpan.Name)
}

func TestEventsAreInjectedIntoAttributes(t *testing.T) {
	rootSpan := newSpan("Root", withEvents([]traces.SpanEvent{
		{Name: "event 1", Attributes: attributesFromMap(map[string]string{"attribute1": "value"})},
		{Name: "event 2", Attributes: attributesFromMap(map[string]string{"attribute2": "value"})},
	}))
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []traces.Span{rootSpan, childSpan1, childSpan2, grandchildSpan}
	trace := traces.NewTrace("trace", spans)

	require.NotEmpty(t, trace.RootSpan.Attributes.Get("span.events"))

	events := []traces.SpanEvent{}
	err := json.Unmarshal([]byte(trace.RootSpan.Attributes.Get("span.events")), &events)
	require.NoError(t, err)

	assert.Equal(t, "event 1", events[0].Name)
	assert.Equal(t, "value", events[0].Attributes.Get("attribute1"))
	assert.Equal(t, "event 2", events[1].Name)
	assert.Equal(t, "value", events[1].Attributes.Get("attribute2"))
}

func TestMergingZeroTraces(t *testing.T) {
	trace := traces.MergeTraces()
	assert.Nil(t, trace)
}

func TestMergingFragmentsFromSameTrace(t *testing.T) {
	traceID := id.NewRandGenerator().TraceID().String()
	rootSpan := newSpan("Root")
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))

	firstFragment := traces.NewTrace(traceID, []traces.Span{childSpan2})
	secondFragment := traces.NewTrace(traceID, []traces.Span{rootSpan, childSpan1})

	trace := traces.MergeTraces(&firstFragment, &secondFragment)
	require.NotNil(t, trace)
	assert.NotEmpty(t, trace.ID)
	assert.Len(t, trace.Flat, 3)
}

type option func(*traces.Span)

func withParent(parent *traces.Span) option {
	return func(s *traces.Span) {
		s.Parent = parent
	}
}

func withEvents(events []traces.SpanEvent) option {
	return func(s *traces.Span) {
		s.Events = events
	}
}

func newSpan(name string, options ...option) traces.Span {
	span := traces.Span{
		ID:         id.NewRandGenerator().SpanID(),
		Name:       name,
		Attributes: traces.NewAttributes(),
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(1 * time.Second),
	}

	for _, option := range options {
		option(&span)
	}

	if span.Parent != nil {
		span.Attributes.Set(traces.TracetestMetadataFieldParentID, span.Parent.ID.String())
	}

	return span
}

func newOtelSpan(name string, parent *v1.Span) *v1.Span {
	id := id.NewRandGenerator().SpanID()
	var parentId []byte = nil
	if parent != nil {
		parentId = parent.SpanId
	}

	return &v1.Span{
		SpanId:            id[:],
		Name:              name,
		ParentSpanId:      parentId,
		StartTimeUnixNano: uint64(time.Now().UnixNano()),
		EndTimeUnixNano:   uint64(time.Now().Add(1 * time.Second).UnixNano()),
	}
}

func TestJSONEncoding(t *testing.T) {

	rootSpan := createSpan("root")
	subSpan1 := createSpan("subSpan1")
	subSubSpan1 := createSpan("subSubSpan1")
	subSpan2 := createSpan("subSpan2")

	subSpan1.Parent = rootSpan
	subSpan2.Parent = rootSpan
	subSubSpan1.Parent = subSpan1

	flat := map[trace.SpanID]*traces.Span{
		// We copy those spans so they don't have children injected into them
		// the flat structure shouldn't have children.
		rootSpan.ID:    copySpan(rootSpan),
		subSpan1.ID:    copySpan(subSpan1),
		subSubSpan1.ID: copySpan(subSubSpan1),
		subSpan2.ID:    copySpan(subSpan2),
	}

	rootSpan.Children = []*traces.Span{subSpan1, subSpan2}
	subSpan1.Children = []*traces.Span{subSubSpan1}

	tid := id.NewRandGenerator().TraceID()
	trace := traces.Trace{
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
		var actual traces.Trace
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

func copySpan(span *traces.Span) *traces.Span {
	newSpan := *span
	return &newSpan
}

func createSpan(name string) *traces.Span {
	return &traces.Span{
		ID:        id.NewRandGenerator().SpanID(),
		Name:      name,
		StartTime: time.Date(2021, 11, 24, 14, 05, 12, 0, time.UTC),
		EndTime:   time.Date(2021, 11, 24, 14, 05, 17, 0, time.UTC),
		Attributes: attributesFromMap(map[string]string{
			"service.name": name,
		}),
		Children: nil,
	}
}

var now = time.Now()

func getTime(n int) time.Time {
	return now.Add(time.Duration(n) * time.Second)
}

func TestSort(t *testing.T) {
	randGenerator := id.NewRandGenerator()
	trace := traces.Trace{
		ID: randGenerator.TraceID(),
		RootSpan: traces.Span{
			Name:       "root",
			StartTime:  getTime(0),
			Attributes: traces.Attributes{},
			Children: []*traces.Span{
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
					Children: []*traces.Span{
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

	expectedTrace := traces.Trace{
		ID: randGenerator.TraceID(),
		RootSpan: traces.Span{
			Name:       "root",
			StartTime:  getTime(0),
			Attributes: traces.Attributes{},
			Children: []*traces.Span{
				{
					Name:      "child 1",
					StartTime: getTime(1),
					Children: []*traces.Span{
						{

							Name:      "grandchild 1",
							StartTime: getTime(2),
							Children:  make([]*traces.Span, 0),
						},
						{

							Name:      "grandchild 2",
							StartTime: getTime(3),
							Children:  make([]*traces.Span, 0),
						},
					},
				},
				{
					Name:      "child 2",
					StartTime: getTime(2),
					Children:  make([]*traces.Span, 0),
				},
				{
					Name:      "child 3",
					StartTime: getTime(3),
					Children:  make([]*traces.Span, 0),
				},
			},
		},
	}

	assert.Equal(t, expectedTrace.RootSpan, sortedTrace.RootSpan)
}

func TestUnmarshalLargeTrace(t *testing.T) {
	bytes, err := os.ReadFile("./data/big-trace-json.json")
	require.NoError(t, err)

	trace := traces.Trace{}

	err = json.Unmarshal(bytes, &trace)
	require.NoError(t, err)

	assert.Greater(t, len(trace.Flat), 0)
}

func TestBrowserSpan(t *testing.T) {
	span := newSpan("click")
	span.Attributes = attributesFromMap(map[string]string{
		"event_type": "click",
		"http.url":   "http://localhost:1663",
	})

	spans := []traces.Span{span}
	trace := traces.NewTrace("trace", spans)

	assert.Equal(t, trace.Spans()[0].Attributes.Get(traces.TracetestMetadataFieldType), "general")
}

func attributesFromMap(input map[string]string) traces.Attributes {
	attributes := traces.NewAttributes()
	for key, value := range input {
		attributes.Set(key, value)
	}

	return attributes
}

func child(t *testing.T, span *traces.Span, index int) *traces.Span {
	if len(span.Children) < index+1 {
		t.FailNow()
	}

	child := span.Children[index]
	if child == nil {
		t.FailNow()
	}

	return child
}

func grandchild(t *testing.T, span *traces.Span, parentIndex int, grandChildIndex int) *traces.Span {
	return child(t, child(t, span, parentIndex), grandChildIndex)
}
