package installer

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const createIssueMsg = "If you need help, please create an issue: https://github.com/kubeshop/tracetest/issues/new/choose"

var DefaultUI UI = &ptermUI{}

type UI interface {
	Banner()

	Panic(error)
	Exit(string)

	Warning(string)
	Info(string)
	Success(string)

	Println(string)
	Title(string)

	Green(string) string
	Red(string) string

	Confirm(prompt string, defaultValue bool) bool
	Select(prompt string, options []option) (selected option)
	TextInput(msg, defaultValue string) string
}

type option struct {
	text string
	fn   func(ui UI)
}

type ptermUI struct{}

func (ui ptermUI) Banner() {
	pterm.Print("\n\n")

	pterm.DefaultBigText.
		WithLetters(putils.LettersFromString("TraceTest")).
		Render()

	pterm.Print("\n\n")

}

func (ui ptermUI) Panic(err error) {
	pterm.Error.WithFatal(true).Println(err)
}

func (ui ptermUI) Exit(msg string) {
	pterm.Error.Println(msg)
	os.Exit(1)
}

func (ui ptermUI) Warning(msg string) {
	pterm.Warning.Println(msg)
}

func (ui ptermUI) Info(msg string) {
	pterm.Info.Println(msg)
}

func (ui ptermUI) Success(msg string) {
	pterm.Success.Println(msg)
}

func (ui ptermUI) Println(msg string) {
	pterm.Println(msg)
}

func (ui ptermUI) Title(msg string) {
	pterm.Println(pterm.Yellow("\n-> ", msg, "\n"))
}

func (ui ptermUI) Green(msg string) string {
	return pterm.Green(msg)
}

func (ui ptermUI) Red(msg string) string {
	return pterm.Red(msg)
}

func (ui ptermUI) Confirm(msg string, defaultValue bool) bool {
	confirm, err := (&pterm.InteractiveConfirmPrinter{
		DefaultValue: defaultValue,
		DefaultText:  msg,
		TextStyle:    &pterm.ThemeDefault.PrimaryStyle,
		ConfirmText:  "Yes",
		ConfirmStyle: &pterm.ThemeDefault.SuccessMessageStyle,
		RejectText:   "No",
		RejectStyle:  &pterm.ThemeDefault.ErrorMessageStyle,
		SuffixStyle:  &pterm.ThemeDefault.SecondaryStyle,
	}).
		Show()
	if err != nil {
		ui.Panic(err)
	}

	return confirm
}

func (ui ptermUI) TextInput(msg, defaultValue string) string {
	text, err := (&pterm.InteractiveTextInputPrinter{
		TextStyle:   &pterm.ThemeDefault.PrimaryStyle,
		DefaultText: fmt.Sprintf("%s [%s]", msg, defaultValue),
		MultiLine:   false,
	}).
		Show()
	ui.Println("")
	if err != nil {
		ui.Panic(err)
	}

	if text == "" {
		return defaultValue
	}

	return text
}

func (ui ptermUI) Select(prompt string, options []option) (selected option) {
	textOpts := make([]string, len(options))
	lookupMap := make(map[string]int)

	for ix, opt := range options {
		textOpts[ix] = opt.text
		if _, ok := lookupMap[opt.text]; ok {
			panic(fmt.Sprintf("duplicated option %s", opt.text))
		}
		lookupMap[opt.text] = ix
	}

	selectedText, err := (&pterm.InteractiveSelectPrinter{
		TextStyle:     &pterm.ThemeDefault.PrimaryStyle,
		DefaultText:   prompt,
		Options:       textOpts,
		OptionStyle:   &pterm.ThemeDefault.DefaultText,
		DefaultOption: textOpts[0],
		MaxHeight:     5,
		Selector:      ">",
		SelectorStyle: &pterm.ThemeDefault.SecondaryStyle,
	}).
		Show()
	if err != nil {
		panic(err)
	}

	selectedIx := lookupMap[selectedText]
	return options[selectedIx]
}
