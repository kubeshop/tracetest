package cmdutil

import (
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/preprocessor"
)

var ()

func GetVariableSetClient(preprocessor preprocessor.Preprocessor) resourcemanager.Client {
	httpClient := &resourcemanager.HTTPClient{}

	variableSetClient := resourcemanager.NewClient(
		httpClient, GetLogger(),
		"variableset", "variablesets",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "DESCRIPTION", Path: "spec.description"},
			},
		}),
		resourcemanager.WithResourceType("VariableSet"),
		resourcemanager.WithApplyPreProcessor(preprocessor.Preprocess),
		resourcemanager.WithDeprecatedAlias("Environment"),
	)

	return variableSetClient
}
