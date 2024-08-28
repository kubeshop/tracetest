package trigger_preprocessor

import (
	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

type graphql struct {
	logger *zap.Logger
}

func GRAPHQL(logger *zap.Logger) graphql {
	return graphql{logger: cmdutil.GetLogger()}
}

func (g graphql) Type() trigger.TriggerType {
	return trigger.TriggerTypeGraphql
}

func (g graphql) Preprocess(input fileutil.File, test openapi.TestResource) (openapi.TestResource, error) {
	// query can be defined separately in a file like: query: ./query.graphql
	rawQuery := test.Spec.Trigger.Graphql.Body.GetQuery()
	if isValidFilePath(rawQuery, input.AbsDir()) {
		queryFilePath := input.RelativeFile(rawQuery)
		g.logger.Debug("script file", zap.String("path", queryFilePath))

		queryFile, err := fileutil.Read(queryFilePath)
		if err == nil {
			g.logger.Debug("script file contents", zap.String("contents", string(queryFile.Contents())))
			test.Spec.Trigger.Graphql.Body.SetQuery(string(queryFile.Contents()))
		}
	}

	// schema can be defined separately in a file like: schema: ./schema.graphql
	rawSchema := test.Spec.Trigger.Graphql.GetSchema()
	if isValidFilePath(rawSchema, input.AbsDir()) {
		schemaFilePath := input.RelativeFile(rawSchema)
		g.logger.Debug("script file", zap.String("path", schemaFilePath))

		schemaFile, err := fileutil.Read(schemaFilePath)
		if err == nil {
			g.logger.Debug("script file contents", zap.String("contents", string(schemaFile.Contents())))
			test.Spec.Trigger.Graphql.SetSchema(string(schemaFile.Contents()))
		}
	}

	return test, nil
}
