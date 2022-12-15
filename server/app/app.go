package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	httpServer "github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/http/websocket"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/kubeshop/tracetest/server/tracing"
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

func jsonEscape(text any) string {
	initial, err := json.Marshal(text)
	if err != nil {
		panic(err)
	}

	encoded, err := json.Marshal(string(initial))
	if err != nil {
		panic(err)
	}

	formatted := string(encoded)
	return strings.Trim(formatted, `"`)
}

func spaHandler(prefix, staticPath, indexPath string, tplVars map[string]string) http.HandlerFunc {
	var fileMatcher = regexp.MustCompile(`\.[a-zA-Z]*$`)
	handler := func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			tpl, err := template.ParseFiles(filepath.Join(staticPath, indexPath))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			if err = tpl.Execute(w, tplVars); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		} else {
			http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
		}
	}

	return http.StripPrefix(prefix, http.HandlerFunc(handler)).ServeHTTP
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

	traceDB, err := tracedb.New(a.config.Config, testDB)
	if err != nil {
		log.Fatal(err)
	}
	a.registerStopFn(func() {
		fmt.Println("stopping traceDB")
		traceDB.Close()
	})

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

	triggerReg := trigger.NewRegsitry(tracer, applicationTracer)
	triggerReg.Add(trigger.HTTP())
	triggerReg.Add(trigger.GRPC())

	subscriptionManager := subscription.NewManager()

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(testDB)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	assertionExecutor := executor.NewAssertionExecutor(tracer)
	outputProcesser := executor.InstrumentedOutputProcessor(tracer)
	assertionRunner := executor.NewAssertionRunner(execTestUpdater, assertionExecutor, outputProcesser, subscriptionManager)
	assertionRunner.Start(5)
	a.registerStopFn(func() {
		fmt.Println("stopping assertionRunner")
		assertionRunner.Stop()
	})

	traceConversionConfig := traces.NewConversionConfig()
	// hardcoded for now. In the future we will get those values from the database
	traceConversionConfig.AddTimeFields(
		"tracetest.span.duration",
	)

	pollerExecutor := executor.NewPollerExecutor(a.config.Config, tracer, execTestUpdater, traceDB)

	tracePoller := executor.NewTracePoller(pollerExecutor, execTestUpdater, assertionRunner, a.config.PoolingRetryDelay(), a.config.MaxWaitTimeForTraceDuration(), subscriptionManager)
	tracePoller.Start(5) // worker count. should be configurable
	a.registerStopFn(func() {
		fmt.Println("stopping tracePoller")
		tracePoller.Stop()
	})

	runner := executor.NewPersistentRunner(triggerReg, testDB, execTestUpdater, tracePoller, tracer, subscriptionManager)
	runner.Start(5) // worker count. should be configurable
	a.registerStopFn(func() {
		fmt.Println("stopping runner")
		runner.Stop()
	})

	transactionRunner := executor.NewTransactionRunner(runner, testDB, subscriptionManager)
	transactionRunner.Start(5)
	a.registerStopFn(func() {
		fmt.Println("stopping transactionRunner")
		transactionRunner.Stop()
	})

	mappers := mappings.New(traceConversionConfig, comparator.DefaultRegistry(), testDB)

	controller := httpServer.NewController(testDB, runner, transactionRunner, assertionRunner, mappers)
	apiApiController := openapi.NewApiApiController(controller)
	customController := httpServer.NewCustomController(controller, apiApiController, openapi.DefaultErrorHandler, tracer)
	httpRouter := customController
	if a.config.Server.PathPrefix != "" {
		httpRouter = httpServer.NewPrefixedRouter(httpRouter, a.config.Server.PathPrefix)
	}

	router := openapi.NewRouter(httpRouter)

	wsRouter := websocket.NewRouter()
	wsRouter.Add("subscribe", websocket.NewSubscribeCommandExecutor(subscriptionManager, mappers))
	wsRouter.Add("unsubscribe", websocket.NewUnsubscribeCommandExecutor(subscriptionManager))

	router.Handle("/ws", wsRouter.Handler())

	router.PathPrefix(a.config.Server.PathPrefix).Handler(
		spaHandler(
			a.config.Server.PathPrefix,
			"./html",
			"index.html",
			map[string]string{
				"AnalyticsKey":         analytics.FrontendKey,
				"AnalyticsEnabled":     fmt.Sprintf("%t", a.config.GA.Enabled),
				"ServerPathPrefix":     fmt.Sprintf("%s/", a.config.Server.PathPrefix),
				"ServerID":             serverID,
				"AppVersion":           Version,
				"Env":                  Env,
				"DemoEnabled":          jsonEscape(a.config.Demo.Enabled),
				"DemoEndpoints":        jsonEscape(a.config.Demo.Endpoints),
				"ExperimentalFeatures": jsonEscape(a.config.ExperimentalFeatures),
			},
		),
	)

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

	port := 11633
	if a.config.Server.HttpPort != 0 {
		port = a.config.Server.HttpPort
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handlers.CompressHandler(router),
	}
	a.registerStopFn(func() {
		fmt.Println("stopping http server")
		httpServer.Shutdown(ctx)
	})

	log.Printf("HTTP Server started")
	go httpServer.ListenAndServe()

	return nil
}
