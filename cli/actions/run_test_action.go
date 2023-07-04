package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cienvironment "github.com/cucumber/ci-environment/go"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/ui"
	"github.com/kubeshop/tracetest/cli/variable"
	"go.uber.org/zap"
)

type RunResourceArgs struct {
	DefinitionFile string
	EnvID          string
	WaitForResult  bool
	JUnit          string
}

func (args RunResourceArgs) Validate() error {
	if args.DefinitionFile == "" {
		return fmt.Errorf("you must specify a definition file to run a test")
	}

	if args.JUnit != "" && !args.WaitForResult {
		return fmt.Errorf("--junit option requires --wait-for-result")
	}

	return nil
}

type runTestAction struct {
	config        config.Config
	logger        *zap.Logger
	openapiClient *openapi.APIClient
	environments  resourcemanager.Client
	tests         resourcemanager.Client
	transactions  resourcemanager.Client
	yamlFormat    resourcemanager.Format
	cliExit       func(int)
}

func NewRunTestAction(
	config config.Config,
	logger *zap.Logger,
	openapiClient *openapi.APIClient,
	tests resourcemanager.Client,
	transactions resourcemanager.Client,
	environments resourcemanager.Client,
	cliExit func(int),
) runTestAction {

	yamlFormat, err := resourcemanager.Formats.Get(resourcemanager.FormatYAML)
	if err != nil {
		panic(fmt.Errorf("could not get yaml format: %w", err))
	}

	return runTestAction{
		config,
		logger,
		openapiClient,
		environments,
		tests,
		transactions,
		yamlFormat,
		cliExit,
	}
}

func (a runTestAction) Run(ctx context.Context, args RunResourceArgs) error {
	if err := args.Validate(); err != nil {
		return err
	}

	a.logger.Debug(
		"Running test from definition",
		zap.String("definitionFile", args.DefinitionFile),
		zap.String("environment", args.EnvID),
		zap.Bool("waitForResults", args.WaitForResult),
		zap.String("junit", args.JUnit),
	)

	f, err := fileutil.Read(args.DefinitionFile)
	if err != nil {
		return fmt.Errorf("cannot read definition file %s: %w", args.DefinitionFile, err)
	}
	df := defFile{f}
	a.logger.Debug("Definition file read", zap.String("absolutePath", df.AbsPath()))

	envID, err := a.resolveEnvID(ctx, args.EnvID)
	if err != nil {
		return fmt.Errorf("cannot resolve environment id: %w", err)
	}
	a.logger.Debug("env resolved", zap.String("ID", envID))

	df, err = a.apply(ctx, df)
	if err != nil {
		return fmt.Errorf("cannot apply definition file: %w", err)
	}

	a.logger.Debug("Definition file applied", zap.String("updated", string(df.Contents())))

	var run *openapi.ExecuteDefinitionResponse
	var envVars []envVar

	// iterate until we have all env vars,
	// or the server returns an actual error
	for {
		run, err = a.run(ctx, df, envID, envVars)
		if err == nil {
			break
		}
		missingEnvVarsErr, ok := err.(missingEnvVarsError)
		if !ok {
			// actual error, return
			return fmt.Errorf("cannot run test: %w", err)
		}

		// missing vars error
		envVars = a.askForMissingVars([]envVar(missingEnvVarsErr))
	}

	if !args.WaitForResult {
		// TODO: show url?
		return nil
	}

	result, err := a.waitForResult(ctx, run)
	if err != nil {
		return fmt.Errorf("cannot wait for test result: %w", err)
	}

	if args.JUnit != "" {
		err := a.writeJUnitReport(ctx, result, args.JUnit)
		if err != nil {
			return fmt.Errorf("cannot write junit report: %w", err)
		}
	}

	return nil
}

func (a runTestAction) resolveEnvID(ctx context.Context, envID string) (string, error) {
	if !fileutil.IsFilePath(envID) {
		a.logger.Debug("envID is not a file path", zap.String("envID", envID))
		return envID, nil
	}

	f, err := fileutil.Read(envID)
	if err != nil {
		return "", fmt.Errorf("cannot read environment file %s: %w", envID, err)
	}

	a.logger.Debug("envID is a file path", zap.String("filePath", envID), zap.Any("file", f))
	updatedEnv, err := a.environments.Apply(ctx, f, a.yamlFormat)
	if err != nil {
		return "", fmt.Errorf("could not read environment file: %w", err)
	}

	var env openapi.EnvironmentResource
	err = yaml.Unmarshal([]byte(updatedEnv), &env)
	if err != nil {
		a.logger.Error("error parsing json", zap.String("content", updatedEnv), zap.Error(err))
		return "", fmt.Errorf("could not unmarshal environment json: %w", err)
	}

	return env.Spec.GetId(), nil
}

