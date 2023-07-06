package test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"go.opentelemetry.io/otel/trace"
)

type testSpecV1 maps.Ordered[SpanQuery, namedAssertions]

func (v1 testSpecV1) valid() bool {
	valid := true
	specs := maps.Ordered[SpanQuery, namedAssertions](v1)
	specs.ForEach(func(key SpanQuery, val namedAssertions) error {
		anyEmptyAssertion := false
		for _, assertion := range val.Assertions {
			if assertion == "" {
				anyEmptyAssertion = true
			}
		}

		if key == "" && anyEmptyAssertion {
			valid = false
		}
		return nil
	})

	return valid
}

type testSpecV2 []TestSpec

func (v2 testSpecV2) valid() bool {
	for _, spec := range v2 {
		// since we can have an empty selector, to check if go
		// sent an empty struct we need to see if we have an empty selector and no assertions
		if spec.Selector == "" && len(spec.Assertions) == 0 {
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
				Selector:   SpanQuery(key),
				Name:       val.Name,
				Assertions: val.Assertions,
			})
			return nil
		})

		return nil
	}

	return fmt.Errorf("test json version is not supported. Expecting version 1 or 2")
}

func (s *SpanQuery) UnmarshalJSON(data []byte) error {
	selectorStruct := struct {
		Query          string       `json:"query"`
		ParsedSelector SpanSelector `json:"parsedSelector"`
	}{}

	err := json.Unmarshal(data, &selectorStruct)
	if err != nil {
		// This is only the query string
		var query string
		err = json.Unmarshal(data, &query)
		if err != nil {
			return err
		}

		*s = SpanQuery(query)
		return nil
	}

	*s = SpanQuery(selectorStruct.Query)
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

func (sar SpanAssertionResult) MarshalJSON() ([]byte, error) {
	sid := ""
	if sar.SpanID != nil {
		sid = sar.SpanID.String()
	}
	return json.Marshal(&struct {
		SpanID        *string
		ObservedValue string
		CompareErr    string
	}{
		SpanID:        &sid,
		ObservedValue: sar.ObservedValue,
		CompareErr:    errToString(sar.CompareErr),
	})
}

func (sar *SpanAssertionResult) UnmarshalJSON(data []byte) error {
	aux := struct {
		SpanID        string
		ObservedValue string
		CompareErr    string
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var sid *trace.SpanID
	if aux.SpanID != "" {
		s, err := trace.SpanIDFromHex(aux.SpanID)
		if err != nil {
			return err
		}
		sid = &s
	}

	sar.SpanID = sid
	sar.ObservedValue = aux.ObservedValue
	if err := stringToErr(aux.CompareErr); err != nil {
		if err.Error() == comparator.ErrNoMatch.Error() {
			err = comparator.ErrNoMatch
		}

		sar.CompareErr = err
	}

	return nil
}

func errToString(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

func stringToErr(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}
