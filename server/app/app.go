package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/config/demoresource"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	httpServer "github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/http/websocket"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
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

	serverID string
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

func provision(provisioner *provisioning.Provisioner, file string) {
	var err error

	if file != "" {
		log.Println("[provisioning] attempting file: ", file)
		err = provisioner.FromFile(file)
		if err != nil {
			log.Fatalf("[provisioning] error: %s", err.Error())
		}
		fmt.Println("[Provisioning]: success")
		return
	}

	err = provisioner.FromEnv()
	log.Println("[provisioning] attempting env var")
	if err != nil {
		if !errors.Is(err, provisioning.ErrEnvEmpty) {
			log.Fatalf("[provisioning] error: %s", err.Error())
		}
		log.Println("[provisioning] TRACETEST_PROVISIONING env var is empty")
	}
	fmt.Println("[Provisioning]: success")
}

func (app *App) subscribeToConfigChanges(sm *subscription.Manager) {
	sm.Subscribe(configresource.ResourceID, subscription.NewSubscriberFunction(
		func(m subscription.Message) error {
			configFromDB, ok := m.Content.(configresource.Config)
			if !ok {
				return fmt.Errorf("cannot read update to configFromDB. unexpected type %T", m.Content)
			}

			return app.initAnalytics(configFromDB)
		}),
	)
}

func (app *App) initAnalytics(configFromDB configresource.Config) error {
	return analytics.Init(configFromDB.IsAnalyticsEnabled(), app.serverID, Version, Env)
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

	subscriptionManager := subscription.NewManager()
	app.subscribeToConfigChanges(subscriptionManager)

	configRepo := configresource.NewRepository(db, configresource.WithPublisher(subscriptionManager))
	configFromDB := configRepo.Current(ctx)

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
	app.serverID = serverID

	err = app.initAnalytics(configFromDB)
	if err != nil {
		return err
	}

	fmt.Println("New install?", isNewInstall)
	if isNewInstall {
		err = analytics.SendEvent("Install", "beacon", "", nil)
		if err != nil {
			return err
		}
	}

	applicationTracer, err := tracing.GetApplicationTracer(ctx, app.cfg)
	if err != nil {
		return fmt.Errorf("could not create trigger span tracer: %w", err)
	}

	triggerRegistry := getTriggerRegistry(tracer, applicationTracer)

	pollingProfileRepo := pollingprofile.NewRepository(db)
	environmentRepo := environment.NewRepository(db)
	dataStoreRepo := datastoreresource.NewRepository(db)

	eventEmitter := executor.NewEventEmitter(testDB, subscriptionManager)
	registerOtlpServer(app, testDB, eventEmitter, dataStoreRepo)

	rf := newRunnerFacades(
		pollingProfileRepo,
		dataStoreRepo,
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

	err = analytics.SendEvent("Server Started", "beacon", "", nil)
	if err != nil {
		return err
	}

	provisioner := provisioning.New()

	router, mappers := controller(app.cfg, testDB, tracer, rf, triggerRegistry)
	registerWSHandler(router, mappers, subscriptionManager)

	apiRouter := router.PathPrefix("/api").Subrouter()
	registerConfigResource(configRepo, apiRouter, db, provisioner)

	registerPollingProfilesResource(pollingProfileRepo, apiRouter, db, provisioner)
	registerEnvironmentResource(&environmentRepo, apiRouter, db, provisioner)

	demoRepo := demoresource.NewRepository(db)
	registerDemosResource(demoRepo, apiRouter, db, provisioner)

	registerDataStoreResource(dataStoreRepo, apiRouter, db, provisioner)

	registerSPAHandler(router, app.cfg, configFromDB.IsAnalyticsEnabled(), serverID)

	if isNewInstall {
		provision(provisioner, app.provisioningFile)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.cfg.ServerPort()),
		Handler: handlers.CompressHandler(router),
	}

	app.registerStopFn(func() {
		fmt.Println("stopping http server")
		httpServer.Shutdown(ctx)
	})

	go httpServer.ListenAndServe()
	log.Printf("HTTP Server started on %s", httpServer.Addr)

	return nil
}

