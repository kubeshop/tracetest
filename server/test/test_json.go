package test

import "encoding/json"

type spec struct {
	Selector   Selector    `json:"selector"`
	Name       string      `json:"name,omitempty"`
	Assertions []Assertion `json:"assertions"`
}

type namedAssertion struct {
	Name       string
	Assertions []Assertion
}

type keyValueSpec struct {
	Key   string
	Value namedAssertion
}

func (ts *TestSpec) UnmarshalJSON(data []byte) error {
	spec := spec{}

	err := json.Unmarshal(data, &spec)
	if err != nil || spec.Selector.Query == "" {
		// old format
		oldTestSpecItem := keyValueSpec{}
		err = json.Unmarshal(data, &oldTestSpecItem)
		if err != nil {
			return err
		}

		spec.Selector.Query = SpanQuery(oldTestSpecItem.Key)
		spec.Name = oldTestSpecItem.Value.Name
		spec.Assertions = oldTestSpecItem.Value.Assertions
	}

	ts.Name = spec.Name
	ts.Assertions = spec.Assertions
	ts.Selector = spec.Selector
	return nil
}
