package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	PokeAPI = "dev"
)

var EmptyDemoEnabled []string

type App struct {
	config  Config
	db      *sql.DB
	stopFns []func()
}

type Config struct {
	config.Config
	Migrations string
}

func New(config Config, db *sql.DB) (*App, error) {
	app := &App{
		config: config,
		db:     db,
	}

	return app, nil
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
	fmt.Printf("Starting tracetest (version %s, env %s)\n", Version, Env)
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

	err = analytics.Init(a.config.GA.Enabled, serverID, Version, Env)
	if err != nil {
		return err
	}

	if isNewInstall {
		err = analytics.SendEvent("Install", "beacon", "")
		if err != nil {
			return err
		}
	}

	applicationTracer, err := tracing.GetApplicationTracer(ctx, a.config.Config)
	if err != nil {
		return fmt.Errorf("could not create trigger span tracer: %w", err)
	}

	subscriptionManager := subscription.NewManager()

	rf := newRunnerFacades(
		a.config.Config,
		testDB,
		applicationTracer,
		tracer,
		subscriptionManager,
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
	)
	a.registerStopFn(func() {
		fmt.Println("stopping http server")
		httpServer.Shutdown(ctx)
	})

	go httpServer.ListenAndServe()
	log.Printf("HTTP Server started")

	return nil
}

func newRunnerFacades(
	conf config.Config,
	testDB model.Repository,
	appTracer trace.Tracer,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
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

	fallbackDS := model.DataStore{}
	dsc, err := conf.DataStore()
	if err == nil && dsc != nil {
		fallbackDS = model.DataStoreFromConfig(*dsc)
	}
	newTraceDB := tracedb.Factory(testDB)

	pollerExecutor := executor.NewPollerExecutor(
		conf.PoolingRetryDelay(),
		conf.MaxWaitTimeForTraceDuration(),
		tracer,
		execTestUpdater,
		tracedb.WithFallback(newTraceDB, fallbackDS),
		testDB,
	)

	tracePoller := executor.NewTracePoller(
		pollerExecutor,
		execTestUpdater,
		assertionRunner,
		conf.PoolingRetryDelay(),
		conf.MaxWaitTimeForTraceDuration(),
		subscriptionManager,
	)

	runner := executor.NewPersistentRunner(
		triggerRegistry(tracer, appTracer),
		testDB,
		execTestUpdater,
		tracePoller,
		tracer,
		subscriptionManager,
	)

	transactionRunner := executor.NewTransactionRunner(runner, testDB, subscriptionManager)

	return &runnerFacade{
		newTraceDB:        newTraceDB,
		runner:            runner,
		transactionRunner: transactionRunner,
		assertionRunner:   assertionRunner,
		tracePoller:       tracePoller,
	}
}

func triggerRegistry(tracer, appTracer trace.Tracer) *trigger.Registry {
	triggerReg := trigger.NewRegsitry(tracer, appTracer)
	triggerReg.Add(trigger.HTTP())
	triggerReg.Add(trigger.GRPC())

	return triggerReg
}

func newHttpServer(
	serverID string,
	conf config.Config,
	testDB model.Repository,
	tracer trace.Tracer,
	subscriptionManager *subscription.Manager,
	rf *runnerFacade,
) *http.Server {
	mappers := mappings.New(tracesConversionConfig(), comparator.DefaultRegistry(), testDB)

	router := openapi.NewRouter(httpRouter(conf, testDB, tracer, rf, mappers))

	wsRouter := websocket.NewRouter()
	wsRouter.Add("subscribe", websocket.NewSubscribeCommandExecutor(subscriptionManager, mappers))
	wsRouter.Add("unsubscribe", websocket.NewUnsubscribeCommandExecutor(subscriptionManager))

	router.Handle("/ws", wsRouter.Handler())

	router.
		PathPrefix(conf.Server.PathPrefix).
		Handler(
			httpServer.SPAHandler(
				conf,
				serverID,
				Version,
				Env,
			),
		)

	port := 11633
	if conf.Server.HttpPort != 0 {
		port = conf.Server.HttpPort
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handlers.CompressHandler(router),
	}

	return httpServer
}

func httpRouter(
	conf config.Config,
	testDB model.Repository,
	tracer trace.Tracer,
	rf *runnerFacade,
	mappers mappings.Mappings,
) openapi.Router {
	controller := httpServer.NewController(testDB, rf, mappers)
	apiApiController := openapi.NewApiApiController(controller)
	customController := httpServer.NewCustomController(controller, apiApiController, openapi.DefaultErrorHandler, tracer)
	httpRouter := customController

	if conf.Server.PathPrefix != "" {
		httpRouter = httpServer.NewPrefixedRouter(httpRouter, conf.Server.PathPrefix)
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
