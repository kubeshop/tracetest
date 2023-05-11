package resources_actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"go.uber.org/zap"
)

type ResourceActionsInterface interface {
	Logger() *zap.Logger
	Formatter() resources_formatters.ResourceFormatter
	FileType() yaml.FileType
	Name() string
	GetID(*file.File) (string, error)
	Apply(context.Context, file.File) (*file.File, error)
	List(context.Context, utils.ListArgs) (*file.File, error)
	Get(context.Context, string) (*file.File, error)
	Delete(context.Context, string) (string, error)
}

type ResourceActions struct {
	actions ResourceActionsInterface
}

func WrapActions(actions ResourceActionsInterface) ResourceActions {
	return ResourceActions{
		actions: actions,
	}
}

func (r *ResourceActions) Name() string {
	return r.actions.Name()
}

func (r *ResourceActions) Apply(ctx context.Context, args ApplyArgs) (*file.File, string, error) {
	if args.File == "" {
		return nil, "", fmt.Errorf("you must specify a file to be applied")
	}

	r.actions.Logger().Debug(
		fmt.Sprintf("applying %s", r.Name()),
		zap.String("file", args.File),
	)

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return nil, "", fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != r.actions.FileType() {
		return nil, "", fmt.Errorf(fmt.Sprintf(`file must be of type "%s"`, r.actions.FileType()))
	}

	f, err := r.actions.Apply(ctx, fileContent)
	if err != nil {
		return nil, "", err
	}

	Id, err := r.actions.GetID(f)
	if err != nil {
		return nil, "", err
	}

	result, err := f.WriteRaw()
	if err != nil {
		return nil, "", err
	}

	return &result, fmt.Sprintf("%s applied successfully (id: %s)", utils.Capitalize(r.actions.Name()), Id), nil
}

func (r *ResourceActions) List(ctx context.Context, args utils.ListArgs) (*file.File, error) {
	resource, err := r.actions.List(ctx, args)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (r *ResourceActions) Get(ctx context.Context, id string) (*file.File, error) {
	resource, err := r.actions.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return resource, err
}

func (r *ResourceActions) Export(ctx context.Context, id string, filePath string) error {
	file, err := r.actions.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = file.WriteRaw()
	return err
}

func (r *ResourceActions) Delete(ctx context.Context, id string) (string, error) {
	return r.actions.Delete(ctx, id)
}

func (r *ResourceActions) Formatter() resources_formatters.ResourceFormatter {
	return r.actions.Formatter()
}
