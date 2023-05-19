package global_types

import "github.com/spf13/cobra"

type command struct {
	cmd *cobra.Command
}

type Command interface {
	Set(*cobra.Command)
	Get() *cobra.Command
}

var _ Command = &command{}

func NewCommand(cmd *cobra.Command) Command {
	return &command{
		cmd: cmd,
	}
}

func (d *command) Set(cmd *cobra.Command) {
	d.cmd = cmd
}

func (d *command) Get() *cobra.Command {
	return d.cmd
}
