package pages

import (
	"github.com/kubeshop/tracetest/agent/ui/dashboard/components"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/rivo/tview"
)

type TestRunPage struct {
	*tview.Grid

	renderScheduler components.RenderScheduler
}

func NewTestRunPage(renderScheduler components.RenderScheduler, sensor sensors.Sensor) *TestRunPage {
	p := &TestRunPage{
		Grid:            tview.NewGrid(),
		renderScheduler: renderScheduler,
	}

	p.Grid.
		SetRows(10, 90).
		SetColumns(30, 0, 30).
		AddItem(components.NewHeader(renderScheduler, sensor), 0, 0, 1, 3, 0, 0, true).
		AddItem(components.NewTestRunList(), 1, 0, 1, 3, 0, 0, true)

	return p
}
