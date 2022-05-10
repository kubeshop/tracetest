package app

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/kubeshop/tracetest/analytics"
	"github.com/kubeshop/tracetest/config"
	"github.com/kubeshop/tracetest/executor"
	httpServer "github.com/kubeshop/tracetest/http"
	"github.com/kubeshop/tracetest/http/websocket"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/subscription"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/tracedb"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type App struct {
	config         config.Config
	db             testdb.Repository
	traceDB        tracedb.TraceDB
	tracerProvider *sdktrace.TracerProvider
	executor       executor.Executor
}

func NewApp(ctx context.Context, config config.Config, options ...Option) (*App, error) {
	app := &App{
		config: config,
	}

	for _, option := range options {
		err := option(app)
		if err != nil {
			return nil, err
		}
	}

	for _, defaultOption := range app.Options(ctx, config) {
		err := defaultOption(app)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	err := analytics.Init(a.config.GA, "tracetest", os.Getenv("VERSION"))
	if err != nil {
		return err
	}

	subscriptionManager := subscription.NewManager()

	tracePoller := executor.NewTracePoller(a.traceDB, a.db, a.config.MaxWaitTimeForTraceDuration(), subscriptionManager)
	tracePoller.Start(5) // worker count. should be configurable
	defer tracePoller.Stop()

	runner := executor.NewPersistentRunner(a.executor, a.db, a.db, tracePoller)
	runner.Start(5) // worker count. should be configurable
	defer runner.Stop()

	controller := httpServer.NewController(a.traceDB, a.db, runner)
	apiApiController := openapi.NewApiApiController(controller)

	router := openapi.NewRouter(apiApiController)
	router.Use(otelmux.Middleware("tracetest"))

	dir := "./html"
	fileServer := http.FileServer(http.Dir(dir))
	fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			serveIndex(w, dir+"/index.html")
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	err = analytics.CreateAndSendEvent("server_started_backend", "beacon")
	if err != nil {
		return err
	}

	go func() {
		wsRouter := websocket.NewRouter()
		wsRouter.Add("subscribe", websocket.NewSubscribeCommandExecutor(subscriptionManager))
		wsRouter.Add("unsubscribe", websocket.NewUnsubscribeCommandExecutor(subscriptionManager))
		log.Printf("WS Server started")
		wsRouter.ListenAndServe(":8081")
	}()

	log.Printf("HTTP Server started")
	log.Fatal(http.ListenAndServe(":8080", router))

	return nil
}

type gaParams struct {
	MeasurementId    string
	AnalyticsEnabled bool
}

func serveIndex(w http.ResponseWriter, path string) {
	templateData := gaParams{
		MeasurementId:    os.Getenv("GOOGLE_ANALYTICS_MEASUREMENT_ID"),
		AnalyticsEnabled: os.Getenv("ANALYTICS_ENABLED") == "true",
	}

	tpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err = tpl.Execute(w, templateData); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (a *App) Stop(ctx context.Context) {
	a.tracerProvider.Shutdown(ctx)
}

func (a *App) Options(ctx context.Context, config config.Config) []Option {
	options := make([]Option, 0)

	if a.db == nil {
		options = append(options, getDefaultDBOption(config))
	}

	if a.traceDB == nil {
		options = append(options, getDefaultTraceDBOption(config))
	}

	if a.tracerProvider == nil {
		options = append(options, getDefaultTracerProviderOption(ctx, config))
	}

	if a.executor == nil {
		options = append(options, getDefaultExecutorOption())
	}

	return options
}

func getDefaultDBOption(config config.Config) Option {
	return func(a *App) error {
		defaultDB, err := testdb.Postgres(
			testdb.WithDSN(config.PostgresConnString),
		)

		if err != nil {
			return fmt.Errorf("could not connect to db: %w", err)
		}

		a.db = defaultDB
		return nil
	}
}

func getDefaultTraceDBOption(config config.Config) Option {
	return func(a *App) error {
		traceDB, err := tracedb.New(config)
		if err != nil {
			return fmt.Errorf("could not connect to trace db: %w", err)
		}

		a.traceDB = traceDB
		return nil
	}
}

func getDefaultTracerProviderOption(ctx context.Context, config config.Config) Option {
	return func(a *App) error {
		tp, err := initOtelTracing(ctx)
		if err != nil {
			return fmt.Errorf("could not create OpenTelemetry tracer: %w", err)
		}

		a.tracerProvider = tp
		return nil
	}
}

func initOtelTracing(ctx context.Context) (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if endpoint == "" {
		endpoint = "opentelemetry-collector:4317"
		exporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
		if err != nil {
			return nil, err
		}
	} else {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		}
		exporter, err = otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
		if err != nil {
			return nil, err
		}
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{}))

	// Set standard attributes per semantic conventions
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tracetest"),
	)

	// Create and set the TraceProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}

func getDefaultExecutorOption() Option {
	return func(a *App) error {
		ex, err := executor.New()
		if err != nil {
			return fmt.Errorf("could not create executor: %w", err)
		}

		a.executor = ex
		return nil
	}

}
