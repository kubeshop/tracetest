package runner

import (
	"context"
	"fmt"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"

	consoleUI "github.com/kubeshop/tracetest/agent/ui"
)

func (s *Runner) RunDesktopStrategy(ctx context.Context, cfg agentConfig.Config, uiEndpoint string) error {
	s.ui.Infof("Starting Agent with name %s...", cfg.Name)

	sensor := sensors.NewSensor()
	dashboardObserver := newDashboardObserver(sensor)
	session, claims, err := s.authenticate(ctx, cfg, dashboardObserver)
	if err != nil {
		return err
	}

	isOpen := true
	message := `Agent is started! Leave the terminal open so tests can be run and traces gathered from this environment.
You can`
	options := []consoleUI.Option{
		{
			Text: "Open Tracetest in a browser to this environment",
			Fn: func(_ consoleUI.ConsoleUI) {
				s.ui.OpenBrowser(fmt.Sprintf("%sorganizations/%s/environments/%s", uiEndpoint, claims["organization_id"], claims["environment_id"]))
			},
		},
		{
			Text: "(Experimental) Open Dashboard",
			Fn: func(ui consoleUI.ConsoleUI) {
				sensor.Reset()
				err := s.RunDashboardStrategy(ctx, cfg, uiEndpoint, sensor)
				if err != nil {
					fmt.Println(err.Error())
				}
			},
		},
		{
			Text: "Stop this agent",
			Fn: func(_ consoleUI.ConsoleUI) {
				isOpen = false
				session.Close()
				s.claims = nil
				s.ui.Finish()
			},
		}}

	for isOpen {
		selected := s.ui.Select(message, options, 0)
		selected.Fn(s.ui)
	}
	return nil
}
