package components

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/styles"
	"github.com/rivo/tview"
)

type HeaderData struct {
	Context AgentContext
	Metrics AgentMetrics
	Message BannerMessage
}

type AgentContext struct {
	OrganizationName       string
	EnvironmentName        string
	LastUsedTracingBackend string
}

type BannerMessage struct {
	Text string
	Type string
}

type AgentMetrics struct {
	Uptime   time.Duration
	TestRuns int64
	Traces   int64
	Spans    int64
}

type Header struct {
	*tview.Flex

	renderScheduler RenderScheduler
	sensor          sensors.Sensor
	data            HeaderData

	messageBanner   *MessageBanner
	uptimeTableCell *tview.TableCell
	runsTableCell   *tview.TableCell
	tracesTableCell *tview.TableCell
	spansTableCell  *tview.TableCell
}

func NewHeader(renderScheduler RenderScheduler, sensor sensors.Sensor) *Header {
	h := &Header{
		Flex:            tview.NewFlex(),
		renderScheduler: renderScheduler,
		sensor:          sensor,
		messageBanner:   NewMessageBanner(renderScheduler),
		uptimeTableCell: tview.NewTableCell("0s").SetStyle(styles.MetricValueStyle),
		runsTableCell:   tview.NewTableCell("0").SetStyle(styles.MetricValueStyle),
		tracesTableCell: tview.NewTableCell("0").SetStyle(styles.MetricValueStyle),
		spansTableCell:  tview.NewTableCell("0").SetStyle(styles.MetricValueStyle),
	}

	h.draw()

	return h
}

func (h *Header) draw() {
	h.Clear()

	flex := tview.NewFlex()

	flex.SetDirection(tview.FlexColumn).
		AddItem(h.getEnvironmentInformationTable(), 0, 1, true).
		AddItem(h.getMetricsTable(), 0, 1, true)

	h.Flex.SetDirection(tview.FlexRow).AddItem(h.messageBanner, 0, 0, true).AddItem(flex, 0, 8, true)

	h.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 's':
			h.messageBanner.SetMessage("Now you see me :D", events.Error)
			h.showMessageBanner()
		case 'w':
			h.messageBanner.SetMessage("This is a warning! :D", events.Warning)
			h.showMessageBanner()
		case 'h':
			h.messageBanner.SetMessage("", events.Error)
			h.hideMessageBanner()
		}

		return event
	})

	h.setupSensors()
}

func (h *Header) onDataChange() {
	h.renderScheduler.Render(func() {
		h.uptimeTableCell.SetText(h.data.Metrics.Uptime.String())
		h.runsTableCell.SetText(fmt.Sprintf("%d", h.data.Metrics.TestRuns))
		h.tracesTableCell.SetText(fmt.Sprintf("%d", h.data.Metrics.Traces))
		h.spansTableCell.SetText(fmt.Sprintf("%d", h.data.Metrics.Spans))
	})
}

func (h *Header) getEnvironmentInformationTable() tview.Primitive {
	table := tview.NewTable()
	table.SetBackgroundColor(styles.HeaderBackgroundColor)
	table.SetBorder(true).SetTitle("Environment").SetTitleColor(styles.HighlighColor)
	table.SetCell(0, 0, tview.NewTableCell("Organization: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(0, 1, tview.NewTableCell("my-company").SetStyle(styles.MetricValueStyle))
	table.SetCell(1, 0, tview.NewTableCell("Environment: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(1, 1, tview.NewTableCell("steve-dev").SetStyle(styles.MetricValueStyle))
	table.SetCell(2, 0, tview.NewTableCell("Last Tracing Backend: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(2, 1, tview.NewTableCell("Jaeger").SetStyle(styles.MetricValueStyle))
	table.SetCell(3, 0, tview.NewTableCell("Version: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(3, 1, tview.NewTableCell("v0.15.5").SetStyle(styles.MetricValueStyle))
	table.SetBorderPadding(1, 1, 2, 1)

	return table
}

func (h *Header) getMetricsTable() tview.Primitive {
	table := tview.NewTable()
	table.SetBackgroundColor(styles.HeaderBackgroundColor)
	table.SetBorder(true).SetTitle("Environment").SetTitleColor(styles.HighlighColor)
	table.SetCell(0, 0, tview.NewTableCell("Uptime: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(0, 1, h.uptimeTableCell)
	table.SetCell(1, 0, tview.NewTableCell("Runs: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(1, 1, h.runsTableCell)
	table.SetCell(2, 0, tview.NewTableCell("Traces: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(2, 1, h.tracesTableCell)
	table.SetCell(3, 0, tview.NewTableCell("Spans: ").SetStyle(styles.MetricNameStyle))
	table.SetCell(3, 1, h.spansTableCell)
	table.SetBorderPadding(1, 1, 2, 1)

	return table
}

func (h *Header) showMessageBanner() {
	h.Flex.ResizeItem(h.messageBanner, 0, 4)
}

func (h *Header) hideMessageBanner() {
	h.Flex.ResizeItem(h.messageBanner, 0, 0)
}

func (h *Header) setupSensors() {
	h.sensor.On(events.UptimeChanged, func(e sensors.Event) {
		var uptime time.Duration
		e.Unmarshal(&uptime)

		h.data.Metrics.Uptime = uptime
		h.onDataChange()
	})
}
