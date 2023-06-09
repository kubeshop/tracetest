package test

import (
	"encoding/json"
	"fmt"

	"github.com/fluidtruck/deepcopy"
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

	return fmt.Errorf("test json version is not supported. Expecting version 1 or 2")
}

func (s *Selector) UnmarshalJSON(data []byte) error {
	selectorStruct := struct {
		Query          string       `json:"query"`
		ParsedSelector SpanSelector `json:"parsedSelector"`
	}{}

	err := json.Unmarshal(data, &selectorStruct)
	if err != nil {
		// This is only the query string
		return json.Unmarshal(data, &s.Query)
	}

	s.Query = SpanQuery(selectorStruct.Query)
	s.ParsedSelector = selectorStruct.ParsedSelector
	return nil
}
