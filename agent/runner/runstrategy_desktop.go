package runner

import (
	"context"
	"errors"
	"fmt"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/cli/config"

	consoleUI "github.com/kubeshop/tracetest/agent/ui"
)

func (s *Runner) RunDesktopStrategy(ctx context.Context, cfg agentConfig.Config, uiEndpoint string) error {
	s.ui.Infof("Starting Agent with name %s...", cfg.Name)

	isStarted := false
	session := &Session{}

	var err error

	for !isStarted {
		session, err = StartSession(ctx, cfg, nil, s.logger)
		if err != nil && errors.Is(err, ErrOtlpServerStart) {
			s.ui.Error("Tracetest Agent binds to the OpenTelemetry ports 4317 and 4318 which are used to receive trace information from your system. The agent tried to bind to these ports, but failed.")
			shouldRetry := s.ui.Enter("Please stop the process currently listening on these ports and press enter to try again.")

			if !shouldRetry {
				s.ui.Finish()
				return err
			}

			continue
		}

		if err != nil {
			return err
		}

		isStarted = true
	}

	claims, err := config.GetTokenClaims(session.Token)
	if err != nil {
		return err
	}

	isOpen := true
	message := `Agent is started! Leave the terminal open so tests can be run and traces gathered from this environment.
You can`
	options := []consoleUI.Option{{
		Text: "Open Tracetest in a browser to this environment",
		Fn: func(_ consoleUI.ConsoleUI) {
			s.ui.OpenBrowser(fmt.Sprintf("%sorganizations/%s/environments/%s/dashboard", uiEndpoint, claims["organization_id"], claims["environment_id"]))
		},
	}, {
		Text: "Stop this agent",
		Fn: func(_ consoleUI.ConsoleUI) {
			isOpen = false
			session.Close()
			s.ui.Finish()
		},
	}}

	for isOpen {
		selected := s.ui.Select(message, options, 0)
		selected.Fn(s.ui)
	}
	return nil
}
