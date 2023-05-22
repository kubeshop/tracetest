package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	yamllib "github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/ui"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

type transactionActions struct {
	resourceArgs

	openapiClient *openapi.APIClient
}

var _ ResourceActions = &transactionActions{}

type RunArguments struct {
	EnvironmentID string
	Metadata      map[string]string
	Variables     map[string]string
	JUnit         string
	WaitForResult bool
}

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

	if newFile.HasID() {
		rawTransaction, err := actions.formatter.ToStruct(&newFile)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file into transaction: %w", err)
		}

		transaction := rawTransaction.(openapi.TransactionResource)
		return actions.resourceClient.Update(ctx, newFile, *transaction.Spec.Id)
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
	return actions.resourceClient.List(ctx, args)
}

// Name implements ResourceActions
func (actions transactionActions) Name() string {
	return "transaction"
}

func (actions transactionActions) Run(ctx context.Context, file file.File, args RunArguments) (any, error) {
	if args.JUnit != "" && !args.WaitForResult {
		return nil, fmt.Errorf("--junit option requires --wait-for-result")
	}

	appliedFile, err := actions.Apply(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("could not apply transaction: %w", err)
	}

	rawTransaction, err := actions.formatter.ToStruct(appliedFile)
	if err != nil {
		return nil, err
	}

	transaction := rawTransaction.(openapi.TransactionResource)

	return actions.RunByID(ctx, *transaction.Spec.Id, args)
}

func (actions transactionActions) RunByID(ctx context.Context, id string, args RunArguments) (any, error) {
	if id == "" {
		return nil, fmt.Errorf("id must be provided")
	}

	if args.JUnit != "" && !args.WaitForResult {
		return nil, fmt.Errorf("--junit option requires --wait-for-result")
	}

	return actions.runTransaction(ctx, id, args)
}

func (actions transactionActions) runTransaction(ctx context.Context, id string, args RunArguments) (any, error) {
	variables := make([]openapi.EnvironmentValue, 0, len(args.Variables))
	for varName, varValue := range args.Variables {
		variables = append(variables, openapi.EnvironmentValue{
			Key:   &varName,
			Value: &varValue,
		})
	}

	runInformation := openapi.RunInformation{
		EnvironmentId: &args.EnvironmentID,
		Metadata:      args.Metadata,
		Variables:     variables,
	}

	transactionRun, response, err := actions.openapiClient.ApiApi.RunTransaction(ctx, id).
		RunInformation(runInformation).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("could not send run request to server: %w", err)
	}

	if response != nil && response.StatusCode == http.StatusUnprocessableEntity {
		filledVariables, err := actions.askForMissingVariables(response)
		if err != nil {
			return nil, err
		}

		for name, value := range filledVariables {
			args.Variables[name] = value
		}

		return actions.runTransaction(ctx, id, args)
	}

	return transactionRun, nil
}

func (actions transactionActions) askForMissingVariables(resp *http.Response) (map[string]string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return map[string]string{}, fmt.Errorf("could not read response body: %w", err)
	}

	var missingVariablesError openapi.MissingVariablesError
	err = json.Unmarshal(body, &missingVariablesError)
	if err != nil {
		return map[string]string{}, fmt.Errorf("could not unmarshal response: %w", err)
	}

	uniqueMissingVariables := map[string]string{}
	for _, missingVariables := range missingVariablesError.MissingVariables {
		for _, variable := range missingVariables.Variables {
			defaultValue := ""
			if variable.DefaultValue != nil {
				defaultValue = *variable.DefaultValue
			}
			uniqueMissingVariables[*variable.Key] = defaultValue
		}
	}

	if len(uniqueMissingVariables) > 0 {
		ui.DefaultUI.Warning("Some variables are required by one or more tests")
		ui.DefaultUI.Info("Fill the values for each variable:")
	}

	filledVariables := map[string]string{}

	for variableName, variableDefaultValue := range uniqueMissingVariables {
		value := ui.DefaultUI.TextInput(variableName, variableDefaultValue)
		filledVariables[variableName] = value
	}

	return filledVariables, nil
}
