package environment

import "github.com/kubeshop/tracetest/server/resourcemanager"

const (
	ResourceName       = "Environment"
	ResourceNamePlural = "Environments"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationCreate,
	resourcemanager.OperationUpsert,
	resourcemanager.OperationDelete,
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}
