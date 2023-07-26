package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/producer-api/config"
	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/producer-api/telemetry"
)

func main() {
	currentConfig := config.Current()

	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("error creating zap logger, error: %v", err)
		return
	}

	spanExporter, err := telemetry.NewSpanExporter(ctx, currentConfig.OtelExporterEndpoint)
	if err != nil {
		logger.Fatal("failed to initialize exporter", zap.Error(err))
	}

	traceProvider := telemetry.NewTraceProvider(spanExporter, currentConfig.OtelServiceName, logger)
	defer func() { _ = traceProvider.Shutdown(ctx) }()
	otel.SetTracerProvider(traceProvider)

	tracer := traceProvider.Tracer(currentConfig.OtelServiceName)

	mux := http.NewServeMux()
	mux.HandleFunc("/publish", publishHandler(tracer))

	fmt.Println("Starting producer-api on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("error running server", zap.Error(err))
	}
}

func publishHandler(tracer trace.Tracer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, span := tracer.Start(req.Context(), "GET /publish")
		defer span.End()

		bodyContent, err := io.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(string(bodyContent))
	}
}
