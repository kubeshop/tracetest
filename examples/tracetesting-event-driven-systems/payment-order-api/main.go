package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kubeshop/tracetest/examples/tracetesting-event-driven-systems/payment-order-api/config"
	"github.com/kubeshop/tracetest/examples/tracetesting-event-driven-systems/payment-order-api/streaming"
	"github.com/kubeshop/tracetest/examples/tracetesting-event-driven-systems/payment-order-api/telemetry"
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

	handler := http.HandlerFunc(paymentHandler(tracer, publisher))
	otelHandler := otelhttp.NewHandler(handler, "POST /payment")

	http.Handle("/payment", otelHandler)

	logger.Info("Starting payment-order-api on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("unable to start server", zap.Error(err))
		return
	}
}

type PaymentOrder struct {
	OriginCustomerID      string  `json:"originCustomerID"`
	DestinationCustomerID string  `json:"destinationCustomerID"`
	Value                 float64 `json:"value"`
}

func paymentHandler(tracer trace.Tracer, publisher *streaming.Publisher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx, span := tracer.Start(req.Context(), "payment handler")
		defer span.End()

		bodyContent, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Printf("error when reading body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var order PaymentOrder
		err = json.Unmarshal(bodyContent, &order)
		if err != nil {
			fmt.Printf("error when parsing body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		valid, validationError := validatePaymentOrder(&order)
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(validationError))
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

func validatePaymentOrder(order *PaymentOrder) (bool, string) {
	if order.Value <= 0 {
		return false, "Order should have a value greater than zero"
	}

	if order.OriginCustomerID == order.DestinationCustomerID {
		return false, "Order origin and destination should be different"
	}

	return true, ""
}
