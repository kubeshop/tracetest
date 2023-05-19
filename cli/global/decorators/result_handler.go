package global_decorators

import (
	"fmt"
	"os"

	global_types "github.com/kubeshop/tracetest/cli/global/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type resultHandler struct {
	Version
}

type ResultHandler interface {
	Version
	Error(err error)
	Response(res string)
}

var _ ResultHandler = &resultHandler{}

func WithResultHandler(command global_types.Command) global_types.Command {
	version, err := command.(Version)
	if !err {
		panic("command must implement Version interface")
	}

	resultHandler := &resultHandler{
		Version: version,
	}

	cmd := command.Get()
	cmd.PreRun = resultHandler.all(cmd.PreRun)
	cmd.Run = resultHandler.all(cmd.Run)
	cmd.PostRun = resultHandler.all(cmd.PostRun)

	resultHandler.Set(cmd)
	return resultHandler
}

func (d *resultHandler) all(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {

		next(cmd, args)
	}
}

func (d *resultHandler) Error(err error) {
	if err != nil {
		d.GetLogger().Error(fmt.Sprintf(`
Version
%s

An error ocurred when executing the command`, d.GetVersionText()), zap.Error(err))
		os.Exit(1)
	}
}

func (d *resultHandler) Response(res string) {
	if res != "" {
		fmt.Println(res)
		os.Exit(0)
	}
}
