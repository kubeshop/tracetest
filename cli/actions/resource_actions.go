package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"go.uber.org/zap"
)

type ResourceActions interface {
	Logger() *zap.Logger
	FileType() yaml.FileType
	Name() string
	Apply(context.Context, file.File) error
	List(context.Context, utils.ListArgs) (string, error)
	Get(context.Context, string) (string, error)
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

func (r *resourceActions) Apply(ctx context.Context, args ApplyArgs) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	r.actions.Logger().Debug(
		fmt.Sprintf("applying %s", r.Name()),
		zap.String("file", args.File),
	)

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != r.actions.FileType() {
		return fmt.Errorf(fmt.Sprintf(`file must be of type "%s"`, r.actions.FileType()))
	}

	return r.actions.Apply(ctx, fileContent)
}

func (r *resourceActions) List(ctx context.Context, args utils.ListArgs) error {
	resources, err := r.actions.List(ctx, args)
	if err != nil {
		return err
	}

	fmt.Println(resources)
	return nil
}

func (r *resourceActions) Get(ctx context.Context, id string) error {
	resource, err := r.actions.Get(ctx, id)
	if err != nil {
		return err
	}

	fmt.Println(resource)
	return nil
}

func (r *resourceActions) Export(ctx context.Context, id string, filePath string) error {
	resource, err := r.actions.Get(ctx, id)
	if err != nil {
		return err
	}

	file, err := file.NewFromRaw(filePath, []byte(resource))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.WriteRaw()
	return err
}

func (r *resourceActions) Delete(ctx context.Context, id string) error {
	return r.actions.Delete(ctx, id)
}
