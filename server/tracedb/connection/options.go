package connection

import "context"

func WithConnectivityTest(step TestStep) TesterOption {
	return func(t *Tester) {
		t.connectivityTestStep = step
	}
}

func WithAuthenticationTest(step TestStep) TesterOption {
	return func(t *Tester) {
		t.authenticationTestStep = step
	}
}

func WithPollingTest(step TestStep) TesterOption {
	return func(t *Tester) {
		t.pollingTestStep = step
	}
}

type functionTestStep struct {
	fn func(ctx context.Context) (string, error)
}

func (s *functionTestStep) TestConnection(ctx context.Context) ConnectionTestStepResult {
	str, err := s.fn(ctx)
	return ConnectionTestStepResult{
		OperationDescription: str,
		Error:                err,
	}
}

func NewTestStep(f func(ctx context.Context) (string, error)) TestStep {
	return &functionTestStep{fn: f}
}
