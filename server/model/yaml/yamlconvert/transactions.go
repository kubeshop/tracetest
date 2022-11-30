package yamlconvert

import (
	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func Transaction(in model.Transaction) yaml.File {
	out := yaml.Transaction{}
	dc.DeepCopy(in, &out)

	stepList := make([]string, len(in.Steps))

	for index, step := range in.Steps {
		stepList[index] = string(step.ID)
	}

	out.Steps = stepList

	return yaml.File{
		Type: yaml.FileTypeTransaction,
		Spec: out,
	}
}
