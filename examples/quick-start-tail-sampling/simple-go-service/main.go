package main

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kubeshop/tracetest/examples/quick-start-tail-sampling/simple-go-service/config"
	"github.com/kubeshop/tracetest/examples/quick-start-tail-sampling/simple-go-service/telemetry"
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

	wrappedHeartbeatHandler := getWrappedEndpointHandler("GET /heartbeat", heartbeatHandler(tracer, logger))
	http.Handle("/heartbeat", wrappedHeartbeatHandler)

	wrappedBusinessLogicHandler := getWrappedEndpointHandler("GET /businessLogic", businessLogicHandler(tracer, logger))
	http.Handle("/businessLogic", wrappedBusinessLogicHandler)

	logger.Info("Starting producer-api on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("unable to start server", zap.Error(err))
		return
	}
}

func getWrappedEndpointHandler(handlerName string, handler func(http.ResponseWriter, *http.Request)) http.Handler {
	httpHandler := http.HandlerFunc(handler)
	return otelhttp.NewHandler(httpHandler, handlerName)
}

func heartbeatHandler(tracer trace.Tracer, logger *zap.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, span := tracer.Start(req.Context(), "executing heartbeat")
		defer span.End()

		w.Write([]byte("I'm alive!"))
	}
}

func businessLogicHandler(tracer trace.Tracer, logger *zap.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logger.Info("GET /businessLogic")
		for k, v := range req.Header {
			logger.Info(fmt.Sprintf("Http Header -  '%s': %s", k, v))
		}

		_, span := tracer.Start(req.Context(), "doing some business rule")
		defer span.End()

		w.Write([]byte("Done!"))
	}
}
