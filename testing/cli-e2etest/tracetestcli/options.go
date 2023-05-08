package tracetestcli

type ExecOption func(*executionState)

type executionState struct {
	cliConfigFile string
}

func composeExecutionState(options ...ExecOption) *executionState {
	state := &executionState{}

	for _, option := range options {
		option(state)
	}

	return state
}

func WithCLIConfig(cliConfig string) ExecOption {
	return func(es *executionState) {
		es.cliConfigFile = cliConfig
	}
}
