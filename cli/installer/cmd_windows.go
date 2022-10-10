//go:build windows

package installer

import (
	"os"
	"os/exec"
	"syscall"
)

func _execCmd(cmd string) error {
	execCmd := exec.Command("powershell", "-nologo", "-noprofile", cmd)
	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout

	err := execCmd.Run()

	execCmd.Stderr = nil
	execCmd.Stdin = nil
	execCmd.Stdout = nil

	return err
}

func getCmdOutput(cmd string) string {
	execCmd := exec.Command("powershell", "-nologo", "-noprofile", cmd)

	out, _ := execCmd.CombinedOutput()

	return string(out)
}

func commandSuccess(probeCmd string) bool {
	cmd := exec.Command("powershell", "-nologo", "-noprofile", probeCmd)
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

// This function is important for installing things on Windows because programs are not
// available right after its installation. You have to refresh the environment variables
// to be able to find the command using the PATH env.
// Instead of closing and opening the CLI, we can execute this command instead.
func refreshEnvVariables() {
	_execCmd("$foo = refreshenv 2>$null")
}
