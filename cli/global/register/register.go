package global_register

import (
	deprecated_register "github.com/kubeshop/tracetest/cli/deprecated/register"
	flows_register "github.com/kubeshop/tracetest/cli/flows/register"
	"github.com/kubeshop/tracetest/cli/global"
	misc_register "github.com/kubeshop/tracetest/cli/misc/register"
	resources_register "github.com/kubeshop/tracetest/cli/resources/register"
)

func Register(root global.Root) {
	flows_register.Register(root)
	misc_register.Register(root)
	deprecated_register.Register(root)
	resources_register.Register(root)
}
