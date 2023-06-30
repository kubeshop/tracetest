package environment

import "github.com/kubeshop/tracetest/server/resourcemanager"

const (
	ResourceName       = "Environment"
	ResourceNamePlural = "Environments"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationCreate,
	resourcemanager.OperationDelete,
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}
