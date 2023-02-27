package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	httpServer "github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/http/websocket"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	Version = "dev"
	Env     = "dev"
)

var EmptyDemoEnabled []string

type App struct {
	cfg              *config.Config
	provisioningFile string
	stopFns          []func()
}

func New(config *config.Config) (*App, error) {
	app := &App{
		cfg: config,
	}

	return app, nil
}

func (app *App) Version() string {
	return fmt.Sprintf("tracetest-server %s (%s)", Version, Env)
}

func (app *App) Stop() {
	for _, fn := range app.stopFns {
		fn()
	}
}

func (app *App) registerStopFn(fn func()) {
	app.stopFns = append(app.stopFns, fn)
}

func (app *App) HotReload() {
	app.Stop()
	app.Start()
}

type appOption func(app *App)

func WithProvisioningFile(path string) appOption {
	return func(app *App) {
		app.provisioningFile = path
	}
}

func (app *App) provision(db model.Repository) {
	p := provisioning.New(db)

	var err error

	if app.provisioningFile != "" {
		log.Println("[provisioning] attempting file: ", app.provisioningFile)
		err = p.FromFile(app.provisioningFile)
		if err != nil {
			log.Fatalf("[provisioning] error: %s", err.Error())
		}
		fmt.Println("[Provisioning]: success")
		return
	}

	err = p.FromEnv()
	log.Println("[provisioning] attempting env var")
	if err != nil {
		if !errors.Is(err, provisioning.ErrEnvEmpty) {
			log.Fatalf("[provisioning] error: %s", err.Error())
		}
		log.Println("[provisioning] TRACETEST_PROVISIONING env var is empty")
	}
	fmt.Println("[Provisioning]: success")
}

func (app *App) Start(opts ...appOption) error {
	for _, opt := range opts {
		opt(app)
	}
	fmt.Println(app.Version())
	fmt.Println("Starting")
	ctx := context.Background()

	db, err := testdb.Connect(app.cfg.PostgresConnString())
	if err != nil {
		return err
	}

	testDB, err := testdb.Postgres(
		testdb.WithDB(db),
	)
	if err != nil {
		log.Fatal(err)
	}

	tracer, err := tracing.NewTracer(ctx, app.cfg)
	if err != nil {
		log.Fatal(err)
	}
	app.registerStopFn(func() {
		fmt.Println("stopping tracer")
		tracing.ShutdownTracer(ctx)
	})

	serverID, isNewInstall, err := testDB.ServerID()
	if err != nil {
		return err
	}

	err = analytics.Init(app.cfg.AnalyticsEnabled(), serverID, Version, Env)
	if err != nil {
		return err
	}

	fmt.Println("New install?", isNewInstall)
	if isNewInstall {
		err = analytics.SendEvent("Install", "beacon", "")
		if err != nil {
			return err
		}

		app.provision(testDB)

	}

	applicationTracer, err := tracing.GetApplicationTracer(ctx, app.cfg)
	if err != nil {
		return fmt.Errorf("could not create trigger span tracer: %w", err)
	}

	subscriptionManager := subscription.NewManager()
	triggerRegistry := getTriggerRegistry(tracer, applicationTracer)

	rf := newRunnerFacades(
		app.cfg,
		testDB,
		applicationTracer,
		tracer,
		subscriptionManager,
		triggerRegistry,
	)

	// worker count. should be configurable
	rf.tracePoller.Start(5)
	rf.runner.Start(5)
	rf.runner.Start(5)
	rf.transactionRunner.Start(5)
	rf.assertionRunner.Start(5)

	app.registerStopFn(func() {
		fmt.Println("stopping tracePoller")
		rf.tracePoller.Stop()
	})
	app.registerStopFn(func() {
		fmt.Println("stopping runner")
		rf.runner.Stop()
	})
	app.registerStopFn(func() {
		fmt.Println("stopping transactionRunner")
		rf.transactionRunner.Stop()
	})
	app.registerStopFn(func() {
		fmt.Println("stopping assertionRunner")
		rf.assertionRunner.Stop()
	})

	err = analytics.SendEvent("Server Started", "beacon", "")
	if err != nil {
		return err
	}

	otlpServer := otlp.NewServer(":21321", testDB)
	go otlpServer.Start()
	app.registerStopFn(func() {
		fmt.Println("stopping otlp server")
		otlpServer.Stop()
	})

	httpServer := newHttpServer(
		serverID,
		app.cfg,
		testDB,
		tracer,
		subscriptionManager,
		rf,
		triggerRegistry,
	)
	app.registerStopFn(func() {
		fmt.Println("stopping http server")
		httpServer.Shutdown(ctx)
	})

	go httpServer.ListenAndServe()
	log.Printf("HTTP Server started on %s", httpServer.Addr)

	return nil
}

