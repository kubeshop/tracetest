package models

import (
	"fmt"

	"github.com/kubeshop/tracetest/extensions/k6/openapi"
)

type Test struct {
	ID   string
	Name string
}

type TracetestRun struct {
	TestId  string
	TestRun *openapi.TestRun
}

func (tr *TracetestRun) Summary(baseUrl string) string {
	runUrl := fmt.Sprintf("%s/test/%s/run/%s", baseUrl, tr.TestId, *tr.TestRun.Id)

	failingSpecs := false
	if tr.TestRun != nil && tr.TestRun.Result != nil && tr.TestRun.Result.AllPassed != nil {
		failingSpecs = !*tr.TestRun.Result.AllPassed
	}

	return fmt.Sprintf("RunState=%s FailingSpecs=%t, TracetestURL= %s", *tr.TestRun.State, failingSpecs, runUrl)
}

func (tr *TracetestRun) IsSuccessful() bool {
	if tr.TestRun != nil && tr.TestRun.Result != nil && tr.TestRun.Result.AllPassed != nil {
		return *tr.TestRun.Result.AllPassed
	}

	return false
}
