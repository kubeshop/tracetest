package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/producer-api/config"
	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/producer-api/messagequeue"
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

	logger.Info("Initializing Message Queue publisher...")
	publisher, err := messagequeue.GetMessageQueuePublisher(currentConfig.MessageQueueConnectionString, currentConfig.MessageQueueName)
	if err != nil {
		logger.Error("Unable to setup Message Queue publisher", zap.Error(err))
		return
	}
	defer publisher.Close()
	logger.Info("Message Queue publisher initialized.")

	handler := http.HandlerFunc(publishHandler(tracer, publisher))
	otelHandler := otelhttp.NewHandler(handler, "GET /publish")

	http.Handle("/publish", otelHandler)

	logger.Info("Starting producer-api on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("unable to start server", zap.Error(err))
		return
	}
}

func publishHandler(tracer trace.Tracer, publisher *messagequeue.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx, span := tracer.Start(req.Context(), "publish handler")
		defer span.End()

		bodyContent, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("error when reading body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = publisher.Publish(ctx, bodyContent)
		if err != nil {
			fmt.Printf("error when publishing message: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
