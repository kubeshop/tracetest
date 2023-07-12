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

func Init(enabled bool, serverID, appVersion, env string) error {
	// ga not enabled, use dumb settings
	if !enabled {
		defaultClient = noopTracker{}
		return nil
	}

	// setup an actual client
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("could not get hostname: %w", err)
	}

	defaultClient = newSegmentTracker(hostname, serverID, appVersion, env)

	return nil
}

var defaultClient Tracker

func Ready() bool {
	return defaultClient != nil && defaultClient.Ready()
}

func SendEvent(name, category, clientID string, data *map[string]string) error {
	fmt.Printf(`sending event "%s" (%s)%s`, name, category, "\n")
	if !Ready() {
		err := fmt.Errorf("uninitalized client. Call analytics.Init")
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

	err := defaultClient.Track(name, eventData)
	if err != nil {
		fmt.Printf(`could not send event "%s" (%s): %s%s`, name, category, err.Error(), "\n")
	} else {
		fmt.Printf(`event sent "%s" (%s)%s`, name, category, "\n")
	}

	return err
}
