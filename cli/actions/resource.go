package actions

import (
	"context"
	"fmt"
	"net/http"

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
	logger  *zap.Logger
	client  http.Client
	request http.Request
}

type supportedResources string

var (
	Config supportedResources = "config"
)

func NewResourceActions(resourceType string, logger *zap.Logger, client http.Client, request http.Request) (ResourceActions, error) {
	args := resourceArgs{
		logger:  logger,
		client:  client,
		request: request,
	}

	switch resourceType {
	case string(Config):
		return NewConfigActions(args), nil
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}
