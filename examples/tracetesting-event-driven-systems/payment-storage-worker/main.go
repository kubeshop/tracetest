package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/kubeshop/tracetest/examples/tracetesting-event-driven-systems/payment-storage-worker/config"
	"github.com/kubeshop/tracetest/examples/tracetesting-event-driven-systems/payment-storage-worker/streaming"
	"github.com/kubeshop/tracetest/examples/tracetesting-event-driven-systems/payment-storage-worker/telemetry"
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
	logger.Info("Kafka reader initialized.")

	logger.Info("Starting worker...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	processor := &MessageProcessor{
		tracer:               tracer,
		logger:               logger,
		paymentOrderDatabase: []*PaymentOrder{},
	}

	err = reader.Read(ctx, processor.ReadMessage)
	if err != nil {
		logger.Error("Unable to read messages from Kafka", zap.Error(err))
		return
	}

	<-ctx.Done()
	logger.Info("Worker stop signal detected")
}

type MessageProcessor struct {
	tracer               trace.Tracer
	logger               *zap.Logger
	paymentOrderDatabase []*PaymentOrder
}

type PaymentOrder struct {
	OriginCustomerID      string  `json:"originCustomerID"`
	DestinationCustomerID string  `json:"destinationCustomerID"`
	Value                 float64 `json:"value"`
}

func (p *MessageProcessor) ReadMessage(ctx context.Context, topic, message string) {
	ctx, span := p.tracer.Start(ctx, "Process incoming paymentOrder")
	defer span.End()

	p.logger.Info("Incoming message", zap.String("topic", topic), zap.String("message", message))

	var order PaymentOrder
	err := json.Unmarshal([]byte(message), &order)
	if err != nil {
		p.logger.Error("paymentOrder with invalid format", zap.Error(err))
	}

	p.StorePaymentOrder(ctx, &order)
}

func (p *MessageProcessor) StorePaymentOrder(ctx context.Context, order *PaymentOrder) {
	ctx, span := p.tracer.Start(ctx, "Storing paymentOrder")
	defer span.End()

	p.paymentOrderDatabase = append(p.paymentOrderDatabase, order)
}
