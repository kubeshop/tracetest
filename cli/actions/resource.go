package actions

import (
	"errors"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/utils"
	"go.uber.org/zap"
)

type ApplyArgs struct {
	File string
}

type resourceArgs struct {
	logger         *zap.Logger
	resourceClient utils.ResourceClient
	config         config.Config
}

func (r resourceArgs) Logger() *zap.Logger {
	return r.logger
}

type ResourceArgsOption = func(args *resourceArgs)
type ResourceRegistry map[string]resourceActions

var (
	ErrResourceNotRegistered      = errors.New("resource not registered")
	ErrNotSupportedResourceAction = errors.New("the specified resource type doesn't support the action")
)

func NewResourceRegistry() ResourceRegistry {
	return ResourceRegistry{}
}

func (r ResourceRegistry) Register(actions ResourceActions) {
	r[actions.Name()] = WrapActions(actions)
}

func (r ResourceRegistry) Get(name string) (resourceActions, error) {
	actions, found := r[name]

	if !found {
		return resourceActions{}, ErrResourceNotRegistered
	}

	return actions, nil
}

func WithClient(client utils.ResourceClient) ResourceArgsOption {
	return func(args *resourceArgs) {
		args.resourceClient = client
	}
}

func WithLogger(logger *zap.Logger) ResourceArgsOption {
	return func(args *resourceArgs) {
		args.logger = logger
	}
}

func WithConfig(config config.Config) ResourceArgsOption {
	return func(args *resourceArgs) {
		args.config = config
	}
}

func NewResourceArgs(options ...ResourceArgsOption) resourceArgs {
	args := resourceArgs{}

	for _, option := range options {
		option(&args)
	}

	return args
}
