package resourcemanager

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
	"golang.org/x/exp/slices"
)

type Operation string

const (
	OperationList   Operation = "list"
	OperationCreate Operation = "create"
	OperationUpdate Operation = "update"
	OperationGet    Operation = "get"
	OperationDelete Operation = "delete"
)

var availableOperations = []Operation{
	OperationList,
	OperationCreate,
	OperationUpdate,
	OperationGet,
	OperationDelete,
}

type Operationer interface {
	Operations() []Operation
}

type SortableHandler interface {
	SortingFields() []string
}

type List[T ResourceSpec] interface {
	SortableHandler
	List(_ context.Context, take, skip int, query, sortBy, sortDirection string) ([]T, error)
	Count(_ context.Context, query string) (int, error)
}

type Create[T ResourceSpec] interface {
	Create(context.Context, T) (T, error)
	SetID(T, id.ID) T
}

type Update[T ResourceSpec] interface {
	Update(context.Context, T) (T, error)
}

type Get[T ResourceSpec] interface {
	Get(context.Context, id.ID) (T, error)
}

type Delete[T ResourceSpec] interface {
	Delete(context.Context, id.ID) error
}

type resourceHandler[T ResourceSpec] struct {
	SetID         func(T, id.ID) T
	List          func(_ context.Context, take, skip int, query, sortBy, sortDirection string) ([]T, error)
	Count         func(_ context.Context, query string) (int, error)
	SortingFields func() []string
	Create        func(context.Context, T) (T, error)
	Update        func(context.Context, T) (T, error)
	Get           func(context.Context, id.ID) (T, error)
	Delete        func(context.Context, id.ID) error
}

func (rh *resourceHandler[T]) bindOperations(enabledOperations []Operation, handler any) error {
	if len(enabledOperations) < 1 {
		return fmt.Errorf("no operations enabled")
	}

	if slices.Contains(enabledOperations, OperationList) {
		err := rh.bindListOperation(handler)
		if err != nil {
			return err
		}
	}

	if slices.Contains(enabledOperations, OperationCreate) {
		err := rh.bindCreateOperation(handler)
		if err != nil {
			return err
		}
	}

	if slices.Contains(enabledOperations, OperationUpdate) {
		err := rh.bindUpdateOperation(handler)
		if err != nil {
			return err
		}
	}

	if slices.Contains(enabledOperations, OperationGet) {
		err := rh.bindGetOperation(handler)
		if err != nil {
			return err
		}
	}

	if slices.Contains(enabledOperations, OperationDelete) {
		err := rh.bindDeleteOperation(handler)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rh *resourceHandler[T]) bindListOperation(handler any) error {
	casted, ok := handler.(List[T])
	if !ok {
		return fmt.Errorf("handler does not implement interface `List[T]`")
	}
	rh.List = casted.List
	rh.Count = casted.Count
	rh.SortingFields = casted.SortingFields

	return nil
}

func (rh *resourceHandler[T]) bindCreateOperation(handler any) error {
	casted, ok := handler.(Create[T])
	if !ok {
		return fmt.Errorf("handler does not implement interface `Create[T]`")
	}

	rh.Create = casted.Create
	rh.SetID = casted.SetID

	return nil
}

func (rh *resourceHandler[T]) bindUpdateOperation(handler any) error {
	casted, ok := handler.(Update[T])
	if !ok {
		return fmt.Errorf("handler does not implement interface `Update[T]`")
	}

	rh.Update = casted.Update

	return nil
}

func (rh *resourceHandler[T]) bindGetOperation(handler any) error {
	casted, ok := handler.(Get[T])
	if !ok {
		return fmt.Errorf("handler does not implement interface `Get[T]`")
	}
	rh.Get = casted.Get

	return nil
}

func (rh *resourceHandler[T]) bindDeleteOperation(handler any) error {
	casted, ok := handler.(Delete[T])
	if !ok {
		return fmt.Errorf("handler does not implement interface `Delete[T]`")
	}
	rh.Delete = casted.Delete

	return nil
}
