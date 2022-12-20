package tracedb

import (
	"net"
	"strings"
	"time"
)

type ConnectionTestResult struct {
	ConnectivityTestResult   ConnectionTestStepResult
	AuthenticationTestResult ConnectionTestStepResult
	TraceRetrivalTestResult  ConnectionTestStepResult
}

func (c ConnectionTestResult) HasSucceed() bool {
	return c.AuthenticationTestResult.HasSucceed() && c.ConnectivityTestResult.HasSucceed() && c.TraceRetrivalTestResult.HasSucceed()
}

type ConnectionTestStepResult struct {
	OperationDescription string
	Error                error
}

func (r ConnectionTestStepResult) HasSucceed() bool {
	return r.Error == nil
}

func (r ConnectionTestStepResult) IsSet() bool {
	return r.OperationDescription != ""
}

func isReachable(endpoint string) (bool, error) {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	_, err := net.DialTimeout("tcp", endpoint, 5*time.Second)
	if err != nil {
		return false, err
	}

	return true, nil
}
