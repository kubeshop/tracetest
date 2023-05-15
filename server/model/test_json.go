package model

import (
	"encoding/json"
	"time"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
)

type jsonSpec struct {
	Name       string      `json:"name"`
	Selector   string      `json:"selector"`
	Assertions []Assertion `json:"assertions"`
}

type jsonOutput struct {
	Name     string `json:"name"`
	Selector string `json:"selector"`
	Value    string `json:"value"`
}

type jsonTest struct {
	ID               id.ID        `json:"id"`
	CreatedAt        time.Time    `json:"createdAt,omitempty"`
	Name             string       `json:"name"`
	Description      string       `json:"description"`
	Version          int          `json:"version,omitempty"`
	ServiceUnderTest Trigger      `json:"serviceUnderTest"`
	SpecsJSON        []jsonSpec   `json:"specs"`
	OutputsJSON      []jsonOutput `json:"outputs"`
	Summary          Summary      `json:"summary,omitempty"`
}

func (t Test) MarshalJSON() ([]byte, error) {
	jt := jsonTest{}
	err := deepcopy.DeepCopy(t, &jt)
	if err != nil {
		return nil, err
	}

	jt.SpecsJSON = make([]jsonSpec, 0, t.Specs.Len())
	t.Specs.ForEach(func(key SpanQuery, val NamedAssertions) error {
		// assertions := make([]string, len(val.Assertions))
		// for i, a := range val.Assertions {
		// 	assertions[i] = string(a)
		// }

		jt.SpecsJSON = append(jt.SpecsJSON, jsonSpec{
			Name: val.Name,
			// Assertions: assertions,
			Assertions: val.Assertions,
			Selector:   string(key),
		})

		return nil
	})

	jt.OutputsJSON = make([]jsonOutput, 0, t.Outputs.Len())
	t.Outputs.ForEach(func(key string, val Output) error {
		jt.OutputsJSON = append(jt.OutputsJSON, jsonOutput{
			Name:     key,
			Selector: string(val.Selector),
			Value:    val.Value,
		})

		return nil
	})

	return json.Marshal(jt)
}

func (t *Test) UnmarshalJSON(data []byte) error {
	jt := jsonTest{}
	err := json.Unmarshal(data, &jt)
	if err != nil {
		return err
	}

	specs := maps.Ordered[SpanQuery, NamedAssertions]{}
	for _, spec := range jt.SpecsJSON {
		specs, err = specs.Add(SpanQuery(spec.Selector), NamedAssertions{
			Name:       spec.Name,
			Assertions: spec.Assertions,
		})
		if err != nil {
			return err
		}
	}

	outputs := maps.Ordered[string, Output]{}
	for _, output := range jt.OutputsJSON {
		outputs, err = outputs.Add(output.Name, Output{
			Selector: SpanQuery(output.Selector),
			Value:    output.Value,
		})
		if err != nil {
			return err
		}
	}

	err = deepcopy.DeepCopy(jt, t)
	if err != nil {
		return err
	}

	t.Specs = specs
	t.Outputs = outputs

	return nil
}

func (t Test) MarshalYAML() ([]byte, error) {
	return t.MarshalJSON()
}

func (t *Test) UnmarshalYAML(data []byte) error {
	return t.UnmarshalJSON(data)
}
