package testfixtures

import (
	"github.com/kubeshop/tracetest/testmock"
)

func getPokeshopApp(args ...interface{}) (*testmock.DemoApp, error) {
	return testmock.GetDemoApplicationInstance()
}

func init() {
	RegisterFixture(POKESHOP_APP, getPokeshopApp)
}
