package components

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/agent/ui/dashboard/models"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/styles"
	"github.com/rivo/tview"
)

var headers = []string{
	"Name",
	"Type",
	"Endpoint",
	"Status",
	"When",
}

type TestRunList struct {
	*tview.Table

	testRuns []models.TestRun

	renderScheduler RenderScheduler
}

func NewTestRunList(renderScheduler RenderScheduler) *TestRunList {
	list := &TestRunList{
		Table:           tview.NewTable(),
		renderScheduler: renderScheduler,
	}

	for i, header := range headers {
		header = strings.ToUpper(header)
		headerCell := tview.NewTableCell(header).SetStyle(styles.MetricNameStyle).SetExpansion(1).SetAlign(tview.AlignLeft)
		list.Table.SetCell(0, i, headerCell)
		list.Table.SetFixed(1, len(headers))
	}

	list.SetBorder(true).SetTitleColor(styles.HighlighColor).SetTitle("Test runs").SetBorderPadding(2, 0, 0, 0)
	list.SetSelectedStyle(styles.SelectedListItem)
	list.renderRuns()

	list.SetSelectable(true, false)
	list.Select(0, 0)
	list.SetSelectedFunc(func(row, column int) {
		fmt.Println(row, column)
	})

	return list
}

func (l *TestRunList) SetTestRuns(runs []models.TestRun) {
	l.testRuns = runs
	l.renderScheduler.Render(func() {
		l.renderRuns()
	})
}

func (l *TestRunList) renderRuns() {
	for i, run := range l.testRuns {
		l.Table.SetCell(i+1, 0, tview.NewTableCell(run.Name).SetStyle(styles.MetricValueStyle).SetExpansion(1).SetAlign(tview.AlignLeft))
		l.Table.SetCell(i+1, 1, tview.NewTableCell(run.Type).SetStyle(styles.MetricValueStyle).SetExpansion(1).SetAlign(tview.AlignLeft))
		l.Table.SetCell(i+1, 2, tview.NewTableCell(run.Endpoint).SetStyle(styles.MetricValueStyle).SetExpansion(1).SetAlign(tview.AlignLeft))
		l.Table.SetCell(i+1, 3, tview.NewTableCell(run.Status).SetStyle(styles.MetricValueStyle).SetExpansion(1).SetAlign(tview.AlignLeft))
		l.Table.SetCell(i+1, 4, tview.NewTableCell(run.When.String()).SetStyle(styles.MetricValueStyle).SetExpansion(1).SetAlign(tview.AlignLeft))
	}
}
