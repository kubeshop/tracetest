package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/producer-api/config"
	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/producer-api/streaming"
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
	defer logger.Sync()

	logger.Info("Setting up server...")

	logger.Info("Initializing OpenTelemetry...")
	tracer, err := telemetry.Setup(ctx, currentConfig.OtelExporterEndpoint, currentConfig.OtelServiceName)
	if err != nil {
		logger.Error("Unable to setup OpenTelemetry", zap.Error(err))
		return
	}
	logger.Info("OpenTelemetry initialized.")

	logger.Info("Initializing Kafka publisher...")
	publisher, err := streaming.GetKafkaPublisher(currentConfig.KafkaBrokerUrl, currentConfig.KafkaTopic)
	if err != nil {
		logger.Error("Unable to setup Kafka publisher", zap.Error(err))
		return
	}
	defer publisher.Close()
	logger.Info("Kafka publisher initialized.")

	mux := http.NewServeMux()
	mux.HandleFunc("/publish", publishHandler(tracer, publisher))

	logger.Info("Starting producer-api on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Error("unable to start server", zap.Error(err))
		return
	}
}

func publishHandler(tracer trace.Tracer, publisher *streaming.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, span := tracer.Start(req.Context(), "GET /publish")
		defer span.End()

		bodyContent, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("error when reading body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = publisher.Publish(req.Context(), bodyContent)
		if err != nil {
			fmt.Printf("error when publishing message: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
