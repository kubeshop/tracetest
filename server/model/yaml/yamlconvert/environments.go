package yamlconvert

import (
	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func Environment(in model.Environment) yaml.File {
	out := yaml.Environment{}
	dc.DeepCopy(in, &out)

	dc.DeepCopy(in.Values, &out.Values)

	return yaml.File{
		Type: yaml.FileTypeEnvironment,
		Spec: out,
	}
}
