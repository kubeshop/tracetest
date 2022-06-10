package testfixtures

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/testmock"
)

type TracetestAppFixtureConfig struct {
	Prefix string
	Port   int
}

type TracetestAppFixtureOption func(config *TracetestAppFixtureConfig)

func WithServerPrefix(prefix string) TracetestAppFixtureOption {
	return func(config *TracetestAppFixtureConfig) {
		config.Prefix = prefix
	}
}

func WithServerPort(port int) TracetestAppFixtureOption {
	return func(config *TracetestAppFixtureConfig) {
		config.Port = port
	}
}

func GetTracetestApp(options ...TracetestAppFixtureOption) (*app.App, error) {
	config := TracetestAppFixtureConfig{}
	for _, option := range options {
		option(&config)
	}

	fixtureOptions := make([]Option, 0)
	fixtureOptions = append(fixtureOptions, WithArguments(config))
	if len(options) > 0 {
		// Disable cache to prevent problems with different configurations
		fixtureOptions = append(fixtureOptions, WithCacheDisabled())
	}

	return GetFixtureValue[*app.App](TRACETEST_APP, fixtureOptions...)
}

func init() {
	RegisterFixture(TRACETEST_APP, getTracetestApp)
}

func getTracetestApp(options FixtureOptions) (*app.App, error) {
	demoApp, err := GetFixtureValue[*testmock.DemoApp](POKESHOP_APP)
	if err != nil {
		return nil, fmt.Errorf("could not get pokeshop app: %w", err)
	}

	appOptions := make([]testmock.TestingAppOption, 0)
	arguments := options.Arguments.(TracetestAppFixtureConfig)
	if arguments.Prefix != "" {
		appOptions = append(appOptions, testmock.WithServerPrefix(arguments.Prefix))
	}

	if arguments.Port != 0 {
		appOptions = append(appOptions, testmock.WithServerPort(arguments.Port))
	}

	tracetestApp, err := testmock.GetTestingApp(demoApp, appOptions...)
	if err != nil {
		return nil, err
	}

	go tracetestApp.Start()

	time.Sleep(1 * time.Second)

	return tracetestApp, nil
}
