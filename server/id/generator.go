package id

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Generator interface {
	UUID() uuid.UUID
	TraceID() trace.TraceID
	SpanID() trace.SpanID
}

func NewRandGenerator() Generator {
	return randGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type randGenerator struct {
	rand *rand.Rand
}

func (g randGenerator) UUID() uuuuid.UUID {
	return uuid.New()
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
