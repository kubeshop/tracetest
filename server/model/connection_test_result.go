package model

type ConnectionResult struct {
	PortCheck      ConnectionTestStep
	Connectivity   ConnectionTestStep
	Authentication ConnectionTestStep
	FetchTraces    ConnectionTestStep
}

func (c ConnectionResult) HasSucceed() bool {
	return c.Connectivity.HasSucceed() && c.Authentication.HasSucceed() && c.FetchTraces.HasSucceed()
}

type ConnectionTestStep struct {
	Passed  bool
	Status  Status
	Message string
	Error   error
}

func (r *ConnectionTestStep) HasSucceed() bool {
	if r == nil {
		return true
	}

	return r.Error == nil
}

func (r *ConnectionTestStep) IsSet() bool {
	if r == nil {
		return false
	}

	return r.Message != ""
}
