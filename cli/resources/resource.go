package resources

import global_decorators "github.com/kubeshop/tracetest/cli/global/decorators"

type Resource interface {
	global_decorators.Logger
	global_decorators.Analytics
	global_decorators.Version
	global_decorators.Config
	global_decorators.ResultHandler
}
