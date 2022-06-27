package app

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/kubeshop/tracetest/server/analytics"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	httpServer "github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/http/websocket"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

var Version = ""

type App struct {
	config  config.Config
	db      model.Repository
	traceDB tracedb.TraceDB
	tracer  trace.Tracer
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

func (a *App) Start() error {
	fmt.Println("Starting tracetest", Version)

	serverID, isNewInstall, err := a.db.ServerID()
	if err != nil {
		return err
	}

	err = analytics.Init(a.config.GA.Enabled, serverID, "tracetest", Version)
	if err != nil {
		return err
	}

	if isNewInstall {
		err = analytics.CreateAndSendEvent("install_server", "beacon")
		if err != nil {
			return err
		}
	}

	ex, err := executor.NewTriggerer(a.tracer)
	if err != nil {
		return fmt.Errorf("could not create executor: %w", err)
	}

	subscriptionManager := subscription.NewManager()

	execTestUpdater := (executor.CompositeUpdater{}).
		Add(executor.NewDBUpdater(a.db)).
		Add(executor.NewSubscriptionUpdater(subscriptionManager))

	assertionExecutor := executor.NewAssertionExecutor(a.tracer)

	assertionRunner := executor.NewAssertionRunner(execTestUpdater, assertionExecutor)
	assertionRunner.Start(5)
	defer assertionRunner.Stop()

	pollerExecutor := executor.NewPollerExecutor(a.tracer, execTestUpdater, a.traceDB, a.config.PoolingRetryDelay(), a.config.MaxWaitTimeForTraceDuration())

	tracePoller := executor.NewTracePoller(pollerExecutor, execTestUpdater, assertionRunner, a.tracer, a.config.PoolingRetryDelay(), a.config.MaxWaitTimeForTraceDuration())
	tracePoller.Start(5) // worker count. should be configurable
	defer tracePoller.Stop()

	runner := executor.NewPersistentRunner(ex, a.db, execTestUpdater, tracePoller)
	runner.Start(5) // worker count. should be configurable
	defer runner.Stop()

	controller := httpServer.NewController(a.db, runner, assertionRunner)
	apiApiController := openapi.NewApiApiController(controller)
	customController := httpServer.NewCustomController(controller, apiApiController, openapi.DefaultErrorHandler, a.tracer)
	httpRouter := customController
	if a.config.Server.PathPrefix != "" {
		httpRouter = httpServer.NewPrefixedRouter(httpRouter, a.config.Server.PathPrefix)
	}

	router := openapi.NewRouter(httpRouter)

	wsRouter := websocket.NewRouter()
	wsRouter.Add("subscribe", websocket.NewSubscribeCommandExecutor(subscriptionManager))
	wsRouter.Add("unsubscribe", websocket.NewUnsubscribeCommandExecutor(subscriptionManager))

	router.Handle("/ws", wsRouter.Handler())

	router.PathPrefix(a.config.Server.PathPrefix).Handler(
		spaHandler(
			a.config.Server.PathPrefix,
			"./html",
			"index.html",
			map[string]string{
				"MeasurementId":    analytics.MeasurementID,
				"AnalyticsEnabled": fmt.Sprintf("%t", a.config.GA.Enabled),
				"ServerPathPrefix": fmt.Sprintf("%s/", a.config.Server.PathPrefix),
			},
		),
	)

	err = analytics.CreateAndSendEvent("server_started_backend", "beacon")
	if err != nil {
		return err
	}

	port := 8080
	if a.config.Server.HttpPort != 0 {
		port = a.config.Server.HttpPort
	}

	log.Printf("HTTP Server started")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handlers.CompressHandler(router)))

	return nil
}
