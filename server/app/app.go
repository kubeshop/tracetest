package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/kubeshop/tracetest/analytics"
	"github.com/kubeshop/tracetest/config"
	"github.com/kubeshop/tracetest/executor"
	httpServer "github.com/kubeshop/tracetest/http"
	"github.com/kubeshop/tracetest/http/websocket"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/subscription"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/tracedb"
)

type App struct {
	config  config.Config
	db      testdb.Repository
	traceDB tracedb.TraceDB
}

func New(config config.Config, db testdb.Repository, tracedb tracedb.TraceDB) (*App, error) {
	app := &App{
		config:  config,
		db:      db,
		traceDB: tracedb,
	}

	return app, nil
}

func spaHandler(staticPath, indexPath string, tplVars map[string]string) http.HandlerFunc {
	var fileMatcher = regexp.MustCompile(`\.[a-zA-Z]*$`)
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (a *App) Start() error {
	err := analytics.Init(a.config.GA, "tracetest", os.Getenv("VERSION"))
	if err != nil {
		return err
	}

	ex, err := executor.New()
	if err != nil {
		return fmt.Errorf("could not create executor: %w", err)
	}

	subscriptionManager := subscription.NewManager()

	assertionRunner := executor.NewAssertionRunner(a.db)
	assertionRunner.Start(5)
	defer assertionRunner.Stop()

	tracePoller := executor.NewTracePoller(a.traceDB, a.db, a.db, a.config.MaxWaitTimeForTraceDuration(), subscriptionManager, assertionRunner)
	tracePoller.Start(5) // worker count. should be configurable
	defer tracePoller.Stop()

	runner := executor.NewPersistentRunner(ex, a.db, a.db, tracePoller)
	runner.Start(5) // worker count. should be configurable
	defer runner.Stop()

	controller := httpServer.NewController(a.traceDB, a.db, runner, assertionRunner)
	apiApiController := openapi.NewApiApiController(controller)

	router := openapi.NewRouter(apiApiController)

	router.PathPrefix("/").Handler(
		spaHandler(
			"./html",
			"index.html",
			map[string]string{
				"MeasurementId":    a.config.GA.MeasurementID,
				"AnalyticsEnabled": fmt.Sprintf("%t", a.config.GA.Enabled),
			},
		),
	)

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
