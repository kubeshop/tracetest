package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/config/demo"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	httpServer "github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/http/websocket"
	"github.com/kubeshop/tracetest/server/linter/analyzer"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testconnection"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testsuite"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/tracing"
	"github.com/kubeshop/tracetest/server/variableset"
	"go.opentelemetry.io/otel/trace"
)

var (
	Version = "dev"
	Env     = "dev"

	pgChannelName = "tracetest_queue"
)

var EmptyDemoEnabled []string

type App struct {
	cfg              *config.AppConfig
	provisioningFile string
	stopFns          []func()

	serverID string
}

func New(config *config.AppConfig) (*App, error) {
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
	sm.Subscribe(config.ResourceID, subscription.NewSubscriberFunction(
		func(m subscription.Message) error {
			configFromDB, ok := m.Content.(config.Config)
			if !ok {
				return fmt.Errorf("cannot read update to configFromDB. unexpected type %T", m.Content)
			}

			return app.initAnalytics(configFromDB)
		}),
	)
}

func (app *App) initAnalytics(configFromDB config.Config) error {
	return analytics.Init(configFromDB.IsAnalyticsEnabled(), app.serverID, Version, Env, app.cfg.AnalyticsServerKey(), app.cfg.AnalyticsFrontendKey())
}

var instanceID = id.GenerateID().String()

