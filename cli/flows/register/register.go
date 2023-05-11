package flows_register

import (
	"github.com/kubeshop/tracetest/cli/flows"
	"github.com/kubeshop/tracetest/cli/global"
)

func Register(root global.Root) {
	flows.NewConfigure(root)

	server := flows.NewServer(root)
	flows.NewServerInstall(server)
}
