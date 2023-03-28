package connection

import "context"

type TestStep interface {
	TestConnection(ctx context.Context) ConnectionTestStepResult
}

type TesterOption func(*Tester)

type Tester struct {
	portLinterStep         TestStep
	connectivityTestStep   TestStep
	authenticationTestStep TestStep
	pollingTestStep        TestStep
}

func NewTester(opts ...TesterOption) Tester {
	tester := Tester{}

	for _, opt := range opts {
		opt(&tester)
	}

	return tester
}

func (t *Tester) TestConnection(ctx context.Context) (res ConnectionTestResult) {
	if t.portLinterStep != nil {
		res.EndpointLintTestResult = t.portLinterStep.TestConnection(ctx)
		if res.EndpointLintTestResult.Error != nil {
			res.EndpointLintTestResult.Status = StatusFailed
			return
		}
	}

	res.ConnectivityTestResult = t.connectivityTestStep.TestConnection(ctx)
	if res.ConnectivityTestResult.Error != nil {
		res.ConnectivityTestResult.Status = StatusFailed
		return
	}

	res.AuthenticationTestResult = t.authenticationTestStep.TestConnection(ctx)
	if res.AuthenticationTestResult.Error != nil {
		res.AuthenticationTestResult.Status = StatusFailed
		return
	}

	res.TraceRetrievalTestResult = t.pollingTestStep.TestConnection(ctx)
	if res.TraceRetrievalTestResult.Error != nil {
		res.TraceRetrievalTestResult.Status = StatusFailed
	}

	return
}
