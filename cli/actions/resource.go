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

type ResourceActions interface {
	Apply(ctx context.Context, args ApplyArgs) error
	List(ctx context.Context) error
	Get(ctx context.Context, ID string) error
	Export(ctx context.Context, ID string) error
	Delete(ctx context.Context, ID string) error
}

type resourceArgs struct {
	logger *zap.Logger
	client *openapi.APIClient
	config config.Config
}

type ResourceArgsOption = func(any)
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
	return func(args any) {
		typedArgs := args.(*resourceArgs)
		typedArgs.client = client
	}
}

func WithLogger(logger *zap.Logger) ResourceArgsOption {
	return func(args any) {
		typedArgs := args.(*resourceArgs)
		typedArgs.logger = logger
	}
}

func WithConfig(config config.Config) ResourceArgsOption {
	return func(args any) {
		typedArgs := args.(*resourceArgs)
		typedArgs.config = config
	}
}
