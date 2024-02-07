package components

import (
	"github.com/kubeshop/tracetest/agent/ui/dashboard/styles"
	"github.com/rivo/tview"
)

type TestRunList struct {
	*tview.Table
}

func NewTestRunList() *TestRunList {
	list := &TestRunList{
		Table: tview.NewTable(),
	}

	list.SetBorder(true).SetTitleColor(styles.HighlighColor).SetTitle("Test runs")

	return list
}
