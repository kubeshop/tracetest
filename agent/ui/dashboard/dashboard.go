package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/agent/ui/dashboard/components"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/pages"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/rivo/tview"
)

type Dashboard struct{}

func startUptimeCounter(sensor sensors.Sensor) {
	ticker := time.NewTicker(time.Second)
	start := time.Now()
	go func() {
		for {
			<-ticker.C
			sensor.Emit(events.UptimeChanged, time.Since(start).Round(time.Second))
		}
	}()
}

func StartDashboard(ctx context.Context) error {
	app := tview.NewApplication()
	renderScheduler := components.NewRenderScheduler(app)
	sensor := sensors.NewSensor()

	startUptimeCounter(sensor)

	router := NewRouter()
	router.AddAndSwitchToPage("home", pages.NewTestRunPage(renderScheduler, sensor))

	if err := app.SetRoot(router, true).SetFocus(router).Run(); err != nil {
		return fmt.Errorf("failed to start dashboard: %w", err)
	}

	return nil
}
