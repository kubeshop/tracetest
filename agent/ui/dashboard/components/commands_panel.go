package components

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kubeshop/tracetest/agent/ui/dashboard/styles"
	"github.com/rivo/tview"
)

type CommandsPanel struct {
	*tview.Table
}

type Command struct {
	Name     string
	Shortcut string
}

func (c Command) GetCommand() string {
	cmd := c.Shortcut
	ctrl := "Ctrl"
	alt := "Alt"
	if runtime.GOOS == "darwin" {
		ctrl = "⌘"
		alt = "Opt"
	}

	cmd = strings.ReplaceAll(cmd, "Ctrl", ctrl)
	cmd = strings.ReplaceAll(cmd, "Alt", alt)

	return cmd
}

func NewCommandsPanel(commands []Command) *CommandsPanel {
	panel := &CommandsPanel{
		Table: tview.NewTable(),
	}
	panel.SetBorderPadding(0, 2, 2, 0)
	panel.SetBorder(true).SetTitle("Shortcuts")
	defaultPadding := "     "

	maxItemsPerColumn := 1

	for i, cmd := range commands {
		padding := defaultPadding
		row := i % maxItemsPerColumn
		column := int(i / maxItemsPerColumn)

		if column == 0 {
			padding = ""
		}

		panel.SetCell(row, column*2, tview.NewTableCell(fmt.Sprintf("%s%s:", padding, cmd.Name)).SetStyle(styles.MetricNameStyle).SetAlign(tview.AlignLeft))
		panel.SetCell(row, column*2+1, tview.NewTableCell(cmd.GetCommand()).SetStyle(styles.MetricValueStyle))
	}

	return panel
}
