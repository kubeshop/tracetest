package misc_parameters

type TestRun struct {
	RunTestFileDefinition,
	RunTestEnvID,
	RunTestJUnit string
	RunTestWaitForResult bool
}

func NewTestRun() *TestRun {
	return &TestRun{}
}
