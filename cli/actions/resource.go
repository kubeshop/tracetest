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
	Export(ctx context.Context, ID string) error
	Delete(ctx context.Context, ID string) error
}

type resourceArgs struct {
	logger    *zap.Logger
	client    *openapi.APIClient
	cliConfig config.Config
}
type resourceArgsOption = func(*resourceArgs)

type ResourceRegistry map[string]ResourceActions
type supportedResources string

var (
	SupportedResourceConfig supportedResources = "config"

	ErrResourceNotRegistered = errors.New("resource not registered")
)

func NewResourceRegistry() ResourceRegistry {
	return map[string]ResourceActions{}
}

func (r ResourceRegistry) Register(resource string, actions ResourceActions) {
	r[resource] = actions
}

func (r ResourceRegistry) Get(resource string) (ResourceActions, error) {
	resourceActions, found := r[resource]

	if !found {
		return nil, ErrResourceNotRegistered
	}

	return resourceActions, nil
}

func WithClient(client *openapi.APIClient) resourceArgsOption {
	return func(args *resourceArgs) {
		args.client = client
	}
}

func WithLogger(logger *zap.Logger) resourceArgsOption {
	return func(args *resourceArgs) {
		args.logger = logger
	}
}
