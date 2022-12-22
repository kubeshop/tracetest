package installer

import (
	"runtime"

	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func chocolateyForWindowsChecker(ui cliUI.UI) {
	if !isWindows() {
		return
	}

	if commandExists("choco.exe") {
		ui.Println(ui.Green("✔ chocolatey already installed"))
		return
	}

	ui.Warning("I didn't find chocolatey in your system")
	exitOption(
		"Check the chocolatey install docs on https://chocolatey.org/install#individual",
	)(ui)
}

func wslChecker(ui cliUI.UI) {
	if !isWindows() {
		return
	}

	if commandExists("wsl.exe") {
		ui.Println(ui.Green("✔ WSL module already installed"))
		return
	}

	ui.Warning("I didn't find WSL installed in your system")
	ui.Exit("WSL is a requirement for running Tracetest on Windows. Install it before proceeding: https://learn.microsoft.com/en-us/windows/wsl/install")
}
