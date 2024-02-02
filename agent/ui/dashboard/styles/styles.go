package styles

import "github.com/gdamore/tcell/v2"

var (
	HeaderBackgroundColor = tcell.NewRGBColor(18, 18, 18)
	HeaderLogoColor       = tcell.NewRGBColor(253, 166, 34)

	MetricNameStyle = tcell.Style{}.
			Foreground(tcell.NewRGBColor(253, 166, 34)).
			Bold(true)

	MetricValueStyle = tcell.Style{}.
				Foreground(tcell.NewRGBColor(255, 255, 255)).
				Bold(true)
)