func (a runTestAction) injectLocalEnvVars(ctx context.Context, df defFile) (defFile, error) {
	variableInjector := variable.NewInjector(variable.WithVariableProvider(
		variable.EnvironmentVariableProvider{},
	))

	injected, err := variableInjector.ReplaceInString(string(df.Contents()))
	if err != nil {
		return df, fmt.Errorf("cannot inject local environment variables: %w", err)
	}

	df = defFile{fileutil.New(df.AbsPath(), []byte(injected))}

	return df, nil
}

func (a runTestAction) apply(ctx context.Context, df defFile) (defFile, error) {
	defType, err := getTypeFromFile(df)
	if err != nil {
		return df, fmt.Errorf("cannot get type from definition file: %w", err)
	}
	a.logger.Debug("definition file type", zap.String("type", defType))

	switch defType {
	case "Test":
		return a.applyTest(ctx, df)
	case "Transaction":
		return a.applyTransaction(ctx, df)
	default:
		return df, fmt.Errorf("unknown type %s", defType)
	}
}

func (a runTestAction) applyTest(ctx context.Context, df defFile) (defFile, error) {
	df, err := a.injectLocalEnvVars(ctx, df)
	if err != nil {
		return df, fmt.Errorf("cannot inject local env vars: %w", err)
	}

	var test openapi.TestResource
	err = yaml.Unmarshal(df.Contents(), &test)
	if err != nil {
		a.logger.Error("error parsing test", zap.String("content", string(df.Contents())), zap.Error(err))
		return df, fmt.Errorf("could not unmarshal test yaml: %w", err)
	}

	test, err = consolidateGRPCFile(df, test)
	if err != nil {
		return df, fmt.Errorf("could not consolidate grpc file: %w", err)
	}

	marshalled, err := yaml.Marshal(test)
	if err != nil {
		return df, fmt.Errorf("could not marshal test yaml: %w", err)
	}
	df = defFile{fileutil.New(df.AbsPath(), marshalled)}

	a.logger.Debug("applying test",
		zap.String("absolutePath", df.AbsPath()),
		zap.String("id", test.Spec.GetId()),
		zap.String("marshalled", string(marshalled)),
	)

	updated, err := a.tests.Apply(ctx, df.File, a.yamlFormat)
	if err != nil {
		return df, fmt.Errorf("could not read test file: %w", err)
	}

	df = defFile{fileutil.New(df.AbsPath(), []byte(updated))}

	err = yaml.Unmarshal(df.Contents(), &test)
	if err != nil {
		a.logger.Error("error parsing updated test", zap.String("content", string(df.Contents())), zap.Error(err))
		return df, fmt.Errorf("could not unmarshal test yaml: %w", err)
	}

	a.logger.Debug("test applied",
		zap.String("absolutePath", df.AbsPath()),
		zap.String("id", test.Spec.GetId()),
	)

	return df, nil
}

func consolidateGRPCFile(df defFile, test openapi.TestResource) (openapi.TestResource, error) {
	if test.Spec.ServiceUnderTest.GetTriggerType() != "grpc" {
		return test, nil
	}

	pbFilePath := df.RelativeFile(test.Spec.ServiceUnderTest.Grpc.GetProtobufFile())

	pbFile, err := fileutil.Read(pbFilePath)
	if err != nil {
		return test, fmt.Errorf(`cannot read protobuf file: %w`, err)
	}

	test.Spec.Trigger.Grpc.SetProtobufFile(string(pbFile.Contents()))

	return test, nil
}

func (a runTestAction) applyTransaction(ctx context.Context, df defFile) (defFile, error) {
	df, err := a.injectLocalEnvVars(ctx, df)
	if err != nil {
		return df, fmt.Errorf("cannot inject local env vars: %w", err)
	}

	var tran openapi.TransactionResource
	err = yaml.Unmarshal(df.Contents(), &tran)
	if err != nil {
		a.logger.Error("error parsing transaction", zap.String("content", string(df.Contents())), zap.Error(err))
		return df, fmt.Errorf("could not unmarshal transaction yaml: %w", err)
	}

	tran, err = a.mapTransactionSteps(ctx, df, tran)
	if err != nil {
		return df, fmt.Errorf("could not map transaction steps: %w", err)
	}

	marshalled, err := yaml.Marshal(tran)
	if err != nil {
		return df, fmt.Errorf("could not marshal test yaml: %w", err)
	}
	df = defFile{fileutil.New(df.AbsPath(), marshalled)}

	a.logger.Debug("applying transaction",
		zap.String("absolutePath", df.AbsPath()),
		zap.String("id", tran.Spec.GetId()),
		zap.String("marshalled", string(marshalled)),
	)

	updated, err := a.transactions.Apply(ctx, df.File, a.yamlFormat)
	if err != nil {
		return df, fmt.Errorf("could not read transaction file: %w", err)
	}

	df = defFile{fileutil.New(df.AbsPath(), []byte(updated))}

	err = yaml.Unmarshal(df.Contents(), &tran)
	if err != nil {
		a.logger.Error("error parsing updated transaction", zap.String("content", updated), zap.Error(err))
		return df, fmt.Errorf("could not unmarshal transaction yaml: %w", err)
	}

	a.logger.Debug("transaction applied",
		zap.String("absolutePath", df.AbsPath()),
		zap.String("updated id", tran.Spec.GetId()),
	)

	return df, nil
}

