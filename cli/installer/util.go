package installer

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/kubeshop/tracetest/cli/installer/posix"
	"github.com/kubeshop/tracetest/cli/installer/win"
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
	if isWindows() {
		return win.CommandSuccess(probeCmd)
	}

	return posix.CommandSuccess(probeCmd)
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
