package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	cienvironment "github.com/cucumber/ci-environment/go"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/ui"
	"github.com/kubeshop/tracetest/cli/utils"
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
	jsonFormat    resourcemanager.Format
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

	jsonFormat, err := resourcemanager.Formats.Get(resourcemanager.FormatJSON)
	if err != nil {
		panic(fmt.Errorf("could not get json format: %w", err))
	}

	return runTestAction{
		config,
		logger,
		openapiClient,
		environments,
		tests,
		transactions,
		yamlFormat,
		jsonFormat,
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

	var run runResult
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
		fmt.Println(a.formatResult(run, false))
		a.cliExit(0)
	}

	result, err := a.waitForResult(ctx, run)
	if err != nil {
		return fmt.Errorf("cannot wait for test result: %w", err)
	}

	fmt.Println(a.formatResult(result, true))
	a.cliExit(a.exitCode(result))

	if args.JUnit != "" {
		err := a.writeJUnitReport(ctx, result, args.JUnit)
		if err != nil {
			return fmt.Errorf("cannot write junit report: %w", err)
		}
	}

	return nil
}

func (a runTestAction) exitCode(res runResult) int {
	switch res.ResourceType {
	case "Test":
		if !res.Run.(openapi.TestRun).Result.GetAllPassed() {
			return 1
		}
	case "Transaction":
		for _, step := range res.Run.(openapi.TransactionRun).Steps {
			if !step.Result.GetAllPassed() {
				return 1
			}
		}
	}
	return 0
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

	a.logger.Debug("applying test",
		zap.String("absolutePath", df.AbsPath()),
	)

	updated, err := a.tests.Apply(ctx, df.File, a.yamlFormat)
	if err != nil {
		return df, fmt.Errorf("could not read test file: %w", err)
	}

	df = defFile{fileutil.New(df.AbsPath(), []byte(updated))}

	return df, nil
}

func (a runTestAction) applyTransaction(ctx context.Context, df defFile) (defFile, error) {
	df, err := a.injectLocalEnvVars(ctx, df)
	if err != nil {
		return df, fmt.Errorf("cannot inject local env vars: %w", err)
	}

	a.logger.Debug("applying transaction",
		zap.String("absolutePath", df.AbsPath()),
	)

	updated, err := a.transactions.Apply(ctx, df.File, a.yamlFormat)
	if err != nil {
		return df, fmt.Errorf("could not read transaction file: %w", err)
	}

	df = defFile{fileutil.New(df.AbsPath(), []byte(updated))}
	return df, nil
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
	Name         string
	DefaultValue string
	UserValue    string
}

func (ev envVar) value() string {
	if ev.UserValue != "" {
		return ev.UserValue
	}

	return ev.DefaultValue
}

type envVars []envVar

func (evs envVars) toOpenapi() []openapi.EnvironmentValue {
	vars := make([]openapi.EnvironmentValue, len(evs))
	for i, ev := range evs {
		vars[i] = openapi.EnvironmentValue{
			Key:   openapi.PtrString(ev.Name),
			Value: openapi.PtrString(ev.value()),
		}
	}

	return vars
}

func (evs envVars) unique() envVars {
	seen := make(map[string]bool)
	vars := make(envVars, 0, len(evs))
	for _, ev := range evs {
		if seen[ev.Name] {
			continue
		}

		seen[ev.Name] = true
		vars = append(vars, ev)
	}

	return vars
}

type missingEnvVarsError envVars

func (e missingEnvVarsError) Error() string {
	return fmt.Sprintf("missing env vars: %v", []envVar(e))
}

type runResult struct {
	ResourceType string
	Resource     any
	Run          any
}

