package trigger_preprocessor

import (
	"fmt"

	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type grpc struct {
	logger *zap.Logger
}

func GRPC(logger *zap.Logger) grpc {
	return grpc{logger}
}

func (g grpc) Type() trigger.TriggerType {
	return trigger.TriggerTypeGRPC
}

func (g grpc) Preprocess(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	// protobuf file can be defined separately in a file like: protobufFile: ./file.proto
	definedPBFile := test.Spec.Trigger.Grpc.GetProtobufFile()
	if !isValidFilePath(definedPBFile, input.AbsDir()) {
		g.logger.Debug("protobuf file is not a file path", zap.String("protobufFile", definedPBFile))
		return test, nil
	}

	pbFilePath := input.RelativeFile(definedPBFile)
	g.logger.Debug("protobuf file", zap.String("path", pbFilePath))

	pbFile, err := fileutil.Read(pbFilePath)
	if err != nil {
		return test, fmt.Errorf(`cannot read protobuf file: %w`, err)
	}
	g.logger.Debug("protobuf file contents", zap.String("contents", string(pbFile.Contents())))

	test.Spec.Trigger.Grpc.SetProtobufFile(string(pbFile.Contents()))

	return test, nil
}

func isValidFilePath(filePath, testFile string) bool {
	if fileutil.LooksLikeRelativeFilePath(filePath) {
		// if looks like a relative file path, test if it exists
		return fileutil.IsFilePathToRelativeDir(filePath, testFile)
	}

	// it could be an absolute file path, test it
	return fileutil.IsFilePath(filePath)
}
