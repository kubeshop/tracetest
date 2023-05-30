package linter_resource

import "github.com/kubeshop/tracetest/server/resourcemanager"

const (
	ResourceName       = "Linter"
	ResourceNamePlural = "Linters"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}