func (a runTestAction) run(ctx context.Context, df defFile, envID string, ev envVars) (runResult, error) {
	res := runResult{}

	defType, err := getTypeFromFile(df)
	if err != nil {
		return res, fmt.Errorf("cannot get type from file: %w", err)
	}
	res.ResourceType = defType

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
			return res, fmt.Errorf("could not unmarshal test yaml: %w", err)
		}

		req := a.openapiClient.ApiApi.
			RunTest(ctx, test.Spec.GetId()).
			RunInformation(runInfo)

		a.logger.Debug("running test", zap.String("id", test.Spec.GetId()))

		run, resp, err := a.openapiClient.ApiApi.RunTestExecute(req)
		err = a.handleRunError(resp, err)
		if err != nil {
			return res, err
		}

		full, err := a.tests.Get(ctx, test.Spec.GetId(), a.jsonFormat)
		if err != nil {
			return res, fmt.Errorf("cannot get full test '%s': %w", test.Spec.GetId(), err)
		}
		err = json.Unmarshal([]byte(full), &test)
		if err != nil {
			return res, fmt.Errorf("cannot get full test '%s': %w", test.Spec.GetId(), err)
		}

		res.Resource = test
		res.Run = *run

	case "Transaction":
		var tran openapi.TransactionResource
		err = yaml.Unmarshal(df.Contents(), &tran)
		if err != nil {
			a.logger.Error("error parsing transaction", zap.String("content", string(df.Contents())), zap.Error(err))
			return res, fmt.Errorf("could not unmarshal transaction yaml: %w", err)
		}

		req := a.openapiClient.ApiApi.
			RunTransaction(ctx, tran.Spec.GetId()).
			RunInformation(runInfo)

		a.logger.Debug("running transaction", zap.String("id", tran.Spec.GetId()))

		run, resp, err := a.openapiClient.ApiApi.RunTransactionExecute(req)
		err = a.handleRunError(resp, err)
		if err != nil {
			return res, err
		}

		full, err := a.transactions.Get(ctx, tran.Spec.GetId(), a.jsonFormat)
		if err != nil {
			return res, fmt.Errorf("cannot get full transaction '%s': %w", tran.Spec.GetId(), err)
		}
		err = json.Unmarshal([]byte(full), &tran)
		if err != nil {
			return res, fmt.Errorf("cannot get full transaction '%s': %w", tran.Spec.GetId(), err)
		}

		res.Resource = tran
		res.Run = *run
	default:
		return res, fmt.Errorf("unknown type: %s", defType)
	}

	a.logger.Debug("definition run",
		zap.String("type", defType),
		zap.String("envID", envID),
		zap.Any("envVars", ev),
		zap.Any("response", res),
	)

	return res, nil
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
				Name:         missingVar.GetKey(),
				DefaultValue: missingVar.GetDefaultValue(),
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
		answer.UserValue = ui.DefaultUI.TextInput(missingVar.Name, missingVar.DefaultValue)
		filledVariables = append(filledVariables, answer)
	}

	a.logger.Debug("filled variables", zap.Any("variables", filledVariables))

	return filledVariables
}

func (a runTestAction) formatResult(result runResult, hasResults bool) string {
	switch result.ResourceType {
	case "Test":
		return a.formatTestResult(result, hasResults)
	case "Transaction":
		return a.formatTransactionResult(result, hasResults)
	}
	return ""
}

func (a runTestAction) formatTestResult(result runResult, hasResults bool) string {
	test := result.Resource.(openapi.TestResource)
	run := result.Run.(openapi.TestRun)

	tro := formatters.TestRunOutput{
		HasResults: hasResults,
		Test:       test.GetSpec(),
		Run:        run,
	}

	formatter := formatters.TestRun(a.config, true)
	return formatter.Format(tro)
}
func (a runTestAction) formatTransactionResult(result runResult, hasResults bool) string {
	tran := result.Resource.(openapi.TransactionResource)
	run := result.Run.(openapi.TransactionRun)

	tro := formatters.TransactionRunOutput{
		HasResults:  hasResults,
		Transaction: tran.GetSpec(),
		Run:         run,
	}

	return formatters.
		TransactionRun(a.config, true).
		Format(tro)
}

