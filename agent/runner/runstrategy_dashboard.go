package runner

import (
	"context"

	"github.com/kubeshop/tracetest/agent/collector"
	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/ui/dashboard"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/models"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/kubeshop/tracetest/server/version"
)

func (s *Runner) RunDashboardStrategy(ctx context.Context, cfg agentConfig.Config, uiEndpoint string) error {
	session, claims, err := s.authenticate(ctx, cfg)
	if err != nil {
		return err
	}

	defer session.Close()

	sensor := sensors.NewSensor()
	if collector := collector.GetActiveCollector(); collector != nil {
		collector.SetSensor(sensor)
	}

	// TODO: convert ids into names
	return dashboard.StartDashboard(ctx, models.EnvironmentInformation{
		OrganizationName: claims["organization_id"].(string),
		EnvironmentName:  claims["environment_id"].(string),
		AgentVersion:     version.Version,
	}, sensor)
}
