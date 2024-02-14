package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/agent/collector"
	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/ui/dashboard"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/models"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/kubeshop/tracetest/server/version"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"go.uber.org/zap"
)

func (s *Runner) RunDashboardStrategy(ctx context.Context, cfg agentConfig.Config, uiEndpoint string, sensor sensors.Sensor) error {
	// This prevents the agent logger from printing lots of messages
	// and override the dashboard UI.
	// By calling enableLogger() at the end of this function, the logger behavior is restored
	enableLogger := s.disableLogger()
	defer enableLogger()

	if collector := collector.GetActiveCollector(); collector != nil {
		collector.SetSensor(sensor)
	}

	claims := s.getCurrentSessionClaims()
	if claims == nil {
		return fmt.Errorf("not authenticated")
	}

	// TODO: convert ids into names
	return dashboard.StartDashboard(ctx, models.EnvironmentInformation{
		OrganizationID: claims["organization_id"].(string),
		EnvironmentID:  claims["environment_id"].(string),
		AgentVersion:   version.Version,
		ServerEndpoint: uiEndpoint,
	}, sensor)
}

func (s *Runner) disableLogger() func() {
	if s.loggerLevel == nil {
		// logger is not active, so it's safe to do nothing
		return func() {}
	}

	oldLevel := s.loggerLevel.Level()
	s.loggerLevel.SetLevel(zap.PanicLevel)

	return func() {
		s.loggerLevel.SetLevel(oldLevel)
	}
}

type dashboardObserver struct {
	runs   map[string]models.TestRun
	sensor sensors.Sensor
}

func (o *dashboardObserver) EndDataStoreConnection(*proto.DataStoreConnectionTestRequest, error) {

}

func (o *dashboardObserver) EndSpanReceive([]*v1.Span, error) {

}

func (o *dashboardObserver) EndStopRequest(*proto.StopRequest, error) {

}

func (o *dashboardObserver) EndTracePoll(*proto.PollingRequest, error) {

}

func (o *dashboardObserver) EndTriggerExecution(*proto.TriggerRequest, error) {

}

func (o *dashboardObserver) Error(error) {
}

func (o *dashboardObserver) StartDataStoreConnection(*proto.DataStoreConnectionTestRequest) {
}

func (o *dashboardObserver) StartSpanReceive([]*v1.Span) {
}

func (o *dashboardObserver) StartStopRequest(request *proto.StopRequest) {
	model := o.getRun(request.TestID, request.RunID)
	model.Status = "Stopped by user"

	o.setRun(model)
	o.sensor.Emit(events.UpdatedTestRun, model)
}

func (o *dashboardObserver) StartTracePoll(request *proto.PollingRequest) {
	model := o.getRun(request.TestID, request.RunID)
	model.Status = "Awaiting Trace"

	o.setRun(model)
	o.sensor.Emit(events.UpdatedTestRun, model)
}

func (o *dashboardObserver) StartTriggerExecution(request *proto.TriggerRequest) {
	model := o.getRun(request.TestID, request.RunID)
	model.TestID = request.TestID
	model.RunID = fmt.Sprintf("%d", request.RunID)
	model.Type = request.Trigger.Type
	model.Endpoint = getEndpoint(request)
	model.Name = "<not set>"
	model.Status = "Triggering"
	model.Started = time.Now()

	o.setRun(model)
	o.sensor.Emit(events.NewTestRun, model)
}

func (o *dashboardObserver) getRun(testID string, runID int32) models.TestRun {
	if model, ok := o.runs[fmt.Sprintf("%s-%d", testID, runID)]; ok {
		return model
	}

	return models.TestRun{TestID: testID, RunID: fmt.Sprintf("%d", runID)}
}

func (o *dashboardObserver) setRun(model models.TestRun) {
	o.runs[fmt.Sprintf("%s-%s", model.TestID, model.RunID)] = model
}

func getEndpoint(request *proto.TriggerRequest) string {
	switch request.Trigger.Type {
	case "http":
		return request.Trigger.Http.Url
	case "grpc":
		return fmt.Sprintf("%s/%s", request.Trigger.Grpc.Address, request.Trigger.Grpc.Service)
	case "kafka":
		return request.Trigger.Kafka.Topic
	case "traceID":
		return request.Trigger.TraceID.Id
	default:
		return "<not set>"
	}
}

func newDashboardObserver(sensor sensors.Sensor) event.Observer {
	return &dashboardObserver{sensor: sensor, runs: make(map[string]models.TestRun)}
}
