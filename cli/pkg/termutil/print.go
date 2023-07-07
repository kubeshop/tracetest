package termutil

import "github.com/pterm/pterm"

func GetGreenText(text string) string {
	return pterm.FgGreen.Sprintf(text)
}

func GetRedText(text string) string {
	return pterm.FgRed.Sprintf(text)
}