func (app *App) Start(opts ...appOption) error {
	// instanceID is a temprary ID for this instance of the server
	// it is regenerated on every start intentionally
	for _, opt := range opts {
		opt(app)
	}
	fmt.Println(app.Version())
	fmt.Println("Starting")
	ctx := context.Background()

	poolcfg, err := pgxpool.ParseConfig(app.cfg.PostgresConnString())
	if err != nil {
		return err
	}
	poolcfg.MaxConns = 20

	pool, err := pgxpool.NewWithConfig(context.Background(), poolcfg)
	if err != nil {
		return err
	}

	db, err := testdb.Connect(app.cfg.PostgresConnString())
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(80)

	testDB, err := testdb.Postgres(
		testdb.WithDB(db),
	)

	if err != nil {
		log.Fatal(err)
	}

	subscriptionManager := subscription.NewManager()
	app.subscribeToConfigChanges(subscriptionManager)

	configRepo := config.NewRepository(db, config.WithPublisher(subscriptionManager))
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

	demoRepo := demo.NewRepository(db)
	pollingProfileRepo := pollingprofile.NewRepository(db)
	dataStoreRepo := datastore.NewRepository(db)
	variableSetRepo := variableset.NewRepository(db)
	linterRepo := analyzer.NewRepository(db)
	testRepo := test.NewRepository(db)
	runRepo := test.NewRunRepository(db, test.NewCache(instanceID))
	testRunnerRepo := testrunner.NewRepository(db)
	tracesRepo := traces.NewTraceRepository(db)

	testSuiteRepository := testsuite.NewRepository(db, testRepo)
	testSuiteRunRepository := testsuite.NewRunRepository(db, runRepo)

	tracedbFactory := tracedb.Factory(tracesRepo)

	if app.cfg.OtlpServerEnabled() {
		eventEmitter := executor.NewEventEmitter(testDB, subscriptionManager)
		registerOtlpServer(app, tracesRepo, runRepo, eventEmitter, dataStoreRepo, tracer)
	}

	testPipeline := buildTestPipeline(
		pool,
		pollingProfileRepo,
		dataStoreRepo,
		linterRepo,
		testRunnerRepo,
		testDB,
		testRepo,
		runRepo,
		tracer,
		subscriptionManager,
		triggerRegistry,
		tracedbFactory,
		app.cfg,
	)
	testPipeline.Start()
	app.registerStopFn(func() {
		testPipeline.Stop()
	})

	testSuitePipeline := buildTestSuitePipeline(
		testSuiteRepository,
		testSuiteRunRepository,
		testPipeline,
		subscriptionManager,
	)

	testSuitePipeline.Start()
	app.registerStopFn(func() {
		testSuitePipeline.Stop()
	})

	dsTestListener := testconnection.NewListener()
	dsTestPipeline := buildDataStoreTestPipeline(
		pool,
		dsTestListener,
		tracer,
		tracedbFactory,
		app.cfg,
	)

	dsTestPipeline.Start()
	app.registerStopFn(func() {
		dsTestPipeline.Stop()
	})

	err = analytics.SendEvent("Server Started", "beacon", "", nil)
	if err != nil {
		return err
	}

	provisioner := provisioning.New()

	router, mappers := controller(app.cfg,
		tracer,

		testPipeline,
		testSuitePipeline,
		dsTestPipeline,

		testDB,
		testSuiteRepository,
		testSuiteRunRepository,
		testRepo,
		runRepo,
		variableSetRepo,
		tracedbFactory,
	)
	registerWSHandler(router, mappers, subscriptionManager)

	// use the analytics middleware on complete router
	router.Use(analyticsMW)

	// use the tenant middleware on complete router
	router.Use(middleware.TenantMiddleware)

	apiRouter := router.
		PathPrefix(app.cfg.ServerPathPrefix()).
		PathPrefix("/api").
		Subrouter()

	registerTestSuiteResource(testSuiteRepository, apiRouter, provisioner, tracer)
	registerConfigResource(configRepo, apiRouter, provisioner, tracer)
	registerPollingProfilesResource(pollingProfileRepo, apiRouter, provisioner, tracer)
	registerVariableSetResource(variableSetRepo, apiRouter, provisioner, tracer)
	registerDemosResource(demoRepo, apiRouter, provisioner, tracer)
	registerDataStoreResource(dataStoreRepo, apiRouter, provisioner, tracer)
	registerAnalyzer(linterRepo, apiRouter, provisioner, tracer)
	registerTestRunner(testRunnerRepo, apiRouter, provisioner, tracer)
	registerTestResource(testRepo, apiRouter, provisioner, tracer)

	isTracetestDev := os.Getenv("TRACETEST_DEV") != ""
	registerSPAHandler(router, app.cfg, configFromDB.IsAnalyticsEnabled(), serverID, isTracetestDev)

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

var (
	matchFirstCap     = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap       = regexp.MustCompile("([a-z0-9])([A-Z])")
	matchResourceName = regexp.MustCompile(`(\w)(\.)(\w)`)
)

func toWords(str string) string {
	if matchResourceName.MatchString(str) {
		return str
	}
	words := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	words = matchAllCap.ReplaceAllString(words, "${1} ${2}")
	return words
}

// Analytics global middleware
func analyticsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routeName := mux.CurrentRoute(r).GetName()
		machineID := r.Header.Get("x-client-id")
		source := r.Header.Get("x-source")
		eventProperties := r.Header.Get("x-event-properties")

		eventData := map[string]string{
			"source": source,
		}
		eventData = analytics.InjectProperties(eventData, eventProperties)

		if routeName != "" {
			analytics.SendEvent(toWords(routeName), "test", machineID, &eventData)
		}

		next.ServeHTTP(w, r)
	})
}

func registerSPAHandler(router *mux.Router, cfg httpServerConfig, analyticsEnabled bool, serverID string, isTracetestDev bool) {
	router.
		PathPrefix(cfg.ServerPathPrefix()).
		Handler(
			httpServer.SPAHandler(
				cfg,
				analyticsEnabled,
				serverID,
				Version,
				Env,
				isTracetestDev,
			),
		)
}

