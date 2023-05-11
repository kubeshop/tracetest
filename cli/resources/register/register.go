package resources_register

import (
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/kubeshop/tracetest/cli/resources"
)

func Register(root global.Root) {
	resources.NewApply(root)
	resources.NewDelete(root)
	resources.NewExport(root)
	resources.NewGet(root)
	resources.NewList(root)
}
