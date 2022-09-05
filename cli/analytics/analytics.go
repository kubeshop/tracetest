package analytics

import (
	"fmt"

	"github.com/denisbrodbeck/machineid"
	"github.com/kubeshop/tracetest/cli/config"
	segment "github.com/segmentio/analytics-go/v3"
)

var (
	SecretKey = ""
	client    segment.Client
	mid       string
)

func Init(conf config.Config) {
	if !conf.AnalyticsEnabled {
		return
	}

	client = segment.New(SecretKey)
	id, err := machineid.ProtectedID("tracetest")
	if err != nil {
		panic(fmt.Errorf("could not get machineID: %w", err))
	}
	mid = id

	client.Enqueue(segment.Identify{
		Traits: segment.NewTraits().
			Set("source", "cli").
			Set("clientID", mid).
			Set("env", config.Env).
			Set("appVersion", config.Version),
	})
}

func Track(name, category string, props map[string]string) error {
	if client == nil {
		return nil
	}

	p := segment.NewProperties().
		Set("source", "cli").
		Set("clientID", mid).
		Set("env", config.Env).
		Set("appVersion", config.Version).
		Set("category", category)

	for k, v := range props {
		p = p.Set(k, v)
	}

	err := client.Enqueue(segment.Track{
		Event:      name,
		UserId:     mid,
		Properties: p,
	})

	return err
}

func Close() {
	if client == nil {
		return
	}

	err := client.Close()
	if err != nil {
		panic(err)
	}
}
