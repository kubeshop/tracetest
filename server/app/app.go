package app

import (
	"context"
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
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/subscription"
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
	config           config.Config
	db               model.Repository
	traceDB          tracedb.TraceDB
	tracer           trace.Tracer
	executionContext *executionContext
}

func New(config config.Config, db model.Repository, tracedb tracedb.TraceDB, tracer trace.Tracer) (*App, error) {
	app := &App{
		config:  config,
		db:      db,
		traceDB: tracedb,
		tracer:  tracer,
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

func (a *App) GetConfig() config.Config {
	return a.config
}

func (a *App) SetConfig(config config.Config) {
	a.config = config
}

func (a *App) Start() error {
	a.executionContext = &executionContext{}
	fmt.Printf("Starting tracetest (version %s, env %s)\n", Version, Env)
	ctx := context.Background()

	serverID, isNewInstall, err := a.db.ServerID()

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

	applicationTracer, err := tracing.GetApplicationTracer(ctx, a.config)
	if err != nil {
		return fmt.Errorf("could not create trigger span tracer: %w", err)
	}

	triggerReg := trigger.NewRegsitry(a.tracer, applicationTracer)
	triggerReg.Add(trigger.HTTP())
	triggerReg.Add(trigger.GRPC())

	subscriptionManager := subscription.NewManager()

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(a.db)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	assertionExecutor := executor.NewAssertionExecutor(a.tracer)
	outputProcesser := executor.InstrumentedOutputProcessor(a.tracer)
	assertionRunner := executor.NewAssertionRunner(execTestUpdater, assertionExecutor, outputProcesser, subscriptionManager)
	assertionRunner.Start(5)
	defer assertionRunner.Stop()

	traceConversionConfig := traces.NewConversionConfig()
	// hardcoded for now. In the future we will get those values from the database
	traceConversionConfig.AddTimeFields(
		"tracetest.span.duration",
	)

	pollerExecutor := executor.NewPollerExecutor(a.config, a.tracer, execTestUpdater, a.traceDB)

	tracePoller := executor.NewTracePoller(pollerExecutor, execTestUpdater, assertionRunner, a.config.PoolingRetryDelay(), a.config.MaxWaitTimeForTraceDuration(), subscriptionManager)
	tracePoller.Start(5) // worker count. should be configurable
	defer tracePoller.Stop()

	runner := executor.NewPersistentRunner(triggerReg, a.db, execTestUpdater, tracePoller, a.tracer, subscriptionManager)
	runner.Start(5) // worker count. should be configurable
	defer runner.Stop()

	transactionRunner := executor.NewTransactionRunner(runner, a.db, subscriptionManager)
	transactionRunner.Start(5)
	defer transactionRunner.Stop()

	mappers := mappings.New(traceConversionConfig, comparator.DefaultRegistry(), a.db)

	controller := httpServer.NewController(a.db, runner, transactionRunner, assertionRunner, mappers)
	apiApiController := openapi.NewApiApiController(controller)
	customController := httpServer.NewCustomController(controller, apiApiController, openapi.DefaultErrorHandler, a.tracer)
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

	// Start otlp endpoint
	otlpServer := otlp.StartServer(21321, a.db)
	a.executionContext.otlpServer = otlpServer

	port := 11633
	if a.config.Server.HttpPort != 0 {
		port = a.config.Server.HttpPort
	}

	log.Printf("HTTP Server started")
	httpServer := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handlers.CompressHandler(router)}
	a.executionContext.httpServer = &httpServer

	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return nil
}

func (a *App) Stop() error {
	ctx := context.Background()
	err := a.executionContext.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("could not stop http server: %w", err)
	}

	a.executionContext.otlpServer.Stop()

	a.executionContext = nil

	return nil
}
