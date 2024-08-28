package trigger_preprocessor

import (
	"fmt"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type playwrightengine struct {
	logger *zap.Logger
}

func PLAYWRIGHTENGINE(logger *zap.Logger) playwrightengine {
	return playwrightengine{logger: cmdutil.GetLogger()}
}

func (g playwrightengine) Type() trigger.TriggerType {
	return trigger.TriggerTypePlaywrightEngine
}

func (g playwrightengine) Preprocess(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	// script can be defined separately in a file like: script: ./script.js
	definedScriptFile := test.Spec.Trigger.PlaywrightEngine.GetScript()
	if !isValidFilePath(definedScriptFile, input.AbsDir()) {
		g.logger.Debug("script file is not a file path", zap.String("protobufFile", definedScriptFile))
		return test, nil
	}

	scriptFilePath := input.RelativeFile(definedScriptFile)
	g.logger.Debug("script file", zap.String("path", scriptFilePath))

	scriptFile, err := fileutil.Read(scriptFilePath)
	if err != nil {
		return test, fmt.Errorf(`cannot read script file: %w`, err)
	}
	g.logger.Debug("script file contents", zap.String("contents", string(scriptFile.Contents())))

	test.Spec.Trigger.PlaywrightEngine.SetScript(string(scriptFile.Contents()))

	return test, nil
}
