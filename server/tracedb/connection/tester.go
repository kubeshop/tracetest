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

func (t *Tester) TestConnection(ctx context.Context) (res ConnectionTestResult) {
	res.ConnectivityTestResult = t.connectivityTestStep.TestConnection(ctx)
	if res.ConnectivityTestResult.Error != nil {
		return
	}

	res.AuthenticationTestResult = t.authenticationTestStep.TestConnection(ctx)
	if res.ConnectivityTestResult.Error != nil {
		return
	}

	res.TraceRetrievalTestResult = t.pollingTestStep.TestConnection(ctx)

	return
}
