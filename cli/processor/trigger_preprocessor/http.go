package trigger_preprocessor

import (
	"fmt"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type http struct {
	logger *zap.Logger
}

func HTTP(logger *zap.Logger) http {
	return http{logger}
}

func (g http) Type() trigger.TriggerType {
	return trigger.TriggerTypeHTTP
}

func (g http) Preprocess(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	// body can be defined separately in a file like: body: ./body.json
	definedBodyFile := test.Spec.Trigger.HttpRequest.GetBody()
	if !isValidFilePath(definedBodyFile, input.AbsDir()) {
		g.logger.Debug("body file is not a file path", zap.String("protobufFile", definedBodyFile))
		return test, nil
	}

	bodyFilePath := input.RelativeFile(definedBodyFile)
	g.logger.Debug("protobuf file", zap.String("path", bodyFilePath))

	bodyFile, err := fileutil.Read(definedBodyFile)
	if err != nil {
		return test, fmt.Errorf(`cannot read protobuf file: %w`, err)
	}
	g.logger.Debug("protobuf file contents", zap.String("contents", string(bodyFile.Contents())))

	test.Spec.Trigger.HttpRequest.SetBody(string(bodyFile.Contents()))

	return test, nil
}
