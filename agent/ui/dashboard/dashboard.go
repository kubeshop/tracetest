package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/agent/ui/dashboard/components"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/models"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/pages"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/styles"
	"github.com/rivo/tview"
)

type Dashboard struct{}

func startUptimeCounter(sensor sensors.Sensor) {
	ticker := time.NewTicker(time.Second)
	fastTicker := time.NewTicker(50 * time.Millisecond)
	start := time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				sensor.Emit(events.UptimeChanged, time.Since(start).Round(time.Second))
			case <-fastTicker.C:
				sensor.Emit(events.NewTestRun, models.TestRun{TestID: "1", RunID: "1", Name: "my test", Type: "HTTP", Endpoint: "http://localhost:11633/api/tests", Status: "Awaiting Traces", When: time.Since(start)})
			}
		}
	}()
}

func StartDashboard(ctx context.Context) error {
	app := tview.NewApplication()
	tview.Styles.PrimitiveBackgroundColor = styles.HeaderBackgroundColor
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
