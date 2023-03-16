package resourcemanager

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

var ErrTypeNotSupported = errors.New("type not supported")

type Provisioner interface {
	Provision(_ context.Context, values map[string]any) error
}

func (m *manager[T]) Provision(ctx context.Context, values map[string]any) error {
	if values["type"] != m.resourceType {
		return ErrTypeNotSupported
	}
	targetResource := Resource[T]{}
	err := mapstructure.Decode(values, &targetResource)
	if err != nil {
		return fmt.Errorf(
			"cannot read provisioning for resource type %s: %w",
			m.resourceType, err,
		)
	}

	return m.rh.Provision(ctx, targetResource.Spec)
}
