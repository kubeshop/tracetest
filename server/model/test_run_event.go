package model

import "time"

type (
	Protocol string
	Status   string
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

type TestRunEvent struct {
	Type                string
	Stage               string
	Description         string
	CreatedAt           time.Time
	TestId              string
	RunId               string
	DataStoreConnection ConnectionResult
	Polling             PollingInfo
	Outputs             []OutputInfo
}

type PollingInfo struct {
	Type                string
	ReasonNextIteration string
	IsComplete          bool
	Periodic            *PeriodicPollingConfig
}

type PeriodicPollingConfig struct {
	NumberSpans      int32
	NumberIterations int32
}

type OutputInfo struct {
	LogLevel   string
	Message    string
	OutputName string
}