type facadeConfig interface {
	PoolingRetryDelay() time.Duration
	PoolingMaxWaitTimeForTraceDuration() time.Duration
}

func newRunnerFacades(
	cfg facadeConfig,
	testDB model.Repository,
	appTracer trace.Tracer,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	triggerRegistry *trigger.Registry,
) *runnerFacade {

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(testDB)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	assertionRunner := executor.NewAssertionRunner(
		execTestUpdater,
		executor.NewAssertionExecutor(tracer),
		executor.InstrumentedOutputProcessor(tracer),
		subscriptionManager,
	)

	retryDelay := cfg.PoolingRetryDelay()
	maxWaitTime := cfg.PoolingMaxWaitTimeForTraceDuration()

	pollerExecutor := executor.NewPollerExecutor(
		retryDelay,
		maxWaitTime,
		tracer,
		execTestUpdater,
		tracedb.Factory(testDB),
		testDB,
	)

	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		execTestUpdater,
		assertionRunner,
		retryDelay,
		maxWaitTime,
		subscriptionManager,
	)

	runner := executor.NewPersistentRunner(
		triggerRegistry,
		testDB,
		execTestUpdater,
		tracePoller,
		tracer,
		subscriptionManager,
		tracedb.Factory(testDB),
		testDB,
	)

	transactionRunner := executor.NewTransactionRunner(runner, testDB, subscriptionManager)

	return &runnerFacade{
		runner:            runner,
		transactionRunner: transactionRunner,
		assertionRunner:   assertionRunner,
		tracePoller:       tracePoller,
	}
}

func getTriggerRegistry(tracer, appTracer trace.Tracer) *trigger.Registry {
	triggerReg := trigger.NewRegsitry(tracer, appTracer)
	triggerReg.Add(trigger.HTTP())
	triggerReg.Add(trigger.GRPC())
	triggerReg.Add(trigger.TRACEID())

	return triggerReg
}

type httpServerConfig interface {
	ServerPathPrefix() string
	ServerPort() int
	AnalyticsEnabled() bool
	DemoEnabled() []string
	DemoEndpoints() map[string]string
	ExperimentalFeatures() []string
}

func newHttpServer(
	serverID string,
	cfg httpServerConfig,
	testDB model.Repository,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	rf *runnerFacade,
	triggerRegistry *trigger.Registry,
) *http.Server {
	mappers := mappings.New(tracesConversionConfig(), comparator.DefaultRegistry(), testDB)

	router := openapi.NewRouter(httpRouter(cfg, testDB, tracer, rf, mappers, triggerRegistry))

	wsRouter := websocket.NewRouter()
	wsRouter.Add("subscribe", websocket.NewSubscribeCommandExecutor(subscriptionManager, mappers))
	wsRouter.Add("unsubscribe", websocket.NewUnsubscribeCommandExecutor(subscriptionManager))

	router.Handle("/ws", wsRouter.Handler())

	router.
		PathPrefix(cfg.ServerPathPrefix()).
		Handler(
			httpServer.SPAHandler(
				cfg,
				serverID,
				Version,
				Env,
			),
		)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort()),
		Handler: handlers.CompressHandler(router),
	}

	return httpServer
}

func httpRouter(
	cfg httpServerConfig,
	testDB model.Repository,
	tracer trace.Tracer,
	rf *runnerFacade,
	mappers mappings.Mappings,
	triggerRegistry *trigger.Registry,
) openapi.Router {
	controller := httpServer.NewController(testDB, tracedb.Factory(testDB), rf, mappers, triggerRegistry)
	apiApiController := openapi.NewApiApiController(controller)
	customController := httpServer.NewCustomController(controller, apiApiController, openapi.DefaultErrorHandler, tracer)
	httpRouter := customController

	if prefix := cfg.ServerPathPrefix(); prefix != "" {
		httpRouter = httpServer.NewPrefixedRouter(httpRouter, prefix)
	}

	return httpRouter
}

func tracesConversionConfig() traces.ConversionConfig {
	tcc := traces.NewConversionConfig()
	// hardcoded for now. In the future we will get those values from the database
	tcc.AddTimeFields(
		"tracetest.span.duration",
	)

	return tcc
}
