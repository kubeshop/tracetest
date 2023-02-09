/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/testdb"
)

var configFile = flag.String("config", "config.yaml", "path to the config file")

func main() {

	flag.Parse()
	cfg := loadConfig()
	db, err := testdb.Connect(cfg.PostgresConnString())
	if err != nil {
		log.Fatal(err)
	}

	appCfg := app.Config{
		Config: cfg,
	}

	appInstance, err := app.New(appCfg, db)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		wg.Done()
		appInstance.Stop()
		os.Exit(1)
	}()

	wg.Add(1)
	err = appInstance.Start()
	if err != nil {
		log.Fatal(err)
	}

	go watchChanges(func() {
		appCfg = app.Config{
			Config: loadConfig(),
		}
		appInstance.HotReload(appCfg)
	})

	wg.Wait()
}

func loadConfig() *config.Config {
	cfg, err := config.FromFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	return cfg

}

func watchChanges(updateFn func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("config updated:", event.Name)
					updateFn()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})

}
