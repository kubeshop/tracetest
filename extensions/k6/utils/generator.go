package utils

import (
	"math/rand"
	"time"

	"go.opentelemetry.io/otel/trace"
)

func TraceID() trace.TraceID {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	tid := trace.TraceID{}

	rand.Read(tid[:])
	return tid
}

func SpanID() trace.SpanID {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	sid := trace.SpanID{}
	rand.Read(sid[:])
	return sid
}
