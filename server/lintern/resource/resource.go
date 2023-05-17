package lintern_resource

import "github.com/kubeshop/tracetest/server/resourcemanager"

const (
	ResourceName       = "Lintern"
	ResourceNamePlural = "Linterns"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}
