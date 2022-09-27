package installer

import (
	"runtime"

	"github.com/kubeshop/tracetest/cli/installer/win"
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func installChocolatey(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries. You can try installing Chocolatey manually. See https://chocolatey.org/install#individual",
		installDocs:   "https://chocolatey.org/install#individual",
		windows:       "Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))",
	}).exec(ui)
}

func chocolateyForWindowsChecker(ui UI) {
	if !isWindows() {
		return
	}

	if commandExists("choco.exe") {
		ui.Println(ui.Green("✔ chocolatey already installed"))
		return
	}

	ui.Warning("I didn't find chocolatey in your system")
	option := ui.Select("What do you want to do?", []option{
		{"Install Chocolatey", installChocolatey},
		{"Fix manually", exitOption(
			"Check the chocolatey install docs on https://chocolatey.org/install#individual",
		)},
	}, 0)

	option.fn(ui)

	refreshEnvVariables()

	if commandExists("choco.exe") {
		ui.Println(ui.Green("✔ chocolatey was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ chocolatey could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func wslChecker(ui UI) {
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

// This function is important for installing things on Windows because programs are not
// available right after its installation. You have to refresh the environment variables
// to be able to find the command using the PATH env.
// Instead of closing and opening the CLI, we can execute this command instead.
func refreshEnvVariables() {
	win.ExecCmd("refreshenv", true)
}
