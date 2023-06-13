package test

import (
	"encoding/json"
	"fmt"

	"github.com/fluidtruck/deepcopy"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/server/pkg/maps"
)

type testSpecV1 maps.Ordered[SpanQuery, namedAssertions]

func (v1 testSpecV1) valid() bool {
	valid := true
	specs := maps.Ordered[SpanQuery, namedAssertions](v1)
	specs.ForEach(func(key SpanQuery, val namedAssertions) error {
		if key == "" {
			valid = false
		}
		return nil
	})

	return valid
}

type testSpecV2 []TestSpec

func (v2 testSpecV2) valid() bool {
	for _, spec := range v2 {
		if spec.Selector.Query == "" {
			return false
		}
	}

	return true
}

type namedAssertions struct {
	Name       string
	Assertions []Assertion
}

func (ts *Specs) UnmarshalJSON(data []byte) error {
	v2 := testSpecV2{}
	err := json.Unmarshal(data, &v2)
	if err != nil {
		return err
	}

	if v2.valid() {
		return deepcopy.DeepCopy(v2, ts)
	}

	v1Map := maps.Ordered[SpanQuery, namedAssertions]{}
	v1Map.UnmarshalJSON(data)

	v1 := testSpecV1(v1Map)
	if v1.valid() {
		specs := maps.Ordered[SpanQuery, namedAssertions](v1Map)
		*ts = make([]TestSpec, 0, specs.Len())
		specs.ForEach(func(key SpanQuery, val namedAssertions) error {
			*ts = append(*ts, TestSpec{
				Selector:   Selector{Query: key},
				Name:       val.Name,
				Assertions: val.Assertions,
			})
			return nil
		})

		return nil
	}

	return fmt.Errorf("test spec json version is not supported. Expecting version 1 or 2")
}

type selectorStruct struct {
	Query          string       `json:"query"`
	ParsedSelector SpanSelector `json:"parsedSelector"`
}

func (s *Selector) UnmarshalYAML(data []byte) error {
	selectorStruct := selectorStruct{}
	err := yaml.Unmarshal(data, &selectorStruct)
	if err != nil {
		// This is only the query string
		return yaml.Unmarshal(data, &s.Query)
	}

	s.Query = SpanQuery(selectorStruct.Query)
	s.ParsedSelector = selectorStruct.ParsedSelector
	return nil
}

func (s *Selector) UnmarshalJSON(data []byte) error {
	selectorStruct := selectorStruct{}

	err := json.Unmarshal(data, &selectorStruct)
	if err != nil {
		// This is only the query string
		return json.Unmarshal(data, &s.Query)
	}

	s.Query = SpanQuery(selectorStruct.Query)
	s.ParsedSelector = selectorStruct.ParsedSelector
	return nil
}

type testOutputV1 maps.Ordered[string, Output]

func (v1 testOutputV1) valid() bool {
	orderedMap := maps.Ordered[string, Output](v1)
	for name, item := range orderedMap.Unordered() {
		if name == "" || string(item.Selector) == "" || item.Value == "" {
			return false
		}
	}
	return true
}

type testOutputV2 []Output

func (v2 testOutputV2) valid() bool {
	for _, item := range v2 {
		if item.Name == "" || string(item.Selector) == "" || item.Value == "" {
			return false
		}
	}
	return true
}

func (o *Outputs) UnmarshalJSON(data []byte) error {
	v2 := testOutputV2{}
	err := json.Unmarshal(data, &v2)
	if err == nil && v2.valid() {
		*o = []Output(v2)
		return nil
	}

	v1 := maps.Ordered[string, Output]{}
	err = json.Unmarshal(data, &v1)
	if err == nil && testOutputV1(v1).valid() {
		newOutputs := make(Outputs, 0)
		v1Map := maps.Ordered[string, Output](v1)

		v1Map.ForEach(func(key string, val Output) error {
			newOutputs = append(newOutputs, Output{
				Name:     key,
				Selector: val.Selector,
				Value:    val.Value,
			})

			return nil
		})

		*o = newOutputs
		return nil
	}

	return fmt.Errorf("test output json version is not supported. Expecting version 1 or 2")
}