func (a runTestAction) mapTransactionSteps(ctx context.Context, df defFile, tran openapi.TransactionResource) (openapi.TransactionResource, error) {
	for i, step := range tran.Spec.GetSteps() {
		a.logger.Debug("mapping transaction step",
			zap.Int("index", i),
			zap.String("step", step),
		)
		if !fileutil.LooksLikeFilePath(step) {
			a.logger.Debug("does not look like a file path",
				zap.Int("index", i),
				zap.String("step", step),
			)
			continue
		}

		f, err := fileutil.Read(df.RelativeFile(step))
		if err != nil {
			return openapi.TransactionResource{}, fmt.Errorf("cannot read test file: %w", err)
		}

		testDF, err := a.applyTest(ctx, defFile{f})
		if err != nil {
			return openapi.TransactionResource{}, fmt.Errorf("cannot apply test '%s': %w", step, err)
		}

		var test openapi.TestResource
		err = yaml.Unmarshal(testDF.Contents(), &test)
		if err != nil {
			return openapi.TransactionResource{}, fmt.Errorf("cannot unmarshal updated test '%s': %w", step, err)
		}

		a.logger.Debug("mapped transaction step",
			zap.Int("index", i),
			zap.String("step", step),
			zap.String("mapped step", test.Spec.GetId()),
		)

		tran.Spec.Steps[i] = test.Spec.GetId()
	}

	return tran, nil
}

func getTypeFromFile(df defFile) (string, error) {
	var raw map[string]any
	err := yaml.Unmarshal(df.Contents(), &raw)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal definition file: %w", err)
	}

	if raw["type"] == nil {
		return "", fmt.Errorf("missing type in definition file")
	}

	defType, ok := raw["type"].(string)
	if !ok {
		return "", fmt.Errorf("type is not a string")
	}

	return defType, nil

}

type envVar struct {
	name         string
	defaultValue string
	userValue    string
}

func (ev envVar) value() string {
	if ev.userValue != "" {
		return ev.userValue
	}

	return ev.defaultValue
}

type envVars []envVar

func (evs envVars) toOpenapi() []openapi.EnvironmentValue {
	vars := make([]openapi.EnvironmentValue, len(evs))
	for i, ev := range evs {
		vars[i] = openapi.EnvironmentValue{
			Key:   openapi.PtrString(ev.name),
			Value: openapi.PtrString(ev.value()),
		}
	}

	return vars
}

func (evs envVars) unique() envVars {
	seen := make(map[string]bool)
	vars := make(envVars, 0, len(evs))
	for _, ev := range evs {
		if seen[ev.name] {
			continue
		}

		seen[ev.name] = true
		vars = append(vars, ev)
	}

	return vars
}

type missingEnvVarsError envVars

func (e missingEnvVarsError) Error() string {
	return fmt.Sprintf("missing env vars: %v", []envVar(e))
}

