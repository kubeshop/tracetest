package actions

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ApplyArgs struct {
	File string
}

type ListArgs struct {
	Take          int32
	Skip          int32
	SortDirection string
	SortBy        string
}

type ResourceActions interface {
	Apply(context.Context, ApplyArgs) error
	List(context.Context, ListArgs) error
	Get(context.Context, string) error
	Export(context.Context, string, string) error
	Delete(context.Context, string) error
}

type resourceArgs struct {
	logger *zap.Logger
	client *openapi.APIClient
	config config.Config
}

type ResourceArgsOption = func(args *resourceArgs)
type ResourceRegistry map[SupportedResources]ResourceActions
type SupportedResources string

var (
	SupportedResourceConfig         SupportedResources = "config"
	SupportedResourcePollingProfile SupportedResources = "pollingprofile"

	ErrResourceNotRegistered      = errors.New("resource not registered")
	ErrNotSupportedResourceAction = errors.New("the specified resource type doesn't support the action")
)

func NewResourceRegistry() ResourceRegistry {
	return map[SupportedResources]ResourceActions{}
}

func (r ResourceRegistry) Register(resource SupportedResources, actions ResourceActions) {
	r[resource] = actions
}

func (r ResourceRegistry) Get(resource SupportedResources) (ResourceActions, error) {
	resourceActions, found := r[resource]

	if !found {
		return nil, ErrResourceNotRegistered
	}

	return resourceActions, nil
}

func WithClient(client *openapi.APIClient) ResourceArgsOption {
	return func(args *resourceArgs) {
		args.client = client
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
