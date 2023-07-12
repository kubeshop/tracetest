package model

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

type (
	Protocol          string
	Status            string
	TestRunEventStage string
	TestRunEventType  string
	PollingType       string
	LogLevel          string
)

var (
	ProtocolHTTP Protocol = "http"
	ProtocolGRPC Protocol = "grpc"
)

var (
	StatusPassed  Status = "passed"
	StatusWarning Status = "warning"
	StatusFailed  Status = "failed"
)

var (
	LogLevelWarn  LogLevel = "warning"
	LogLevelError LogLevel = "error"
)

var (
	StageTrigger TestRunEventStage = "trigger"
	StageTrace   TestRunEventStage = "trace"
	StageTest    TestRunEventStage = "test"
)

var (
	PollingTypePeriodic PollingType = "periodic"
)

type TestRunEvent struct {
	ID                  int64
	Type                string
	Stage               TestRunEventStage
	Title               string
	Description         string
	CreatedAt           time.Time
	TestID              id.ID
	RunID               int
	DataStoreConnection ConnectionResult
	Polling             PollingInfo
	Outputs             []OutputInfo
}

func (e TestRunEvent) ResourceID() string {
	return fmt.Sprintf("test/%s/run/%d/event", e.TestID, e.RunID)
}

type PollingInfo struct {
	Type       PollingType
	IsComplete bool
	Periodic   *PeriodicPollingConfig
}

type PeriodicPollingConfig struct {
	NumberSpans      int
	NumberIterations int
}

type OutputInfo struct {
	LogLevel   LogLevel
	Message    string
	OutputName string
}
