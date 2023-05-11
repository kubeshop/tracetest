package misc_register

import (
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/kubeshop/tracetest/cli/misc"
)

func Register(root global.Root) {
	test := misc.NewTest(root)

	misc.NewTestExport(test)
	misc.NewTestList(test)
	misc.NewTestRun(test)

	misc.NewDashboard(root)
	misc.NewDocGen(root)
	misc.NewVersion(root)
}
