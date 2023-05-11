package resources_actions

import (
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
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
	formatter      resources_formatters.ResourceFormatter
}

func (r resourceArgs) Logger() *zap.Logger {
	return r.logger
}

func (r resourceArgs) Formatter() resources_formatters.ResourceFormatter {
	return r.formatter
}

type ResourceArgsOption = func(args *resourceArgs)
type ResourceRegistry map[string]ResourceActions

var (
	ErrResourceNotRegistered      = errors.New("resource not registered")
	ErrNotSupportedResourceAction = errors.New("the specified resource type doesn't support the action")
)

func NewResourceRegistry() ResourceRegistry {
	return ResourceRegistry{}
}

func (r ResourceRegistry) Register(actions ResourceActionsInterface) {
	r[actions.Name()] = WrapActions(actions)
}

func (r ResourceRegistry) Get(name string) (ResourceActions, error) {
	actions, found := r[name]

	if !found {
		return ResourceActions{}, fmt.Errorf("resource not found: %s", name)
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

func WithFormatter(formatter resources_formatters.ResourceFormatter) ResourceArgsOption {
	return func(args *resourceArgs) {
		args.formatter = formatter
	}
}

func NewResourceArgs(options ...ResourceArgsOption) resourceArgs {
	args := resourceArgs{}

	for _, option := range options {
		option(&args)
	}

	return args
}
