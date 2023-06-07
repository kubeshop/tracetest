package analytics

import (
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/kubeshop/tracetest/cli/config"
	segment "github.com/segmentio/analytics-go/v3"
)

var (
	SecretKey = ""
	client    segment.Client
	mid       string
)

func ClientID() string {
	return mid
}

func Init(conf config.Config) {
	if !conf.AnalyticsEnabled || os.Getenv("TRACETEST_DEV") != "" {
		// non-empty TRACETEST_DEV variable means it's running by a dev,
		// and we should totally ignore analytics
		return
	}

	client, _ = segment.NewWithConfig(SecretKey, segment.Config{
		BatchSize: 1,
	})

	id, err := machineid.ProtectedID("tracetest")
	if err == nil {
		// only use id if available.
		mid = id
	} // ignore errors and continue with an empty ID if necessary

	client.Enqueue(segment.Identify{
		UserId: mid,
		Traits: segment.NewTraits().
			Set("source", "cli").
			Set("clientID", mid).
			Set("env", config.Env).
			Set("appVersion", config.Version),
		Context: &segment.Context{
			Direct: true,
		},
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
		Context: &segment.Context{
			Direct: true,
		},
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
