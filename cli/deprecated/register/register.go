package deprecated_register

import (
	"github.com/kubeshop/tracetest/cli/deprecated"
	"github.com/kubeshop/tracetest/cli/global"
)

func Register(root global.Root) {
	deprecated.NewDatastoreLegacy(root)
	deprecated.NewEnvironmentLegacy(root)
}