func registerOtlpServer(app *App, tracesRepo *traces.TraceRepository, runRepository test.RunRepository, eventEmitter executor.EventEmitter, dsRepo *datastore.Repository, tracer trace.Tracer) {
	ingester := otlp.NewIngester(tracesRepo, runRepository, eventEmitter, dsRepo, tracer)
	grpcOtlpServer := otlp.NewGrpcServer(":4317", ingester, tracer)
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

func registerAnalyzer(linterRepo *analyzer.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[analyzer.Linter](
		analyzer.ResourceName,
		analyzer.ResourceNamePlural,
		linterRepo,
		resourcemanager.DisableDelete(),
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerTestRunner(testRunnerRepo *testrunner.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[testrunner.TestRunner](
		testrunner.ResourceName,
		testrunner.ResourceNamePlural,
		testRunnerRepo,
		resourcemanager.DisableDelete(),
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerTestSuiteResource(repo *testsuite.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[testsuite.TestSuite](
		testsuite.TestSuiteResourceName,
		testsuite.TestSuiteResourceNamePlural,
		repo,
		resourcemanager.CanBeAugmented(),
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerConfigResource(configRepo *config.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[config.Config](
		config.ResourceName,
		config.ResourceNamePlural,
		configRepo,
		resourcemanager.DisableDelete(),
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerPollingProfilesResource(repository *pollingprofile.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[pollingprofile.PollingProfile](
		pollingprofile.ResourceName,
		pollingprofile.ResourceNamePlural,
		repository,
		resourcemanager.DisableDelete(),
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerVariableSetResource(repository *variableset.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[variableset.VariableSet](
		variableset.ResourceName,
		variableset.ResourceNamePlural,
		repository,
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerDemosResource(repository *demo.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[demo.Demo](
		demo.ResourceName,
		demo.ResourceNamePlural,
		repository,
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerDataStoreResource(repository *datastore.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[datastore.DataStore](
		datastore.ResourceName,
		datastore.ResourceNamePlural,
		repository,
		resourcemanager.WithTracer(tracer),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func registerTestResource(repository test.Repository, router *mux.Router, provisioner *provisioning.Provisioner, tracer trace.Tracer) {
	manager := resourcemanager.New[test.Test](
		test.ResourceName,
		test.ResourceNamePlural,
		repository,
		resourcemanager.WithTracer(tracer),
		resourcemanager.CanBeAugmented(),
	)
	manager.RegisterRoutes(router)
	provisioner.AddResourceProvisioner(manager)
}

func getTriggerRegistry(tracer, appTracer trace.Tracer) *trigger.Registry {
	triggerReg := trigger.NewRegistry(tracer, appTracer)
	triggerReg.Add(trigger.HTTP())
	triggerReg.Add(trigger.GRPC())
	triggerReg.Add(trigger.TRACEID())
	triggerReg.Add(trigger.Kafka())

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

	tracer trace.Tracer,

	testRunner *executor.TestPipeline,
	testSuitesRunner *executor.TestSuitesPipeline,

	dsTestRunner *testconnection.DataStoreTestPipeline,

	testRunEvents model.TestRunEventRepository,
	transactionRepo *testsuite.Repository,
	transactionRunRepo *testsuite.RunRepository,
	testRepo test.Repository,
	testRunRepo test.RunRepository,
	variablesetRepo *variableset.Repository,
	tracedbFactory tracedb.FactoryFunc,
) (*mux.Router, mappings.Mappings) {
	mappers := mappings.New(tracesConversionConfig(), comparator.DefaultRegistry())

	router := openapi.NewRouter(httpRouter(
		cfg,

		tracer,

		testRunner,
		testSuitesRunner,
		dsTestRunner,

		testRunEvents,
		transactionRepo,
		transactionRunRepo,
		testRepo,
		testRunRepo,
		variablesetRepo,
		tracedbFactory,

		mappers,
	))

	return router, mappers
}

func httpRouter(
	cfg httpServerConfig,

	tracer trace.Tracer,

	testRunner *executor.TestPipeline,
	testSuitesRunner *executor.TestSuitesPipeline,
	dsTestRunner *testconnection.DataStoreTestPipeline,

	testRunEvents model.TestRunEventRepository,
	testSuiteRepo *testsuite.Repository,
	testSuiteRunRepo *testsuite.RunRepository,
	testRepo test.Repository,
	testRunRepo test.RunRepository,
	variableSetRepo *variableset.Repository,
	tracedbFactory tracedb.FactoryFunc,

	mappers mappings.Mappings,
) openapi.Router {
	controller := httpServer.NewController(
		tracer,

		testRunner,
		testSuitesRunner,
		dsTestRunner,

		testRunEvents,
		testSuiteRepo,
		testSuiteRunRepo,
		testRepo,
		testRunRepo,
		variableSetRepo,

		tracedbFactory,
		mappers,
		Version,
	)
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
