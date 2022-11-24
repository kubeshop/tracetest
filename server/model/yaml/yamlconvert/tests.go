package yamlconvert

import (
	dc "github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

func Test(in model.Test) yaml.File {
	out := yaml.Test{}
	dc.DeepCopy(in, &out)
	dc.DeepCopy(in.ServiceUnderTest, &out.Trigger)

	in.Specs.ForEach(func(key model.SpanQuery, val model.NamedAssertions) error {
		spec := yaml.TestSpec{
			Selector: string(key),
			Name:     val.Name,
		}
		dc.DeepCopy(val.Assertions, &spec.Assertions)
		out.Specs = append(out.Specs, spec)

		return nil
	})

	in.Outputs.ForEach(func(key string, val model.Output) error {
		out.Outputs = append(out.Outputs, yaml.Output{
			Name:     key,
			Selector: string(val.Selector),
			Value:    val.Value,
		})

		return nil
	})

	return yaml.File{
		Type: yaml.FileTypeTest,
		Spec: out,
	}
}
