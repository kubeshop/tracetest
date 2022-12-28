package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/adnanrahic/bookstore/part3.2/lib/instrumentation"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const svcName = "availability"

var tracer trace.Tracer

func main() {
	ctx := context.Background()

	exp, err := instrumentation.NewExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := instrumentation.NewTraceProvider(svcName, exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer(svcName)

	r := mux.NewRouter()
	r.Use(otelmux.Middleware(svcName))

	r.HandleFunc("/{bookID}", stockHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

var books = map[string]string{
	"1": "10",
	"2": "1",
	"3": "5",
	"4": "0",
}

func stockHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "Availability Check")
	defer span.End()

	vars := mux.Vars(r)
	bookID, ok := vars["bookID"]
	if !ok {
		span.SetStatus(codes.Error, "no bookID in URL")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "missing bookID in URL")
		return
	}

	span.SetAttributes(
		attribute.String("bookID", bookID),
	)

	stock, ok := books[bookID]
	if !ok {
		span.SetStatus(codes.Error, "book not found")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "book not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, stock)
}
