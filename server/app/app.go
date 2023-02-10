package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
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
	config  Config
	db      *sql.DB
	stopFns []func()
}

type Config struct {
	*config.Config
	Migrations string
}

func New(config Config, db *sql.DB) (*App, error) {
	app := &App{
		config: config,
		db:     db,
	}

	return app, nil
}

func (a *App) Version() string {
	return fmt.Sprintf("tracetest-server %s (%s)", Version, Env)
}

func (a *App) Stop() {
	for _, fn := range a.stopFns {
		fn()
	}
}

func (a *App) registerStopFn(fn func()) {
	a.stopFns = append(a.stopFns, fn)
}

func (a *App) HotReload(c Config) {
	a.config = c
	a.Stop()
	a.Start()
}

func (a *App) Start() error {
	fmt.Println(a.Version())
	fmt.Println("Starting")
	ctx := context.Background()

	testDB, err := testdb.Postgres(
		testdb.WithDB(a.db),
		testdb.WithMigrations(a.config.Migrations),
	)
	if err != nil {
		log.Fatal(err)
	}

	tracer, err := tracing.NewTracer(ctx, a.config.Config)
	if err != nil {
		log.Fatal(err)
	}
	a.registerStopFn(func() {
		fmt.Println("stopping tracer")
		tracing.ShutdownTracer(ctx)
	})

	serverID, isNewInstall, err := testDB.ServerID()
	if err != nil {
		return err
	}

	err = analytics.Init(a.config.AnalyticsEnabled(), serverID, Version, Env)
	if err != nil {
		return err
	}

	fmt.Println("New install?", isNewInstall)
	if isNewInstall {
		err = analytics.SendEvent("Install", "beacon", "")
		if err != nil {
			return err
		}

		ensureFirstTimeDataSources(a.config.Config, testDB)
	}

	applicationTracer, err := tracing.GetApplicationTracer(ctx, a.config.Config)
	if err != nil {
		return fmt.Errorf("could not create trigger span tracer: %w", err)
	}

	subscriptionManager := subscription.NewManager()
	triggerRegistry := getTriggerRegistry(tracer, applicationTracer)

	rf := newRunnerFacades(
		a.config.Config,
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

	a.registerStopFn(func() {
		fmt.Println("stopping tracePoller")
		rf.tracePoller.Stop()
	})
	a.registerStopFn(func() {
		fmt.Println("stopping runner")
		rf.runner.Stop()
	})
	a.registerStopFn(func() {
		fmt.Println("stopping transactionRunner")
		rf.transactionRunner.Stop()
	})
	a.registerStopFn(func() {
		fmt.Println("stopping assertionRunner")
		rf.assertionRunner.Stop()
	})

	err = analytics.SendEvent("Server Started", "beacon", "")
	if err != nil {
		return err
	}

	otlpServer := otlp.NewServer(":21321", testDB)
	go otlpServer.Start()
	a.registerStopFn(func() {
		fmt.Println("stopping otlp server")
		otlpServer.Stop()
	})

	httpServer := newHttpServer(
		serverID,
		a.config.Config,
		testDB,
		tracer,
		subscriptionManager,
		rf,
		triggerRegistry,
	)
	a.registerStopFn(func() {
		fmt.Println("stopping http server")
		httpServer.Shutdown(ctx)
	})

	go httpServer.ListenAndServe()
	log.Printf("HTTP Server started")

	return nil
}

type dataStoreConfig interface {
	DataStore() (*config.TracingBackendDataStoreConfig, error)
}

func ensureFirstTimeDataSources(cfg dataStoreConfig, repo model.DataStoreRepository) {
	dsc, err := cfg.DataStore()
	if err != nil {
		panic(fmt.Errorf("cannot parse DataStore from config file: %w", err))
	}

	if dsc == nil {
		return
	}

	ds := model.DataStoreFromConfig(*dsc)
	ds.IsDefault = true

	if err := ds.Validate(); err != nil {
		panic(fmt.Errorf("invalid DataStore config from config file: %w", err))
	}

	ds, err = repo.CreateDataStore(context.Background(), ds)
	if err != nil {
		panic(fmt.Errorf("cannot persist DataStore from config file: %w", err))
	}

	fmt.Println("persisted initial DataStore from config file:")
	spew.Dump(ds)

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
