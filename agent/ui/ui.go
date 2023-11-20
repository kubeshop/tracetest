package ui

var DefaultUI ConsoleUI = &ptermUI{}

type ConsoleUI interface {
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
	Fn   func(ui ConsoleUI)
}
