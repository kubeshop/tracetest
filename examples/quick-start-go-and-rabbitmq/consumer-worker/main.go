package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/consumer-worker/config"
	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/consumer-worker/messagequeue"
	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/consumer-worker/telemetry"
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

	logger.Info("Setting up worker...")

	logger.Info("Initializing OpenTelemetry...")
	tracer, err := telemetry.Setup(ctx, currentConfig.OtelExporterEndpoint, currentConfig.OtelServiceName)
	if err != nil {
		logger.Error("Unable to setup OpenTelemetry", zap.Error(err))
		return
	}
	logger.Info("OpenTelemetry initialized.")

	logger.Info("Initializing Message Queue reader...")
	reader, err := messagequeue.GetMessageQueueReader(currentConfig.MessageQueueConnectionString, currentConfig.MessageQueueName)
	if err != nil {
		logger.Error("Unable to setup Message Queue reader", zap.Error(err))
		return
	}
	logger.Info("Message Queue reader initialized.")

	logger.Info("Starting worker...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = reader.Read(ctx, getMessageReader(tracer, logger))
	if err != nil {
		logger.Error("Unable to read messages from Message Queue", zap.Error(err))
		return
	}

	<-ctx.Done()
	logger.Info("Worker stop signal detected")
}

func getMessageReader(tracer trace.Tracer, logger *zap.Logger) func(context.Context, string, string) {
	return func(readerContext context.Context, message string) {
		_, span := tracer.Start(readerContext, "Process incoming message")
		defer span.End()

		logger.Info("Incoming message", zap.String("message", message))
	}
}
