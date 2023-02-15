package id

import (
	"encoding/binary"
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
	var r [16]byte
	epoch := time.Now().Unix()
	binary.BigEndian.PutUint32(r[0:4], uint32(epoch))
	_, err := rand.Read(r[4:])
	if err != nil {
		panic(err)
	}

	return trace.TraceID(r)
}
func (g randGenerator) SpanID() trace.SpanID {
	sid := trace.SpanID{}
	g.rand.Read(sid[:])
	return sid
}
