package actions

import (
	"context"
)

type Action[T any] interface {
	Run(ctx context.Context, args T) error
}
