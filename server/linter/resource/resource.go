package linter_resource

import "github.com/kubeshop/tracetest/server/resourcemanager"

const (
	ResourceName       = "linter"
	ResourceNamePlural = "linters"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}
