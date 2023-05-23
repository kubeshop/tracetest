package actions

import (
	"context"
	"fmt"
	"path/filepath"

	yamllib "github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

type transactionActions struct {
	resourceArgs

	openapiClient *openapi.APIClient
}

var _ ResourceActions = &transactionActions{}

func NewTransactionsActions(openapiClient *openapi.APIClient, options ...ResourceArgsOption) transactionActions {
	args := NewResourceArgs(options...)

	return transactionActions{
		resourceArgs:  args,
		openapiClient: openapiClient,
	}
}

// Apply implements ResourceActions
func (actions transactionActions) Apply(ctx context.Context, file file.File) (*file.File, error) {
	newFile, err := actions.convertTestPathToId(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("could not process tests from transaction: %w", err)
	}

	rawTransaction, err := actions.formatter.ToStruct(&newFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file into transaction: %w", err)
	}

	transaction := rawTransaction.(openapi.TransactionResource)

	if newFile.HasID() {
		_, err := actions.Get(ctx, transaction.Spec.GetId())
		if err != nil {
			// doesn't exist, so create it
			return actions.resourceClient.Create(ctx, newFile)
		}

		return actions.resourceClient.Update(ctx, newFile, transaction.Spec.GetId())
	}

	return actions.resourceClient.Create(ctx, newFile)
}

func (actions transactionActions) convertTestPathToId(ctx context.Context, f file.File) (file.File, error) {
	rawTransaction, err := actions.formatter.ToStruct(&f)
	if err != nil {
		return file.File{}, err
	}

	transaction := rawTransaction.(openapi.TransactionResource)
	for i, testPath := range transaction.Spec.Steps {
		if !utils.StringReferencesFile(testPath) {
			// not referencing a file, keep the value
			continue
		}

		id, err := actions.applyTestFile(ctx, f.AbsDir(), testPath)
		if err != nil {
			return file.File{}, fmt.Errorf(`could not apply test "%s": %w`, testPath, err)
		}

		transaction.Spec.Steps[i] = id
	}

	yamlContent, err := yamllib.Marshal(transaction)
	if err != nil {
		return file.File{}, fmt.Errorf("could not convert transaction to YAML: %w", err)
	}

	return file.NewFromRaw(f.Path(), yamlContent)
}

func (actions transactionActions) applyTestFile(ctx context.Context, workingDir string, filePath string) (string, error) {
	path := filepath.Join(workingDir, filePath)
	f, err := file.Read(path)
	if err != nil {
		return "", err
	}

	body, _, err := actions.openapiClient.ApiApi.
		UpsertDefinition(ctx).
		TextDefinition(openapi.TextDefinition{
			Content: openapi.PtrString(f.Contents()),
		}).
		Execute()

	if err != nil {
		return "", fmt.Errorf("could not upsert test definition: %w", err)
	}

	return body.GetId(), nil
}

// Delete implements ResourceActions
func (actions transactionActions) Delete(ctx context.Context, id string) (string, error) {
	return "Transaction successfully deleted", actions.resourceClient.Delete(ctx, id)
}

// FileType implements ResourceActions
func (actions transactionActions) FileType() yaml.FileType {
	return yaml.FileTypeTransaction
}

// Get implements ResourceActions
func (actions transactionActions) Get(ctx context.Context, id string) (*file.File, error) {
	return actions.resourceClient.Get(ctx, id)
}

// GetID implements ResourceActions
func (actions transactionActions) GetID(file *file.File) (string, error) {
	transaction, err := actions.formatter.ToStruct(file)
	if err != nil {
		return "", fmt.Errorf("could not convert file into struct: %w", err)
	}

	return *transaction.(openapi.TransactionResource).Spec.Id, nil
}

// List implements ResourceActions
func (actions transactionActions) List(ctx context.Context, args utils.ListArgs) (*file.File, error) {
	ctx = context.WithValue(ctx, "X-Tracetest-Augmented", true)
	return actions.resourceClient.List(ctx, args)
}

// Name implements ResourceActions
func (actions transactionActions) Name() string {
	return "transaction"
}
