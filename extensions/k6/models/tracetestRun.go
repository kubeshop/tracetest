package models

import "encoding/json"

type Test struct {
	ID   string
	Name string
}

type TestRun struct {
	ID      string
	TraceID string
}

type TracetestRun struct {
	Test    Test
	TestRun TestRun
}

func NewRun(cliResponse string) (*TracetestRun, error) {
	var tracetestRun TracetestRun
	err := json.Unmarshal([]byte(cliResponse), &tracetestRun)

	if err != nil {
		return nil, err
	}

	return &tracetestRun, nil
}
