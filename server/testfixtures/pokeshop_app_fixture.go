package testfixtures

import (
	"github.com/kubeshop/tracetest/server/testmock"
)

func init() {
	RegisterFixture(POKESHOP_APP, getPokeshopApp)
}

func GetPokeshopApp() (*testmock.DemoApp, error) {
	return GetFixtureValue[*testmock.DemoApp](POKESHOP_APP)
}

func getPokeshopApp(options FixtureOptions) (*testmock.DemoApp, error) {
	return testmock.GetDemoApplicationInstance()
}
