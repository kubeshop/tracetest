package connection

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

const reachabilityTimeout = 5 * time.Second

type Protocol string

var (
	ProtocolHTTP Protocol = "http"
	ProtocolGRPC Protocol = "grpc"
)

type ConnectionTestResult struct {
	ConnectivityTestResult   ConnectionTestStepResult
	AuthenticationTestResult ConnectionTestStepResult
	TraceRetrievalTestResult ConnectionTestStepResult
}

var (
	ErrTraceNotFound        = errors.New("trace not found")
	ErrInvalidConfiguration = errors.New("invalid data store configuration")
	ErrConnectionFailed     = errors.New("could not connect to data store")
)

func (c ConnectionTestResult) HasSucceed() bool {
	return c.AuthenticationTestResult.HasSucceed() && c.ConnectivityTestResult.HasSucceed() && c.TraceRetrievalTestResult.HasSucceed()
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

func CheckReachability(endpoint string, protocol Protocol) error {
	if protocol == ProtocolHTTP {
		address, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")

		if address.Scheme == "https" && address.Port() == "" {
			endpoint = fmt.Sprintf("%s:443", address.Hostname())
		}

		if address.Scheme == "http" && address.Port() == "" {
			endpoint = fmt.Sprintf("%s:80", address.Hostname())
		}
	}

	_, err := net.DialTimeout("tcp", endpoint, reachabilityTimeout)
	if err != nil {
		return err
	}

	return nil
}