func (a runTestAction) run(ctx context.Context, df defFile, envID string, ev envVars) (*openapi.ExecuteDefinitionResponse, error) {
	wrapperResp := &openapi.ExecuteDefinitionResponse{}

	defType, err := getTypeFromFile(df)
	if err != nil {
		return wrapperResp, fmt.Errorf("cannot get type from file: %w", err)
	}
	wrapperResp.SetType(defType)

	a.logger.Debug("running definition",
		zap.String("type", defType),
		zap.String("envID", envID),
		zap.Any("envVars", ev),
	)

	runInfo := openapi.RunInformation{
		EnvironmentId: openapi.PtrString(envID),
		Variables:     ev.toOpenapi(),
		Metadata:      getMetadata(),
	}

	switch defType {
	case "Test":
		var test openapi.TestResource
		err = yaml.Unmarshal(df.Contents(), &test)
		if err != nil {
			a.logger.Error("error parsing test", zap.String("content", string(df.Contents())), zap.Error(err))
			return wrapperResp, fmt.Errorf("could not unmarshal test yaml: %w", err)
		}

		req := a.openapiClient.ApiApi.
			RunTest(ctx, test.Spec.GetId()).
			RunInformation(runInfo)

		a.logger.Debug("running test", zap.String("id", test.Spec.GetId()))

		run, resp, err := a.openapiClient.ApiApi.RunTestExecute(req)
		err = a.handleRunError(resp, err)
		if err != nil {
			return wrapperResp, err
		}

		wrapperResp.SetId(test.Spec.GetId())
		wrapperResp.SetRunId(run.GetId())
		wrapperResp.SetUrl(a.config.URL() + "/test/" + test.Spec.GetId() + "/run/" + run.GetId())

	case "Transaction":
		var tran openapi.TransactionResource
		err = yaml.Unmarshal(df.Contents(), &tran)
		if err != nil {
			a.logger.Error("error parsing transaction", zap.String("content", string(df.Contents())), zap.Error(err))
			return wrapperResp, fmt.Errorf("could not unmarshal transaction yaml: %w", err)
		}

		req := a.openapiClient.ApiApi.
			RunTransaction(ctx, tran.Spec.GetId()).
			RunInformation(runInfo)

		a.logger.Debug("running transaction", zap.String("id", tran.Spec.GetId()))

		run, resp, err := a.openapiClient.ApiApi.RunTransactionExecute(req)
		err = a.handleRunError(resp, err)
		if err != nil {
			return wrapperResp, err
		}

		wrapperResp.SetId(tran.Spec.GetId())
		wrapperResp.SetRunId(run.GetId())
		wrapperResp.SetUrl(a.config.URL() + "/transaction/" + tran.Spec.GetId() + "/run/" + run.GetId())
	default:
		return wrapperResp, fmt.Errorf("unknown type: %s", defType)
	}

	a.logger.Debug("definition run",
		zap.String("type", defType),
		zap.String("envID", envID),
		zap.Any("envVars", ev),
		zap.Any("response", wrapperResp),
	)

	return wrapperResp, nil
}

func (a runTestAction) handleRunError(resp *http.Response, reqErr error) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found in server")
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return buildMissingEnvVarsError(body)
	}

	if reqErr != nil {
		a.logger.Error("error running transaction", zap.Error(err), zap.String("body", string(body)))
		return fmt.Errorf("could not run transaction: %w", err)
	}

	return nil
}

func buildMissingEnvVarsError(body []byte) error {
	var missingVarsErrResp openapi.MissingVariablesError
	err := json.Unmarshal(body, &missingVarsErrResp)
	if err != nil {
		return fmt.Errorf("could not unmarshal response body: %w", err)
	}

	missingVars := envVars{}

	for _, missingVarErr := range missingVarsErrResp.MissingVariables {
		for _, missingVar := range missingVarErr.Variables {
			missingVars = append(missingVars, envVar{
				name:         missingVar.GetKey(),
				defaultValue: missingVar.GetDefaultValue(),
			})
		}
	}

	return missingEnvVarsError(missingVars.unique())
}

func (a runTestAction) askForMissingVars(missingVars []envVar) []envVar {
	ui.DefaultUI.Warning("Some variables are required by one or more tests")
	ui.DefaultUI.Info("Fill the values for each variable:")

	filledVariables := make([]envVar, 0, len(missingVars))

	for _, missingVar := range missingVars {
		answer := missingVar
		answer.userValue = ui.DefaultUI.TextInput(missingVar.name, missingVar.defaultValue)
		filledVariables = append(filledVariables, answer)
	}

	a.logger.Debug("filled variables", zap.Any("variables", filledVariables))

	return filledVariables
}

func (a runTestAction) waitForResult(ctx context.Context, run *openapi.ExecuteDefinitionResponse) (*openapi.TestRun, error) {
	return openapi.NewTestRun(), fmt.Errorf("not implemented")
}

func (a runTestAction) writeJUnitReport(ctx context.Context, result *openapi.TestRun, junit string) error {
	return fmt.Errorf("not implemented")
}

func getMetadata() map[string]string {
	ci := cienvironment.DetectCIEnvironment()
	if ci == nil {
		return map[string]string{}
	}

	metadata := map[string]string{
		"name":        ci.Name,
		"url":         ci.URL,
		"buildNumber": ci.BuildNumber,
	}

	if ci.Git != nil {
		metadata["branch"] = ci.Git.Branch
		metadata["tag"] = ci.Git.Tag
		metadata["revision"] = ci.Git.Revision
	}

	return metadata
}

type defFile struct {
	fileutil.File
}
