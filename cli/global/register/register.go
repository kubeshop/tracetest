package global_register

import (
	"github.com/kubeshop/tracetest/cli/global"
	resources_register "github.com/kubeshop/tracetest/cli/resources/register"
)

func Register(root global.Root) {
	// flows_register.Register(root)
	// misc_register.Register(root)
	// deprecated_register.Register(root)
	resources_register.Register(root)
}
