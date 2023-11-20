package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

var DefaultUI UI = &ptermUI{}

type UI interface {
	Banner(version string)

	Panic(error)
	Exit(string)

	Errorf(string, ...any)
	Warningf(string, ...any)
	Infof(string, ...any)
	Printlnf(string, ...any)

	Error(string)
	Warning(string)
	Info(string)
	Success(string)
	Finish()

	Println(string)
	Title(string)
	OpenBrowser(string) error

	Green(string) string
	Red(string) string

	Confirm(prompt string, defaultValue bool) bool
	Enter(msg string) bool
	Select(prompt string, options []Option, defaultIndex int) (selected Option)
	TextInput(msg, defaultValue string) string
}

type Option struct {
	Text string
	Fn   func(ui UI)
}

type ptermUI struct{}

func (ui ptermUI) Banner(version string) {
	pterm.Print("\n\n")

	pterm.DefaultBigText.
		WithLetters(putils.LettersFromString("TraceTest")).
		Render()

	pterm.Print(fmt.Sprintf("Version: %s", version))

	pterm.Print("\n\n")

}

func (ui ptermUI) Panic(err error) {
	pterm.Error.WithFatal(true).Println(err)
}

func (ui ptermUI) Finish() {
	os.Exit(0)
}

func (ui ptermUI) Exit(msg string) {
	pterm.Error.Println(msg)
	os.Exit(1)
}

func (ui ptermUI) Errorf(msg string, args ...any) {
	ui.Error(fmt.Sprintf(msg, args...))
}

func (ui ptermUI) Warningf(msg string, args ...any) {
	ui.Warning(fmt.Sprintf(msg, args...))
}

func (ui ptermUI) Infof(msg string, args ...any) {
	ui.Info(fmt.Sprintf(msg, args...))
}

func (ui ptermUI) Printlnf(msg string, args ...any) {
	ui.Println(fmt.Sprintf(msg, args...))
}

func (ui ptermUI) Error(msg string) {
	pterm.Error.Println(msg)
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
		TextStyle:    &pterm.ThemeDefault.DefaultText,
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

func (ui ptermUI) Enter(msg string) bool {
	confirm, err := (&pterm.InteractiveConfirmPrinter{
		DefaultText:  msg,
		DefaultValue: true,
		TextStyle:    &pterm.ThemeDefault.DefaultText,
		ConfirmText:  "Enter",
		RejectText:   "Cancel",
		RejectStyle:  &pterm.ThemeDefault.ErrorMessageStyle,
		ConfirmStyle: &pterm.ThemeDefault.SuccessMessageStyle,
		SuffixStyle:  &pterm.ThemeDefault.SecondaryStyle,
	}).
		Show()
	if err != nil {
		ui.Panic(err)
	}

	return confirm
}

func (ui ptermUI) TextInput(msg, defaultValue string) string {
	text := msg
	if defaultValue != "" {
		text = fmt.Sprintf("%s (default: %s)", msg, defaultValue)
	}
	text, err := (&pterm.InteractiveTextInputPrinter{
		TextStyle:   &pterm.ThemeDefault.DefaultText,
		DefaultText: text,
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

func (ui ptermUI) Select(prompt string, options []Option, defaultIndex int) (selected Option) {
	textOpts := make([]string, len(options))
	lookupMap := make(map[string]int)

	for ix, opt := range options {
		textOpts[ix] = opt.Text
		if _, ok := lookupMap[opt.Text]; ok {
			panic(fmt.Sprintf("duplicated option %s", opt.Text))
		}
		lookupMap[opt.Text] = ix
	}

	selectedText, err := (&pterm.InteractiveSelectPrinter{
		TextStyle:     &pterm.ThemeDefault.DefaultText,
		DefaultText:   prompt,
		Options:       textOpts,
		OptionStyle:   &pterm.ThemeDefault.DefaultText,
		DefaultOption: textOpts[defaultIndex],
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

func (ui ptermUI) OpenBrowser(u string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", u).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start()
	case "darwin":
		return exec.Command("open", u).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
