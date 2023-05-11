package main

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/global"
	global_register "github.com/kubeshop/tracetest/cli/global/register"
)

func Execute() {
	root := global.NewRoot()
	global_register.Register(root)

	if err := root.Cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
