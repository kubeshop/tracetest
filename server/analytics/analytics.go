package analytics

import (
	"fmt"
	"os"
)

type Tracker interface {
	Track(event string, props map[string]string) error
	Ready() bool
}

type noopTracker struct{}

func (t noopTracker) Track(event string, props map[string]string) error { return nil }
func (t noopTracker) Ready() bool                                       { return true }

var DefaultAnalyticsTracker *AnalyticsTracker

func NewAnalyticsTracker(enabled bool, serverID, appVersion, env string) (*AnalyticsTracker, error) {
	client, error := getClient(enabled, serverID, appVersion, env)
	if error != nil {
		return nil, error
	}

	DefaultAnalyticsTracker = &AnalyticsTracker{client: client}
	return DefaultAnalyticsTracker, nil
}

type AnalyticsTracker struct {
	client Tracker
}

func (at *AnalyticsTracker) Ready() bool {
	return at.client != nil && at.client.Ready()
}

func (at *AnalyticsTracker) SendEvent(name, category, clientID string, data *map[string]string) error {
	fmt.Printf(`sending event "%s" (%s)%s`, name, category, "\n")
	if !at.Ready() {
		err := fmt.Errorf("uninitialized client. Call analytics.Init")
		fmt.Printf(`could not send event "%s" (%s): %s%s`, name, category, err.Error(), "\n")
		return err
	}

	eventData := map[string]string{
		"category": category,
		"clientID": clientID,
	}

	if data != nil {
		for k, v := range *data {
			eventData[k] = v
		}
	}

	err := at.client.Track(name, eventData)
	if err != nil {
		fmt.Printf(`could not send event "%s" (%s): %s%s`, name, category, err.Error(), "\n")
	} else {
		fmt.Printf(`event sent "%s" (%s)%s`, name, category, "\n")
	}

	return err
}

func getClient(enabled bool, serverID, appVersion, env string) (Tracker, error) {
	// ga not enabled, use dumb settings
	if !enabled {
		defaultClient := noopTracker{}
		return defaultClient, nil
	}

	// setup an actual client
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get hostname: %w", err)
	}

	defaultClient := newSegmentTracker(hostname, serverID, appVersion, env)
	return defaultClient, nil
}
