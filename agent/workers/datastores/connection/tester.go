package connection

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

type TestStep interface {
	TestConnection(ctx context.Context) model.ConnectionTestStep
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

func (t *Tester) TestConnection(ctx context.Context) (res model.ConnectionResult) {
	if t.portLinterStep != nil {
		res.PortCheck = t.portLinterStep.TestConnection(ctx)
		if res.PortCheck.Error != nil {
			res.PortCheck.Status = model.StatusFailed
			return
		}
	}

	res.Connectivity = t.connectivityTestStep.TestConnection(ctx)
	if res.Connectivity.Error != nil {
		res.Connectivity.Status = model.StatusFailed
		return
	}

	res.Authentication = t.authenticationTestStep.TestConnection(ctx)
	if res.Authentication.Error != nil {
		res.Authentication.Status = model.StatusFailed
		return
	}

	res.FetchTraces = t.pollingTestStep.TestConnection(ctx)
	if res.FetchTraces.Error != nil {
		res.FetchTraces.Status = model.StatusFailed
	}

	return
}