func (a runTestAction) waitForResult(ctx context.Context, run runResult) (runResult, error) {
	switch run.ResourceType {
	case "Test":
		tr, err := a.waitForTestResult(ctx, run)
		if err != nil {
			return run, err
		}

		run.Run = tr
		return run, nil

	case "Transaction":
		tr, err := a.waitForTransactionResult(ctx, run)
		if err != nil {
			return run, err
		}

		run.Run = tr
		return run, nil
	}
	return run, fmt.Errorf("unknown resource type: %s", run.ResourceType)
}

func (a runTestAction) waitForTestResult(ctx context.Context, result runResult) (openapi.TestRun, error) {
	var (
		testRun   openapi.TestRun
		lastError error
		wg        sync.WaitGroup
	)

	test := result.Resource.(openapi.TestResource)
	run := result.Run.(openapi.TestRun)

	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: change to websockets
	go func() {
		for range ticker.C {
			readyTestRun, err := a.isTestReady(ctx, test.Spec.GetId(), run.GetId())
			if err != nil {
				lastError = err
				wg.Done()
				return
			}

			if readyTestRun != nil {
				testRun = *readyTestRun
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.TestRun{}, lastError
	}

	return testRun, nil
}

func (a runTestAction) isTestReady(ctx context.Context, testID string, testRunID string) (*openapi.TestRun, error) {
	runID, err := strconv.Atoi(testRunID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction run id format: %w", err)
	}

	req := a.openapiClient.ApiApi.GetTestRun(ctx, testID, int32(runID))
	run, _, err := a.openapiClient.ApiApi.GetTestRunExecute(req)
	if err != nil {
		return &openapi.TestRun{}, fmt.Errorf("could not execute GetTestRun request: %w", err)
	}

	if utils.RunStateIsFinished(run.GetState()) {
		return run, nil
	}

	return nil, nil
}

func (a runTestAction) waitForTransactionResult(ctx context.Context, result runResult) (openapi.TransactionRun, error) {
	var (
		transactionRun openapi.TransactionRun
		lastError      error
		wg             sync.WaitGroup
	)

	tran := result.Resource.(openapi.TransactionResource)
	run := result.Run.(openapi.TransactionRun)

	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: change to websockets
	go func() {
		for range ticker.C {
			readyTransactionRun, err := a.isTransactionReady(ctx, tran.Spec.GetId(), run.GetId())
			if err != nil {
				lastError = err
				wg.Done()
				return
			}

			if readyTransactionRun != nil {
				transactionRun = *readyTransactionRun
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.TransactionRun{}, lastError
	}

	return transactionRun, nil
}

func (a runTestAction) isTransactionReady(ctx context.Context, transactionID, transactionRunID string) (*openapi.TransactionRun, error) {
	runID, err := strconv.Atoi(transactionRunID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction run id format: %w", err)
	}

	req := a.openapiClient.ApiApi.GetTransactionRun(ctx, transactionID, int32(runID))
	run, _, err := a.openapiClient.ApiApi.GetTransactionRunExecute(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute GetTestRun request: %w", err)
	}

	if utils.RunStateIsFinished(run.GetState()) {
		return run, nil
	}

	return nil, nil
}

func (a runTestAction) writeJUnitReport(ctx context.Context, result runResult, outputFile string) error {
	test := result.Resource.(openapi.TestResource)
	run := result.Run.(openapi.TestRun)
	runID, err := strconv.Atoi(run.GetId())
	if err != nil {
		return fmt.Errorf("invalid run id format: %w", err)
	}

	req := a.openapiClient.ApiApi.GetRunResultJUnit(ctx, test.Spec.GetId(), int32(runID))
	junit, _, err := a.openapiClient.ApiApi.GetRunResultJUnitExecute(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	f, err := os.Create(junit)
	if err != nil {
		return fmt.Errorf("could not create junit output file: %w", err)
	}

	_, err = f.WriteString(outputFile)

	return err

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
