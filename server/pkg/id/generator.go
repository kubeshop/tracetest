package id

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"go.opentelemetry.io/otel/trace"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type GeneratorFunc func() ID

func GenerateID() ID {
	return ID(shortid.MustGenerate())
}

func SlugFromString(input string) ID {
	id := strings.ReplaceAll(input, " ", "-")
	id = strings.ToLower(id)

	return ID(id)
}

type Generator interface {
	UUID() uuid.UUID
	ID() ID
	TraceID() trace.TraceID
	SpanID() trace.SpanID
}

func NewRandGenerator() Generator {
	return randGenerator{}
}

type randGenerator struct{}

func (g randGenerator) UUID() uuid.UUID {
	return uuid.New()
}

func (g randGenerator) ID() ID {
	return GenerateID()
}

func (g randGenerator) TraceID() trace.TraceID {
	tid := trace.TraceID{}
	rndSeed := rand.NewSource(time.Now().UnixNano())
	rand.New(rndSeed).Read(tid[:])
	return tid
}

func (g randGenerator) SpanID() trace.SpanID {
	sid := trace.SpanID{}
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(sid[:])
	return sid
}
