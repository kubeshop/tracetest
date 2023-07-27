package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/consumer-worker/config"
	"github.com/kubeshop/tracetest/examples/quick-start-go-and-kafka/consumer-worker/streaming"
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

	logger.Info("Initializing Kafka reader...")
	reader, err := streaming.GetKafkaReader(currentConfig.KafkaBrokerUrl, currentConfig.KafkaTopic)
	if err != nil {
		logger.Error("Unable to setup Kafka reader", zap.Error(err))
		return
	}
	defer reader.Close()
	logger.Info("Kafka reader initialized.")

	// based on https://github.com/zaynkorai/go-kafka-example/blob/main/worker/worker.go
	logger.Info("Starting worker...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-reader.PartitionConsumer().Errors():
				logger.Error("Error on reader", zap.Error(err))
			case msg := <-reader.PartitionConsumer().Messages():
				_, span := tracer.Start(ctx, "Incoming message")
				logger.Info("Incoming message", zap.String("topic", string(msg.Topic)), zap.String("message", string(msg.Value)))
				span.End()
			case <-sigchan:
				logger.Info("Worker stop signal detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}
