package pages

import (
	"fmt"

	"github.com/kubeshop/tracetest/agent/ui/dashboard/components"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/models"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/rivo/tview"
)

const maxTestRuns = 25

type TestRunPage struct {
	*tview.Grid

	header      *components.Header
	testRunList *components.TestRunList

	renderScheduler components.RenderScheduler
	testRuns        []models.TestRun
}

func NewTestRunPage(renderScheduler components.RenderScheduler, sensor sensors.Sensor) *TestRunPage {
	p := &TestRunPage{
		Grid:            tview.NewGrid(),
		renderScheduler: renderScheduler,
		testRuns:        make([]models.TestRun, 0, 30),
	}

	p.header = components.NewHeader(renderScheduler, sensor)
	p.testRunList = components.NewTestRunList(renderScheduler)

	p.Grid.
		SetRows(10, 30, 30, 30).
		SetColumns(30, 0, 30).
		AddItem(p.header, 0, 0, 1, 3, 0, 0, false).
		AddItem(p.testRunList, 1, 0, 2, 3, 0, 0, true)

	sensor.On(events.NewTestRun, func(e sensors.Event) {
		var testRun models.TestRun
		err := e.Unmarshal(&testRun)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(p.testRuns) < maxTestRuns {
			p.testRuns = append(p.testRuns, testRun)
		} else {
			p.testRuns = append(p.testRuns[1:], testRun)
		}

		p.testRunList.SetTestRuns(p.testRuns)
	})

	sensor.On(events.UpdatedTestRun, func(e sensors.Event) {
		var testRun models.TestRun
		err := e.Unmarshal(&testRun)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for i, run := range p.testRuns {
			if run.TestID == testRun.TestID && run.RunID == testRun.RunID {
				p.testRuns[i] = testRun
			}
		}

		p.testRunList.SetTestRuns(p.testRuns)
	})

	return p
}

func (p *TestRunPage) Focus(delegate func(p tview.Primitive)) {
	delegate(p.testRunList)
}