func registerSPAHandler(router *mux.Router, cfg httpServerConfig, analyticsEnabled bool, serverID string) {
	router.
		PathPrefix(cfg.ServerPathPrefix()).
		Handler(
			httpServer.SPAHandler(
				cfg,
				analyticsEnabled,
				serverID,
				Version,
				Env,
			),
		)
}

func registerOtlpServer(app *App, testDB model.Repository, eventEmitter executor.EventEmitter, dsRepo *datastoreresource.Repository) {
	ingester := otlp.NewIngester(testDB, eventEmitter, dsRepo)
	grpcOtlpServer := otlp.NewGrpcServer(":4317", ingester)
	httpOtlpServer := otlp.NewHttpServer(":4318", ingester)
	go grpcOtlpServer.Start()
	go httpOtlpServer.Start()

	fmt.Println("OTLP server started on :4317 (grpc) and :4318 (http)")

	app.registerStopFn(func() {
		fmt.Println("stopping otlp server")
		grpcOtlpServer.Stop()
		httpOtlpServer.Stop()
	})
}

func registerConfigResource(configRepo *configresource.Repository, router *mux.Router, db *sql.DB, provisioner *provisioning.Provisioner) {
	manager := resourcemanager.New[configresource.Config](
		configresource.ResourceName,
		configresource.ResourceNamePlural,
		configRepo,
		resourcemanager.WithOperations(configresource.Operations...),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerPollingProfilesResource(repository *pollingprofile.Repository, router *mux.Router, db *sql.DB, provisioner *provisioning.Provisioner) {
	manager := resourcemanager.New[pollingprofile.PollingProfile](
		pollingprofile.ResourceName,
		pollingprofile.ResourceNamePlural,
		repository,
		resourcemanager.WithOperations(pollingprofile.Operations...),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerEnvironmentResource(repository *environment.Repository, router *mux.Router, db *sql.DB, provisioner *provisioning.Provisioner) {
	manager := resourcemanager.New[environment.Environment](
		environment.ResourceName,
		environment.ResourceNamePlural,
		repository,
		resourcemanager.WithOperations(environment.Operations...),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerDemosResource(repository *demoresource.Repository, router *mux.Router, db *sql.DB, provisioner *provisioning.Provisioner) {
	manager := resourcemanager.New[demoresource.Demo](
		demoresource.ResourceName,
		demoresource.ResourceNamePlural,
		repository,
		resourcemanager.WithOperations(demoresource.Operations...),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerDataStoreResource(repository *datastoreresource.Repository, router *mux.Router, db *sql.DB, provisioner *provisioning.Provisioner) {
	manager := resourcemanager.New[datastoreresource.DataStore](
		datastoreresource.ResourceName,
		datastoreresource.ResourceNamePlural,
		repository,
		resourcemanager.WithOperations(datastoreresource.Operations...),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
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
	DemoEnabled() []string
	DemoEndpoints() map[string]string
	ExperimentalFeatures() []string
}

func registerWSHandler(router *mux.Router, mappers mappings.Mappings, subscriptionManager *subscription.Manager) {
	wsRouter := websocket.NewRouter()
	wsRouter.Add("subscribe", websocket.NewSubscribeCommandExecutor(subscriptionManager, mappers))
	wsRouter.Add("unsubscribe", websocket.NewUnsubscribeCommandExecutor(subscriptionManager))

	router.Handle("/ws", wsRouter.Handler())
}

func controller(
	cfg httpServerConfig,
	testDB model.Repository,
	tracer trace.Tracer,
	rf *runnerFacade,
	triggerRegistry *trigger.Registry,
) (*mux.Router, mappings.Mappings) {
	mappers := mappings.New(tracesConversionConfig(), comparator.DefaultRegistry(), testDB)

	router := openapi.NewRouter(httpRouter(cfg, testDB, tracer, rf, mappers, triggerRegistry))

	return router, mappers
}

func httpRouter(
	cfg httpServerConfig,
	testDB model.Repository,
	tracer trace.Tracer,
	rf *runnerFacade,
	mappers mappings.Mappings,
	triggerRegistry *trigger.Registry,
) openapi.Router {
	controller := httpServer.NewController(testDB, tracedb.Factory(testDB), rf, mappers, triggerRegistry, tracer)
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
