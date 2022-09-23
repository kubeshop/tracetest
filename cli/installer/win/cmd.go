package win

import (
	"os"
	"os/exec"
)

func ExecCmd(cmd string) error {
	execCmd := exec.Command("powershell", "-nologo", "-noprofile", cmd)
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout

	err := execCmd.Run()

	execCmd.Stderr = nil
	execCmd.Stdin = nil
	execCmd.Stdout = nil

	return err
}
