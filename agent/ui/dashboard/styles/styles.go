package styles

import "github.com/gdamore/tcell/v2"

var (
	HeaderBackgroundColor = tcell.NewRGBColor(18, 18, 18)
	HeaderLogoColor       = tcell.NewRGBColor(253, 166, 34)

	ErrorMessageBackgroundColor   = tcell.NewRGBColor(102, 0, 0)
	ErrorMessageForegroundColor   = tcell.NewRGBColor(255, 255, 255)
	WarningMessageBackgroundColor = tcell.NewRGBColor(227, 149, 30)
	WarningMessageForegroundColor = tcell.NewRGBColor(0, 0, 0)

	TableSelectionColor = tcell.NewRGBColor(0, 0, 255)

	HighlighColor = tcell.NewRGBColor(253, 166, 34)

	MetricNameStyle = tcell.Style{}.
			Foreground(HighlighColor).
			Bold(true)

	MetricValueStyle = tcell.Style{}.
				Foreground(tcell.NewRGBColor(255, 255, 255)).
				Bold(true)

	SelectedListItem = tcell.Style{}.
				Foreground(tcell.NewRGBColor(255, 255, 255)).
				Background(tcell.NewRGBColor(114, 159, 207)).
				Bold(true)
)
