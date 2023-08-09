package testsuite

import (
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
)

const (
	TestSuiteResourceName       = "TestSuite"
	TestSuiteResourceNamePlural = "TestSuites"
)

type TestSuite struct {
	ID          id.ID         `json:"id"`
	CreatedAt   *time.Time    `json:"createdAt,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Version     *int          `json:"version,omitempty"`
	StepIDs     []id.ID       `json:"steps"`
	Steps       []test.Test   `json:"fullSteps,omitempty"`
	Summary     *test.Summary `json:"summary,omitempty"`
}

func setVersion(t *TestSuite, v int) {
	t.Version = &v
}

func (t TestSuite) GetVersion() int {
	if t.Version == nil {
		return 0
	}
	return *t.Version
}

func setCreatedAt(t *TestSuite, d time.Time) {
	t.CreatedAt = &d
}

func (t TestSuite) GetCreatedAt() time.Time {
	if t.CreatedAt == nil {
		return time.Time{}
	}
	return *t.CreatedAt
}

func (t TestSuite) HasID() bool {
	return t.ID != ""
}

func (t TestSuite) GetID() id.ID {
	return t.ID
}

func (t TestSuite) Validate() error {
	return nil
}

func (t TestSuite) NewRun() TestSuiteRun {

	return TestSuiteRun{
		TestSuiteID:      t.ID,
		TestSuiteVersion: t.GetVersion(),
		CreatedAt:        time.Now().UTC(),
		State:            TestSuiteStateCreated,
		Steps:            make([]test.Run, 0, len(t.StepIDs)),
		CurrentTest:      0,
	}
}

type TestSuiteRunState string

const (
	TestSuiteStateCreated   TestSuiteRunState = "CREATED"
	TestSuiteStateExecuting TestSuiteRunState = "EXECUTING"
	TestSuiteStateFailed    TestSuiteRunState = "FAILED"
	TestSuiteStateFinished  TestSuiteRunState = "FINISHED"
)

func (rs TestSuiteRunState) IsFinal() bool {
	return rs == TestSuiteStateFailed || rs == TestSuiteStateFinished
}
