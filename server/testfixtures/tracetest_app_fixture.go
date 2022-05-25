package testfixtures

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/app"
	"github.com/kubeshop/tracetest/testmock"
)

func getTracetestApp(args ...interface{}) (*app.App, error) {
	demoApp, err := GetFixtureValue[*testmock.DemoApp](POKESHOP_APP)
	if err != nil {
		return nil, fmt.Errorf("could not get pokeshop app: %w", err)
	}

	tracetestApp, err := testmock.GetTestingApp(demoApp)
	if err != nil {
		return nil, err
	}

	go tracetestApp.Start()

	time.Sleep(1 * time.Second)

	return tracetestApp, nil
}

var _ Generator[*app.App] = getTracetestApp

func init() {
	RegisterFixture(TRACETEST_APP, getTracetestApp)
}
