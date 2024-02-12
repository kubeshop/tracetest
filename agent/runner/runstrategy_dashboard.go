package runner

import (
	"context"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/ui/dashboard"
)

func (s *Runner) RunDashboardStrategy(ctx context.Context, cfg agentConfig.Config, uiEndpoint string) error {
	session, _, err := s.authenticate(ctx, cfg)
	if err != nil {
		return err
	}

	defer session.Close()

	return dashboard.StartDashboard(ctx)
}
