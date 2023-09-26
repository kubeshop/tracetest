package analytics

import (
	segment "github.com/segmentio/analytics-go/v3"
)

var (
	SecretKey   = ""
	FrontendKey = ""
)

func newSegmentTracker(hostname, serverID, appVersion, env, secretKey, frontendKey string) Tracker {
	if secretKey != "" {
		SecretKey = secretKey
	}

	if frontendKey != "" {
		FrontendKey = frontendKey
	}

	client, _ := segment.NewWithConfig(SecretKey, segment.Config{
		BatchSize: 1,
	})

	client.Enqueue(segment.Identify{
		UserId: serverID,
		Traits: segment.NewTraits().
			Set("source", "server").
			Set("serverID", serverID).
			Set("env", env).
			Set("appVersion", appVersion).
			Set("hostname", hostname),
		Context: &segment.Context{
			Direct: true,
		},
	})
	return segmentTracker{
		client:     client,
		env:        env,
		appVersion: appVersion,
		hostname:   hostname,
		serverID:   serverID,
	}
}

type segmentTracker struct {
	client     segment.Client
	env        string
	appVersion string
	hostname   string
	serverID   string
}

func (t segmentTracker) Ready() bool {
	return t.appVersion != "" &&
		t.hostname != "" &&
		t.serverID != "" &&
		t.env != ""

}

func (t segmentTracker) Track(name string, props map[string]string) error {
	p := segment.NewProperties().
		Set("source", "server").
		Set("serverID", t.serverID).
		Set("env", t.env).
		Set("appVersion", t.appVersion).
		Set("hostname", t.hostname)

	for k, v := range props {
		p = p.Set(k, v)
	}

	return t.client.Enqueue(segment.Track{
		Event:      name,
		UserId:     t.serverID,
		Properties: p,
		Context: &segment.Context{
			Direct: true,
		},
	})
}
