package misc_parameters

type TestExport struct {
	ExportTestId         string
	ExportTestOutputFile string
	Version              int32
}

func NewTestExport() *TestExport {
	return &TestExport{}
}
