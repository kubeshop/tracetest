package envelope

import (
	agentProto "github.com/kubeshop/tracetest/agent/proto"
	"google.golang.org/protobuf/proto"
)

// EnvelopeSpans get a list of spans and batch them into a packet that does not
// surpasses the maxPacketSize restriction. When maxPacketSize is reached, no
// more spans are added to the packet,
func EnvelopeSpans(spans []*agentProto.Span, maxPacketSize int) []*agentProto.Span {
	envelope := make([]*agentProto.Span, 0, len(spans))
	currentSize := 0

	// There's a weird scenario that must be covered here: imagine a span so big it is bigger than maxPacketSize.
	// It is impossible to send a span like this, so in this case, we classify those spans as "large spans" and we allow
	// `largeSpansPerPacket` per packet.
	//
	// It is important to ensure a limit of large spans per packet because if your whole trace is composed by
	// large spans, this would mean a packet would hold the entiry trace and we don't want that to happen.
	const largeSpansPerPacket = 1
	numberLargeSpansAddedToPacket := 0

	for _, span := range spans {
		spanSize := proto.Size(span)
		isLargeSpan := spanSize > maxPacketSize
		if currentSize+spanSize < maxPacketSize || isLargeSpan {
			if isLargeSpan {
				if numberLargeSpansAddedToPacket >= largeSpansPerPacket {
					// there is already the limit of large spans in the packet, skip this one
					continue
				}

				numberLargeSpansAddedToPacket++
			}
			envelope = append(envelope, span)
			currentSize += spanSize
		}
	}

	return envelope
}
