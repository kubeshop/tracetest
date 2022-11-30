package yamlconvert

import (
	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func Environment(in model.Environment) yaml.File {
	out := yaml.Environment{}
	dc.DeepCopy(in, &out)

	valueList := make([]yaml.EnvironmentValue, len(in.Values))

	for index, value := range in.Values {
		valueList[index] = yaml.EnvironmentValue{
			Key:   value.Key,
			Value: value.Value,
		}
	}

	out.Values = valueList

	return yaml.File{
		Type: yaml.FileTypeEnvironment,
		Spec: out,
	}
}
