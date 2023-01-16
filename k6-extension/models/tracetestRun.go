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

func NewRun(cliResponse string) *TracetestRun {
	var tracetestRun TracetestRun
	json.Unmarshal([]byte(cliResponse), &tracetestRun)

	return &tracetestRun
}
