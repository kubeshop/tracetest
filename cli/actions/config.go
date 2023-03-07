package actions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kubeshop/tracetest/cli/file"
	"go.uber.org/zap"
)

type configActions struct {
	logger  *zap.Logger
	client  http.Client
	request http.Request
}

var _ ResourceActions = &configActions{}

func NewConfigActions(resourceArgs resourceArgs) configActions {
	return configActions{
		logger:  resourceArgs.logger,
		client:  resourceArgs.client,
		request: resourceArgs.request,
	}
}

func (config configActions) Apply(ctx context.Context, args ApplyArgs) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	config.logger.Debug(
		"applying analytics config",
		zap.String("file", args.File),
	)

	fileContent, err := file.Read(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "Config" {
		return fmt.Errorf(`file must be of type "Config"`)
	}

	body := bytes.NewBufferString(fileContent.Contents())

	config.request.Body = io.NopCloser(body)
	config.request.Method = http.MethodPost

	res, err := config.client.Do(&config.request)

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create config: %s", res.Status)
	}

	return err
}

func (config configActions) List(ctx context.Context) error {
	return nil
}

func (config configActions) Export(ctx context.Context, ID string) error {
	return nil
}

func (config configActions) Delete(ctx context.Context, ID string) error {
	return nil
}
