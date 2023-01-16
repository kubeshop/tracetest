package utils

import (
	"bytes"
	"os/exec"
)

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}
