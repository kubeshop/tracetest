package misc_actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type Action[T any] interface {
	Run(ctx context.Context, args T) error
}

type actionAgs struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

type ActionArgsOption = func(args *actionAgs)

func ActionWithClient(client *openapi.APIClient) ActionArgsOption {
	return func(args *actionAgs) {
		args.client = client
	}
}

func ActionWithLogger(logger *zap.Logger) ActionArgsOption {
	return func(args *actionAgs) {
		args.logger = logger
	}
}

func ActionWithConfig(config config.Config) ActionArgsOption {
	return func(args *actionAgs) {
		args.config = config
	}
}

func NewActionArgs(options ...ActionArgsOption) actionAgs {
	args := actionAgs{}

	for _, option := range options {
		option(&args)
	}

	return args
}
