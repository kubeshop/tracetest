package collector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubeshop/tracetest/server/otlp"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	HTTPPort          int
	GRPCPort          int
	BatchTimeout      time.Duration
	RemoteServerURL   string
	RemoteServerToken string
}

func Start(ctx context.Context, config Config, tracer trace.Tracer) error {
	ingester, err := newForwardIngester(ctx, config.BatchTimeout, remoteIngesterConfig{
		URL:   config.RemoteServerURL,
		Token: config.RemoteServerToken,
	})
	if err != nil {
		return fmt.Errorf("could not start local collector: %w", err)
	}

	grpcServer := otlp.NewGrpcServer(fmt.Sprintf("0.0.0.0:%d", config.GRPCPort), ingester, tracer)
	httpServer := otlp.NewHttpServer(fmt.Sprintf("0.0.0.0:%d", config.HTTPPort), ingester)

	onProcessTermination(func() {
		grpcServer.Stop()
		httpServer.Stop()
		if stoppableIngester, ok := ingester.(stoppable); ok {
			stoppableIngester.Stop()
		}
	})

	if err = grpcServer.Start(); err != nil {
		return fmt.Errorf("could not start gRPC OTLP listener: %w", err)
	}

	if err = httpServer.Start(); err != nil {
		return fmt.Errorf("could not start HTTP OTLP listener: %w", err)
	}

	return nil
}

func onProcessTermination(callback func()) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		callback()
	}()
}
