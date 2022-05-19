package actions

import (
	"context"
)

type Action interface {
	Run(ctx context.Context, args []string) error
}
