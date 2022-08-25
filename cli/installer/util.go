package installer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"text/template"
)

func exitOption(msg string) func(ui UI) {
	return func(ui UI) {
		ui.Exit(msg)
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func commandSuccess(probeCmd string) bool {
	cmd := exec.Command("/bin/sh", "-c", probeCmd)
	err := cmd.Run()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			return ws.ExitStatus() == 0
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		return ws.ExitStatus() == 0
	}

	return false

}

type pkgManager int

const (
	undefined pkgManager = iota
	apt
	dnf
	yum
	homebrew
	macIntelChipManual
	macAppleChipManual
	windows
	other
)

var lastDetectedPkgManager pkgManager = undefined

func detectPkgManager() pkgManager {
	if lastDetectedPkgManager != undefined {
		return lastDetectedPkgManager
	}

	switch true {
	case commandExists("brew"):
		lastDetectedPkgManager = homebrew
	case commandExists("apt-get"):
		lastDetectedPkgManager = apt
	case commandExists("dnf"):
		lastDetectedPkgManager = dnf
	case commandExists("yum"):
		lastDetectedPkgManager = yum
	case runtime.GOOS == "darwin" && detectArchitecture() == "amd64":
		lastDetectedPkgManager = macIntelChipManual
	case runtime.GOOS == "darwin" && detectArchitecture() == "arm64":
		lastDetectedPkgManager = macAppleChipManual
	case runtime.GOOS == "windows":
		lastDetectedPkgManager = windows
	default:
		lastDetectedPkgManager = other
	}

	return lastDetectedPkgManager
}

func detectArchitecture() string {
	return runtime.GOARCH
}

type cmd struct {
	args               interface{}
	sudo               bool
	notConfirmMsg      string
	apt                string
	yum                string
	dnf                string
	homebrew           string
	macIntelChipManual string
	macAppleChipManual string
	windows            string
	other              string
}

func (c cmd) exec(ui UI, args ...interface{}) interface{} {
	cmd := ""
	switch detectPkgManager() {
	case apt:
		cmd = c.apt
	case yum:
		cmd = c.yum
	case dnf:
		cmd = c.dnf
	case homebrew:
		cmd = c.homebrew
	case macIntelChipManual:
		cmd = c.macIntelChipManual
	case macAppleChipManual:
		cmd = c.macAppleChipManual
	case windows:
		ui.Exit(
			fmt.Sprintf("We don't support windows yet =(. %s", c.windows),
		)
	case other:
		ui.Exit(
			fmt.Sprintf("OS not supported. %s", c.other),
		)
	}

	renderedCmd := &strings.Builder{}

	err := template.
		Must(template.New("cmd").Parse(cmd)).
		Execute(renderedCmd, c.args)

	if err != nil {
		ui.Panic(err)
	}

	ui.Warning("I'm going to run the following command:")
	ui.Println(renderedCmd.String())

	if c.sudo {
		ui.Warning("During the execution, you might be asked your `sudo` password.")
	}

	if !ui.Confirm("Do you want to execute?", false) {
		ui.Println(c.notConfirmMsg)
		os.Exit(1)
	}

	execCmd := exec.Command("/bin/sh", "-o", "xtrace", "-c", renderedCmd.String())
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout

	err = execCmd.Run()

	execCmd.Stderr = nil
	execCmd.Stdin = nil
	execCmd.Stdout = nil
	if err != nil {
		ui.Panic(err)
	}

	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
