package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const svcName = "books"

var tracer trace.Tracer

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "otel-collector:4317", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return traceExporter, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(svcName),
		),
	)

	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}

func main() {
	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer(svcName)

	r := mux.NewRouter()
	r.Use(otelmux.Middleware(svcName))

	r.HandleFunc("/books", booksListHandler)

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8001", nil))
}

func httpError(span trace.Span, w http.ResponseWriter, msg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, msg)
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
}

func booksListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "Books List")
	defer span.End()

	books, err := getAvailableBooks(ctx)
	if err != nil {
		httpError(span, w, "cannot read books DB", err)
		return
	}

	span.SetAttributes(
		attribute.Int("books.list.count", len(books)),
	)

	jsonBooks, err := json.Marshal(books)
	if err != nil {
		httpError(span, w, "cannot json encode books", err)
		return
	}

	w.Write(jsonBooks)
}

func getAvailableBooks(ctx context.Context) ([]book, error) {
	books, err := getBooks(ctx)
	if err != nil {
		return nil, err
	}

	availableBook := make([]book, 0, len(books))
	for _, book := range books {
		available, err := isBookAvailable(ctx, book.ID)
		if err != nil {
			return nil, err
		}

		if !available {
			continue
		}
		availableBook = append(availableBook, book)
	}

	return availableBook, nil
}

var httpClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}

func isBookAvailable(ctx context.Context, bookID string) (bool, error) {
	ctx, span := tracer.Start(ctx, "Availability Request", trace.WithAttributes(
		attribute.String("bookID", bookID),
	))
	defer span.End()

	url := "http://availability:8000/" + bookID
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "cannot do request")
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		span.SetStatus(codes.Error, "not found")
		return false, nil
	}

	stockBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "cannot read response body")
		return false, err
	}

	stock, err := strconv.Atoi(string(stockBytes))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "cannot parse stock value")
		return false, err
	}

	return stock > 0, nil
}

type book struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func getBooks(ctx context.Context) ([]book, error) {
	return []book{
		{"1", "Harry Potter", 0},
		{"2", "Foundation", 0},
		{"3", "Moby Dick", 0},
		{"4", "The art of war", 0}, // Add this book
	}, nil
}
