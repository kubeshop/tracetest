package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"go.uber.org/zap"
)

type ResourceActions interface {
	Logger() *zap.Logger
	Formatter() formatters.ResourceFormatter
	FileType() yaml.FileType
	Name() string
	Apply(context.Context, file.File) (*file.File, error)
	List(context.Context, utils.ListArgs) (*file.File, error)
	Get(context.Context, string) (*file.File, error)
	Delete(context.Context, string) error
}

type resourceActions struct {
	actions ResourceActions
}

func WrapActions(actions ResourceActions) resourceActions {
	return resourceActions{
		actions: actions,
	}
}

func (r *resourceActions) Name() string {
	return r.actions.Name()
}

func (r *resourceActions) Apply(ctx context.Context, args ApplyArgs) (*file.File, error) {
	if args.File == "" {
		return nil, fmt.Errorf("you must specify a file to be applied")
	}

	r.actions.Logger().Debug(
		fmt.Sprintf("applying %s", r.Name()),
		zap.String("file", args.File),
	)

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != r.actions.FileType() {
		return nil, fmt.Errorf(fmt.Sprintf(`file must be of type "%s"`, r.actions.FileType()))
	}

	file, err := r.actions.Apply(ctx, fileContent)
	if err != nil {
		return nil, err
	}

	result, err := file.WriteRaw()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *resourceActions) List(ctx context.Context, args utils.ListArgs) (*file.File, error) {
	resource, err := r.actions.List(ctx, args)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *resourceActions) Get(ctx context.Context, id string) (*file.File, error) {
	resource, err := r.actions.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return resource, err
}

func (r *resourceActions) Export(ctx context.Context, id string, filePath string) error {
	file, err := r.actions.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = file.WriteRaw()
	return err
}

func (r *resourceActions) Delete(ctx context.Context, id string) error {
	return r.actions.Delete(ctx, id)
}

func (r *resourceActions) Formatter() formatters.ResourceFormatter {
	return r.actions.Formatter()
}
