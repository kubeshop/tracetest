package id

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"go.opentelemetry.io/otel/trace"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type Generator interface {
	UUID() uuid.UUID
	ID() ID
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

func (g randGenerator) UUID() uuid.UUID {
	return uuid.New()
}

func (g randGenerator) ID() ID {
	return ID(shortid.MustGenerate())
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
