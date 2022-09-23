package installer

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
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
	case commandExists("apt-get"):
		lastDetectedPkgManager = apt
	case commandExists("dnf"):
		lastDetectedPkgManager = dnf
	case commandExists("yum"):
		lastDetectedPkgManager = yum
	case commandExists("brew"):
		lastDetectedPkgManager = homebrew
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

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
