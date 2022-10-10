//go:build !windows

package installer

import (
	"os"
	"os/exec"
	"syscall"
)

func _execCmd(cmd string) error {
	execCmd := exec.Command("/bin/sh", "-o", "xtrace", "-c", cmd)
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout

	err := execCmd.Run()

	execCmd.Stderr = nil
	execCmd.Stdin = nil
	execCmd.Stdout = nil

	return err
}

func getCmdOutput(cmd string) string {
	execCmd := exec.Command("/bin/sh", "-c", cmd)

	out, _ := execCmd.CombinedOutput()

	return string(out)
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

func refreshEnvVariables() {}
