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

	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
)

var cfg = flag.String("config", "config.yaml", "path to the config file")

func main() {
	flag.Parse()
	cfg, err := config.FromFile(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	testDB, err := testdb.Postgres(testdb.WithDSN(cfg.PostgresConnString))
	if err != nil {
		log.Fatal(err)
	}

	traceDB, err := tracedb.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.New(cfg, testDB, traceDB)
	if err != nil {
		log.Fatal(err)
	}

	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
