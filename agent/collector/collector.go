package collector

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubeshop/tracetest/server/otlp"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/semaphore"
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

	var semaphore = semaphore.NewWeighted(2)
	semaphore.Acquire(ctx, 2)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	go func() {
		semaphore.Release(1)
		err = grpcServer.Start()
		if err != nil {
			log.Println("ERROR: could not start gRPC OTLP listener: %w", err)
		}
	}()

	go func() {
		semaphore.Release(1)
		err = httpServer.Start()
		if err != nil {
			log.Println("ERROR: could not start HTTP OTLP listener: %w", err)
		}
	}()

	// Wait until semaphore is released
	semaphore.Acquire(ctx, 2)

	return err
}

func onProcessTermination(callback func()) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		callback()
	}()
}
