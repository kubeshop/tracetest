package envelope_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers/envelope"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestSpanEnvelope(t *testing.T) {
	// each of these spans take 73 bytes
	spans := []*proto.Span{
		createSpan(),
		createSpan(),
	}

	envelope := envelope.EnvelopeSpans(spans, 150)
	// 2 * 73 < 150, so all spans should be included in the envelope
	require.Len(t, envelope, 2)
}

func TestSpanEnvelopeShouldIgnoreRestOfSpansIfTheyDontFit(t *testing.T) {
	// each of these spans take 73 bytes
	span1, span2, span3 := createSpan(), createSpan(), createSpan()
	spans := []*proto.Span{
		span1,
		span2,
		span3,
	}

	envelope := envelope.EnvelopeSpans(spans, 150)
	// 3 * 73 > 150, but 2 spans fit the envelope, so take 2 spans instead
	require.Len(t, envelope, 2)
	assert.Equal(t, envelope[0].Id, span1.Id)
	assert.Equal(t, envelope[1].Id, span2.Id)
}

func TestSpanEnvelopeShouldFitAsManySpansAsPossible(t *testing.T) {
	// these spans take 73 bytes, 73 bytes, and 33 bytes respectively
	span1, span2, span3 := createSpan(), createSpan(), createSmallSpan()
	spans := []*proto.Span{
		span1,
		span2,
		span3,
	}

	envelope := envelope.EnvelopeSpans(spans, 110)
	// 73+73+33 > 110, 73+73 is also bigger than 110. But we can add 73 (span1) + 33 (span3) to fit the envelope
	require.Len(t, envelope, 2)
	assert.Equal(t, envelope[0].Id, span1.Id)
	assert.Equal(t, envelope[1].Id, span3.Id)
}

func TestSpanEnvelopeShouldAllowLargeSpans(t *testing.T) {
	// a large span is 682 bytes long, in theory, it should not fit the envelope, however,
	// we should allow 1 per envelope just to make sure ALL spans are sent.
	spans := []*proto.Span{
		createLargeSpan(),
	}

	envelope := envelope.EnvelopeSpans(spans, 110)
	require.Len(t, envelope, 1)
}

func TestSpanEnvelopeShouldOnlyAllowOneLargeSpan(t *testing.T) {
	// a large span is 682 bytes long, in theory, it should not fit the envelope, however,
	// we should allow 1 per envelope just to make sure ALL spans are sent.
	largeSpan1, largeSpan2 := createLargeSpan(), createLargeSpan()
	spans := []*proto.Span{
		largeSpan1,
		largeSpan2,
	}

	envelope := envelope.EnvelopeSpans(spans, 110)
	require.Len(t, envelope, 1)
	assert.Equal(t, largeSpan1.Id, envelope[0].Id)
}

func TestSpanEnvelopeShouldAddALargeSpanEvenIfThereAreMoreSpansInIt(t *testing.T) {
	// Given a list of small spans that should be able to fit the envelope, but also a large span that doesn't fit an envelope,
	// it should include the largeSpan with the 2 first spans and leave the third small span out
	smallSpan1, smallSpan2, largeSpan1, smallSpan3 := createSmallSpan(), createSmallSpan(), createLargeSpan(), createSmallSpan()
	spans := []*proto.Span{
		smallSpan1,
		smallSpan2,
		largeSpan1,
		smallSpan3,
	}

	envelope := envelope.EnvelopeSpans(spans, 100)
	require.Len(t, envelope, 3)
	assert.Equal(t, smallSpan1.Id, envelope[0].Id)
	assert.Equal(t, smallSpan2.Id, envelope[1].Id)
	assert.Equal(t, largeSpan1.Id, envelope[2].Id)
}

func TestSpanEnvelopeShouldIgnoreExtraLargeSpan(t *testing.T) {
	// Given a list of small spans that should be able to fit the envelope, but also a large span that doesn't fit an envelope,
	// it should include the largeSpan with the 2 first spans and leave the third small span out
	smallSpan1, smallSpan2, largeSpan1, smallSpan3, largeSpan2 := createSmallSpan(), createSmallSpan(), createLargeSpan(), createSmallSpan(), createLargeSpan()
	spans := []*proto.Span{
		smallSpan1,
		smallSpan2,
		largeSpan1,
		smallSpan3,
		largeSpan2,
	}

	envelope := envelope.EnvelopeSpans(spans, 100)
	require.Len(t, envelope, 3)
	assert.Equal(t, smallSpan1.Id, envelope[0].Id)
	assert.Equal(t, smallSpan2.Id, envelope[1].Id)
	assert.Equal(t, largeSpan1.Id, envelope[2].Id)
}

func createSpan() *proto.Span {
	return &proto.Span{
		Id:        id.NewRandGenerator().SpanID().String(),
		Name:      "span name",
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Add(2 * time.Second).Unix(),
		Kind:      "internal",
		Attributes: []*proto.KeyValuePair{
			{
				Key:   "service.name",
				Value: "core",
			},
		},
	}
}

func createSmallSpan() *proto.Span {
	return &proto.Span{
		Id:         id.NewRandGenerator().SpanID().String(),
		Name:       "s",
		StartTime:  time.Now().Unix(),
		EndTime:    time.Now().Add(2 * time.Second).Unix(),
		Kind:       "",
		Attributes: []*proto.KeyValuePair{},
	}
}

func createLargeSpan() *proto.Span {
	loremIpsum := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum eu fermentum elit. Ut convallis elit nisl, et porttitor ante dignissim quis. Curabitur porttitor molestie iaculis. Suspendisse potenti. Curabitur sollicitudin finibus mollis. Nunc at tincidunt dolor. Nam eleifend ante in elit vulputate lacinia. Donec sem orci, luctus ut eros id, tincidunt elementum nulla. Nulla et nibh pharetra, pretium odio nec, posuere est. Curabitur a felis ut risus fermentum ornare vitae sed dolor. Mauris non velit at nulla ultricies mattis. "
	return &proto.Span{
		Id:        id.NewRandGenerator().SpanID().String(),
		Name:      loremIpsum,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Add(2 * time.Second).Unix(),
		Kind:      "internal",
		Attributes: []*proto.KeyValuePair{
			{Key: "service.name", Value: "core"},
			{Key: "service.team", Value: "ranchers"},
			{Key: "go.version", Value: "1.22.3"},
			{Key: "go.os", Value: "Linux"},
			{Key: "go.arch", Value: "amd64"},
		},
	}
}
