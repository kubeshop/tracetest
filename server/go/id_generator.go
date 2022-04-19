package openapi

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type IDGenerator interface {
	UUID() string
	TraceID() trace.TraceID
	SpanID() trace.SpanID
}

func NewRandGenerator() IDGenerator {
	return randGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type randGenerator struct {
	rand *rand.Rand
}

func (g randGenerator) UUID() string {
	return uuid.New().String()
}

func (g randGenerator) TraceID() trace.TraceID {
	tid := trace.TraceID{}
	g.rand.Read(tid[:])
	return tid
}
func (g randGenerator) SpanID() trace.SpanID {
	sid := trace.SpanID{}
	g.rand.Read(sid[:])
	return sid
}
