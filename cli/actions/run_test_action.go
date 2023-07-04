package actions

import (
	"context"
	"fmt"

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
	jsonFormat    resourcemanager.Format
	cliExit       func(int)
}

func NewRunTestAction(
	config config.Config,
	logger *zap.Logger,
	openapiClient *openapi.APIClient,
	tests resourcemanager.Client,
	environments resourcemanager.Client,
	cliExit func(int),
) runTestAction {
	jsonFormat, err := resourcemanager.Formats.Get(resourcemanager.FormatJSON)
	if err != nil {
		panic(fmt.Errorf("could not get json format: %w", err))
	}
	return runTestAction{
		config,
		logger,
		openapiClient,
		tests,
		environments,
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
	return "", fmt.Errorf("not implemented")
}

func (a runTestAction) apply(ctx context.Context, defFile defFile) (defFile, error) {
	return defFile, nil
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

func (a runTestAction) run(ctx context.Context, defFile defFile, envID string, envVars []envVar) (*openapi.TestRun, error) {
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
