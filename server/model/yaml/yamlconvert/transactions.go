package yamlconvert

import (
	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func Transaction(in model.Transaction) yaml.File {
	out := yaml.Transaction{}
	dc.DeepCopy(in, &out)

	out.Steps = make([]string, 0, len(in.Steps))

	for _, step := range in.Steps {
		out.Steps = append(out.Steps, step.ID.String()
	}

	return yaml.File{
		Type: yaml.FileTypeTransaction,
		Spec: out,
	}
}
