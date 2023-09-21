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

type CollectorOption func(*remoteIngesterConfig)

func WithTraceCache(traceCache TraceCache) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.traceCache = traceCache
	}
}

type collector struct {
	grpcServer stoppable
	httpServer stoppable
}

// Stop implements stoppable.
func (c *collector) Stop() {
	c.grpcServer.Stop()
	c.httpServer.Stop()
}

func Start(ctx context.Context, config Config, tracer trace.Tracer, opts ...CollectorOption) (stoppable, error) {
	ingesterConfig := remoteIngesterConfig{
		URL:   config.RemoteServerURL,
		Token: config.RemoteServerToken,
	}

	for _, opt := range opts {
		opt(&ingesterConfig)
	}

	ingester, err := newForwardIngester(ctx, config.BatchTimeout, ingesterConfig)
	if err != nil {
		return nil, fmt.Errorf("could not start local collector: %w", err)
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
		return nil, fmt.Errorf("could not start gRPC OTLP listener: %w", err)
	}

	if err = httpServer.Start(); err != nil {
		return nil, fmt.Errorf("could not start HTTP OTLP listener: %w", err)
	}

	return &collector{grpcServer: grpcServer, httpServer: httpServer}, nil
}

func onProcessTermination(callback func()) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		callback()
	}()
}
