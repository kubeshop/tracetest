package connection

import "context"

type TestStep interface {
	TestConnection(ctx context.Context) ConnectionTestStepResult
}

type TesterOption func(*Tester)

type Tester struct {
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

func (t *Tester) TestConnection(ctx context.Context) ConnectionTestResult {
	return ConnectionTestResult{
		ConnectivityTestResult:   t.connectivityTestStep.TestConnection(ctx),
		AuthenticationTestResult: t.authenticationTestStep.TestConnection(ctx),
		TraceRetrievalTestResult: t.pollingTestStep.TestConnection(ctx),
	}
}
