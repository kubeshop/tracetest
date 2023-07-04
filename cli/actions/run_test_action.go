package actions

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
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
	defFile := defFile{f}
	a.logger.Debug("Definition file read", zap.String("absolutePath", defFile.AbsPath()))

	envID, err := a.resolveEnvID(ctx, args.EnvID)
	if err != nil {
		return fmt.Errorf("cannot resolve environment id: %w", err)
	}
	a.logger.Debug("env resolved", zap.String("ID", envID))

	defFile, err = a.apply(ctx, defFile)
	if err != nil {
		return fmt.Errorf("cannot apply definition file: %w", err)
	}

	var run *openapi.TestRun
	var envVars []envVar

	// iterate until we have all env vars,
	// or the server returns an actual error
	for {
		run, err = a.run(ctx, defFile, envID, envVars)
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

	filePath := envID

	a.logger.Debug("envID is a file path", zap.String("filePath", filePath))
	updatedEnv, err := a.environments.Apply(ctx, filePath, a.yamlFormat)
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
	var test openapi.TestResource
	err := yaml.Unmarshal(df.Contents(), &test)
	if err != nil {
		a.logger.Error("error parsing test", zap.String("content", string(df.Contents())), zap.Error(err))
		return df, fmt.Errorf("could not unmarshal test yaml: %w", err)
	}

	test, err = consolidateGRPCFile(df, test)
	if err != nil {
		return df, fmt.Errorf("could not consolidate grpc file: %w", err)
	}

	a.logger.Debug("applying test",
		zap.String("absolutePath", df.AbsPath()),
		zap.String("id", test.Spec.GetId()),
	)
	updated, err := a.tests.Apply(ctx, df.AbsPath(), a.yamlFormat)
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
	var tran openapi.TransactionResource
	err := yaml.Unmarshal(df.Contents(), &tran)
	if err != nil {
		a.logger.Error("error parsing transaction", zap.String("content", string(df.Contents())), zap.Error(err))
		return df, fmt.Errorf("could not unmarshal transaction yaml: %w", err)
	}

	a.logger.Debug("applying transaction",
		zap.String("absolutePath", df.AbsPath()),
		zap.String("id", tran.Spec.GetId()),
	)

	tran, err = a.mapTransactionSteps(ctx, df, tran)
	if err != nil {
		return df, fmt.Errorf("could not map transaction steps: %w", err)
	}

	updated, err := a.transactions.Apply(ctx, df.AbsPath(), a.yamlFormat)
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

type missingEnvVarsError []envVar

func (e missingEnvVarsError) Error() string {
	return fmt.Sprintf("missing env vars: %v", []envVar(e))
}

func (a runTestAction) run(ctx context.Context, df defFile, envID string, envVars []envVar) (*openapi.TestRun, error) {
	return openapi.NewTestRun(), fmt.Errorf("not implemented")
}

func (a runTestAction) askForMissingVars(missingVars []envVar) []envVar {
	return []envVar{}
}

func (a runTestAction) waitForResult(ctx context.Context, run *openapi.TestRun) (*openapi.TestRun, error) {
	return openapi.NewTestRun(), fmt.Errorf("not implemented")
}

func (a runTestAction) writeJUnitReport(ctx context.Context, result *openapi.TestRun, junit string) error {
	return fmt.Errorf("not implemented")
}

type defFile struct {
	fileutil.File
}
