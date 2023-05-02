package resourcemanager

import (
	"context"
	"errors"
	"fmt"
)

var ErrTypeNotSupported = errors.New("type not supported")

type Provisioner interface {
	Provision(_ context.Context, values map[string]any) error
}

func (m *manager[T]) Provision(ctx context.Context, values map[string]any) error {
	if values["type"] != m.resourceTypeSingular {
		return ErrTypeNotSupported
	}

	targetResource := Resource[T]{}
	err := decode(values, &targetResource)
	if err != nil {
		return fmt.Errorf(
			"cannot read provisioning for resource type %s: %w",
			m.resourceTypeSingular, err,
		)
	}

	if !targetResource.Spec.HasID() {
		targetResource.Spec = m.rh.SetID(
			targetResource.Spec,
			m.config.idgen(),
		)
	}

	return m.rh.Provision(ctx, targetResource.Spec)
}
