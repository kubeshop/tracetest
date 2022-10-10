package installer

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/kubeshop/tracetest/cli/ui"
)

type cmd struct {
	args               interface{}
	sudo               bool
	notConfirmMsg      string
	installDocs        string
	apt                string
	yum                string
	dnf                string
	homebrew           string
	macIntelChipManual string
	macAppleChipManual string
	windows            string
	other              string
}

func (c cmd) exec(ui ui.UI, args ...interface{}) interface{} {
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
		cmd = c.windows
	case other:
		cmd = c.other
	}

	if cmd == "" {
		ui.Exit(fmt.Sprintf(
			`We don't support your system for this action. Try again using another method.
			If you want to manually fix this issue, chech the install docs: %s`,
			c.installDocs,
		))
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

	execCmd(renderedCmd.String(), ui)

	return nil
}

func execCmd(cmd string, ui ui.UI) {
	if err := _execCmd(cmd); err != nil {
		ui.Panic(err)
	}
}

func execCmdIgnoreError(cmd string, ui ui.UI) {
	_execCmd(cmd)
}

func execCmdIgnoreErrors(cmd string) {
	_execCmd(cmd)
}

func getCmdOutputClean(cmd string) string {
	out := getCmdOutput(cmd)

	return strings.TrimSpace(out)
}

func commandExists(cmd string) bool {
	refreshEnvVariables()

	_, err := exec.LookPath(cmd)
	return err == nil
}
