package expression

import (
	"fmt"
	"sync"
)

type OperationRegistry struct {
	operations map[string]ExpressionOperation
	mutex      sync.Mutex
}

var operationRegistry *OperationRegistry = nil

func getOperationRegistry() *OperationRegistry {
	if operationRegistry != nil {
		return operationRegistry
	}

	registry := OperationRegistry{
		operations: map[string]ExpressionOperation{},
	}

	registry.Add("+", sum)
	registry.Add("-", subtract)
	registry.Add("*", multiply)
	registry.Add("/", divide)

	operationRegistry = &registry

	return operationRegistry
}

func (r *OperationRegistry) Add(name string, operation ExpressionOperation) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.operations[name] = operation
}

func (r *OperationRegistry) Get(name string) (ExpressionOperation, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	operation, found := r.operations[name]
	if !found {
		return nil, fmt.Errorf(`unsupported operation "%s"`, name)
	}

	return operation, nil
}
