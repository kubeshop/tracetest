package collector

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	HTTPPort          int
	GRPCPort          int
	BatchTimeout      time.Duration
	RemoteServerURL   string
	RemoteServerToken string
}

func Start(config Config) {
	ingester := &forwardIngester{
		BatchTimeout: config.BatchTimeout,
		RemoteIngester: remoteIngesterConfig{
			URL:   config.RemoteServerURL,
			Token: config.RemoteServerToken,
		},
	}

	grpcServer := newGrpcServer(fmt.Sprintf("0.0.0.0:%d", config.GRPCPort), ingester)
	httpServer := newHttpServer(fmt.Sprintf("0.0.0.0:%d", config.HTTPPort), ingester)

	onProcessTermination(func() {
		grpcServer.Stop()
		httpServer.Stop()
	})

	go func() {
		err := grpcServer.Start()
		if err != nil {
			log.Println("ERROR: could not start gRPC OTLP listener: %w", err)
		}
	}()

	go func() {
		err := httpServer.Start()
		if err != nil {
			log.Println("ERROR: could not start HTTP OTLP listener: %w", err)
		}
	}()
}

func onProcessTermination(callback func()) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		callback()
	}()
}
