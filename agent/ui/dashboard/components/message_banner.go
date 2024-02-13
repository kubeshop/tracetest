package components

import (
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/styles"
	"github.com/rivo/tview"
)

type MessageBanner struct {
	*tview.TextView

	renderScheduler RenderScheduler
}

func NewMessageBanner(renderScheduler RenderScheduler) *MessageBanner {
	banner := &MessageBanner{
		TextView:        tview.NewTextView(),
		renderScheduler: renderScheduler,
	}

	banner.TextView.SetBackgroundColor(styles.HeaderBackgroundColor)
	banner.TextView.SetMaxLines(5)
	banner.TextView.SetTextAlign(tview.AlignCenter).SetTextAlign(tview.AlignCenter)
	banner.TextView.SetWrap(true)
	banner.TextView.SetWordWrap(true)
	banner.SetBorderPadding(1, 0, 0, 0)
	banner.SetText("")

	return banner
}

func (b *MessageBanner) SetMessage(text string, messageType events.MessageType) {
	if messageType == events.Warning {
		b.SetBackgroundColor(styles.WarningMessageBackgroundColor)
		b.SetTextColor(styles.WarningMessageForegroundColor)
	}

	if messageType == events.Error {
		b.SetBackgroundColor(styles.ErrorMessageBackgroundColor)
		b.SetTextColor(styles.ErrorMessageForegroundColor)
	}
	b.TextView.SetText(text)
}
