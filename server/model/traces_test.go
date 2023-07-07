package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestTraces(t *testing.T) {
	rootSpan := newSpan("Root", nil)
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []model.Span{rootSpan, childSpan1, childSpan2, grandchildSpan}
	trace := model.NewTrace("trace", spans)

	assert.Len(t, trace.Flat, 4)
	assert.Equal(t, "Root", trace.RootSpan.Name)
	assert.Equal(t, "child 1", trace.RootSpan.Children[0].Name)
	assert.Equal(t, "child 2", trace.RootSpan.Children[1].Name)
	assert.Equal(t, "grandchild", trace.RootSpan.Children[1].Children[0].Name)
}

func TestTraceWithMultipleRoots(t *testing.T) {
	root1 := newSpan("Root 1", nil)
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan("Root 2", nil)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3", nil)
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []model.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := model.NewTrace("trace", spans)

	// agreggate root + 3 roots + 3 child
	assert.Len(t, trace.Flat, 7)
	assert.Equal(t, model.TemporaryRootSpanName, trace.RootSpan.Name)
	assert.Equal(t, "Root 1", trace.RootSpan.Children[0].Name)
	assert.Equal(t, "Root 2", trace.RootSpan.Children[1].Name)
	assert.Equal(t, "Root 3", trace.RootSpan.Children[2].Name)
	assert.Equal(t, "Child from root 1", trace.RootSpan.Children[0].Children[0].Name)
	assert.Equal(t, "Child from root 2", trace.RootSpan.Children[1].Children[0].Name)
	assert.Equal(t, "Child from root 3", trace.RootSpan.Children[2].Children[0].Name)
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
	assert.Equal(t, model.TemporaryRootSpanName, trace.RootSpan.Name)
	assert.Equal(t, "Root 1", trace.RootSpan.Children[0].Name)
	assert.Equal(t, "Root 2", trace.RootSpan.Children[1].Name)
	assert.Equal(t, "Root 3", trace.RootSpan.Children[2].Name)
	assert.Equal(t, "Child from root 1", trace.RootSpan.Children[0].Children[0].Name)
	assert.Equal(t, "Child from root 2", trace.RootSpan.Children[1].Children[0].Name)
	assert.Equal(t, "Child from root 3", trace.RootSpan.Children[2].Children[0].Name)
}

func TestInjectingNewRootWhenSingleRoot(t *testing.T) {
	rootSpan := newSpan("Root", nil)
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []model.Span{rootSpan, childSpan1, childSpan2, grandchildSpan}
	trace := model.NewTrace("trace", spans)

	newRoot := newSpan("new Root", nil)
	newTrace := trace.InsertRootSpan(newRoot)

	assert.Len(t, newTrace.Flat, 5)
	assert.Equal(t, "new Root", newTrace.RootSpan.Name)
	assert.Len(t, newTrace.RootSpan.Children, 1)
	assert.Equal(t, "Root", newTrace.RootSpan.Children[0].Name)
}

func TestInjectingNewRootWhenMultipleRoots(t *testing.T) {
	root1 := newSpan("Root 1", nil)
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan("Root 2", nil)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3", nil)
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []model.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := model.NewTrace("trace", spans)

	for _, oldRoot := range trace.RootSpan.Children {
		require.NotNil(t, oldRoot.Parent)
	}

	newRoot := newSpan("new Root", nil)
	newTrace := trace.InsertRootSpan(newRoot)

	assert.Len(t, newTrace.Flat, 7)
	assert.Equal(t, "new Root", newTrace.RootSpan.Name)
	assert.Len(t, newTrace.RootSpan.Children, 3)
	assert.Equal(t, "Root 1", newTrace.RootSpan.Children[0].Name)
	assert.Equal(t, "Root 2", newTrace.RootSpan.Children[1].Name)
	assert.Equal(t, "Root 3", newTrace.RootSpan.Children[2].Name)

	for _, oldRoot := range trace.RootSpan.Children {
		require.NotNil(t, oldRoot.Parent)
		assert.Equal(t, newRoot.ID.String(), oldRoot.Parent.ID.String())
	}
}

func TestNoTemporaryRootIfTracetestRootExists(t *testing.T) {
	root1 := newSpan("Root 1", nil)
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan(model.TriggerSpanName, nil)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3", nil)
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []model.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := model.NewTrace("trace", spans)

	assert.Equal(t, root2.ID, trace.RootSpan.ID)
	assert.Equal(t, root2.Name, trace.RootSpan.Name)
}

func TestNoTemporaryRootIfATemporaryRootExists(t *testing.T) {
	root1 := newSpan("Root 1", nil)
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan(model.TemporaryRootSpanName, nil)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3", nil)
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []model.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := model.NewTrace("trace", spans)

	assert.Equal(t, root2.ID, trace.RootSpan.ID)
	assert.Equal(t, root2.Name, trace.RootSpan.Name)
}

func TestTriggerSpanShouldBeRootWhenTemporaryRootExistsToo(t *testing.T) {
	root1 := newSpan(model.TriggerSpanName, nil)
	root1Child := newSpan("Child from root 1", withParent(&root1))
	root2 := newSpan(model.TemporaryRootSpanName, nil)
	root2Child := newSpan("Child from root 2", withParent(&root2))
	root3 := newSpan("Root 3", nil)
	root3Child := newSpan("Child from root 3", withParent(&root3))

	spans := []model.Span{root1, root1Child, root2, root2Child, root3, root3Child}
	trace := model.NewTrace("trace", spans)

	assert.Equal(t, root1.ID, trace.RootSpan.ID)
	assert.Equal(t, root1.Name, trace.RootSpan.Name)
}

func TestEventsAreInjectedIntoAttributes(t *testing.T) {
	rootSpan := newSpan("Root", withEvents([]model.SpanEvent{
		{Name: "event 1", Attributes: model.Attributes{"attribute1": "value"}},
		{Name: "event 2", Attributes: model.Attributes{"attribute2": "value"}},
	}))
	childSpan1 := newSpan("child 1", withParent(&rootSpan))
	childSpan2 := newSpan("child 2", withParent(&rootSpan))
	grandchildSpan := newSpan("grandchild", withParent(&childSpan2))

	spans := []model.Span{rootSpan, childSpan1, childSpan2, grandchildSpan}
	trace := model.NewTrace("trace", spans)

	require.NotEmpty(t, trace.RootSpan.Attributes["events"])

	events := []model.SpanEvent{}
	err := json.Unmarshal([]byte(trace.RootSpan.Attributes["events"]), &events)
	require.NoError(t, err)

	assert.Equal(t, "event 1", events[0].Name)
	assert.Equal(t, "value", events[0].Attributes["attribute1"])
	assert.Equal(t, "event 2", events[1].Name)
	assert.Equal(t, "value", events[1].Attributes["attribute2"])
}

type option func(*model.Span)

func withParent(parent *model.Span) option {
	return func(s *model.Span) {
		s.Parent = parent
	}
}

func withEvents(events []model.SpanEvent) option {
	return func(s *model.Span) {
		s.Events = events
	}
}

func newSpan(name string, options ...option) model.Span {
	span := model.Span{
		ID:         id.NewRandGenerator().SpanID(),
		Name:       name,
		Attributes: make(model.Attributes),
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(1 * time.Second),
	}

	for _, option := range options {
		option(&span)
	}

	if span.Parent != nil {
		span.Attributes[model.TracetestMetadataFieldParentID] = span.Parent.ID.String()
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
