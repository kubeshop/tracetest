package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// add instrumentation lib
	"github.com/kubeshop/tracetest/examples/observability-driven-development-go-tracetest/bookstore/part3.2/lib/instrumentation"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const svcName = "books"

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
		{"4", "The art of war", 0},
	}, nil
}
