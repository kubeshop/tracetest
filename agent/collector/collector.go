package collector

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/kubeshop/tracetest/server/otlp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var activeCollector Collector

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

func WithTraceMode(traceMode bool) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.traceMode = traceMode
	}
}

func WithStartRemoteServer(startRemoteServer bool) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.startRemoteServer = startRemoteServer
	}
}

func WithLogger(logger *zap.Logger) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.logger = logger
	}
}

func WithObserver(observer event.Observer) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.observer = observer
	}
}

func WithSensor(sensor sensors.Sensor) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.sensor = sensor
	}
}

func WithTraceModeForwarder(traceModeForwarder TraceModeForwarder) CollectorOption {
	return func(ric *remoteIngesterConfig) {
		ric.traceModeForwarder = traceModeForwarder
	}
}

type collector struct {
	grpcServer stoppable
	httpServer stoppable

	ingester ingester
}

type Collector interface {
	stoppable

	Statistics() Statistics
	ResetStatistics()

	SetSensor(sensors.Sensor)
}

// Stop implements stoppable.
func (c *collector) Stop() {
	c.grpcServer.Stop()
	c.httpServer.Stop()
}

func (c *collector) Statistics() Statistics {
	return c.ingester.Statistics()
}

func (c *collector) ResetStatistics() {
	c.ingester.ResetStatistics()
}

func (c *collector) SetSensor(sensor sensors.Sensor) {
	c.ingester.SetSensor(sensor)
}

func GetActiveCollector() Collector {
	return activeCollector
}

func Start(ctx context.Context, config Config, tracer trace.Tracer, opts ...CollectorOption) (Collector, error) {
	ingesterConfig := remoteIngesterConfig{
		URL:      config.RemoteServerURL,
		Token:    config.RemoteServerToken,
		logger:   zap.NewNop(),
		observer: event.NewNopObserver(),
		sensor:   sensors.NewSensor(),
	}

	for _, opt := range opts {
		opt(&ingesterConfig)
	}

	ingester, err := newForwardIngester(ctx, config.BatchTimeout, ingesterConfig, ingesterConfig.traceModeForwarder, ingesterConfig.startRemoteServer)
	if err != nil {
		return nil, fmt.Errorf("could not start local collector: %w", err)
	}

	grpcServer := otlp.NewGrpcServer(fmt.Sprintf("0.0.0.0:%d", config.GRPCPort), ingester, tracer)
	grpcServer.SetLogger(ingesterConfig.logger)
	httpServer := otlp.NewHttpServer(fmt.Sprintf("0.0.0.0:%d", config.HTTPPort), ingester)
	httpServer.SetLogger(ingesterConfig.logger)

	onProcessTermination(func() {
		ingesterConfig.logger.Debug("Stopping collector")
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

	activeCollector = &collector{grpcServer: grpcServer, httpServer: httpServer, ingester: ingester}
	return activeCollector, nil
}

func onProcessTermination(callback func()) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		callback()
	}()
}
