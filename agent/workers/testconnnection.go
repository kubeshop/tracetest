package workers

import (
	"context"
	"fmt"
	"log"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type TestConnectionWorker struct {
	client *client.Client
	tracer trace.Tracer
}

func NewTestConnectionWorker(client *client.Client) *TestConnectionWorker {
	// TODO: use a real tracer
	tracer := trace.NewNoopTracerProvider().Tracer("noop")

	return &TestConnectionWorker{
		client: client,
		tracer: tracer,
	}
}

func (w *TestConnectionWorker) Test(ctx context.Context, request *proto.DataStoreConnectionTestRequest) error {
	fmt.Println("Data Store Test Connection handled by agent")
	datastoreConfig, err := convertProtoToDataStore(request.Datastore)
	if err != nil {
		return err
	}

	if datastoreConfig == nil {
		return fmt.Errorf("invalid datastore: nil")
	}

	dsFactory := tracedb.Factory(nil)
	ds, err := dsFactory(*datastoreConfig)
	if err != nil {
		log.Printf("Invalid datastore: %s", err.Error())
		return err
	}

	if testableTraceDB, ok := ds.(tracedb.TestableTraceDB); ok {
		connectionResult := testableTraceDB.TestConnection(ctx)
		success, steps := convertConnectionResultToProto(connectionResult)

		w.client.SendDataStoreConnectionResult(ctx, &proto.DataStoreConnectionTestResponse{
			RequestID:  request.RequestID,
			Successful: success,
			Steps:      steps,
		})
	} else {
		w.client.SendDataStoreConnectionResult(ctx, &proto.DataStoreConnectionTestResponse{
			RequestID:  request.RequestID,
			Successful: false,
			Steps:      nil,
		})
	}

	return nil
}

func convertConnectionResultToProto(connectionResult model.ConnectionResult) (bool, *proto.DataStoreConnectionTestSteps) {
	steps := &proto.DataStoreConnectionTestSteps{
		PortCheck:      convertConnectionResultStepToProto(connectionResult.PortCheck),
		Connectivity:   convertConnectionResultStepToProto(connectionResult.Connectivity),
		Authentication: convertConnectionResultStepToProto(connectionResult.Authentication),
		FetchTraces:    convertConnectionResultStepToProto(connectionResult.FetchTraces),
	}

	return connectionResult.HasSucceed(), steps
}

func convertConnectionResultStepToProto(step model.ConnectionTestStep) *proto.DataStoreConnectionTestStep {
	errorMsg := ""
	if step.Error != nil {
		errorMsg = step.Error.Error()
	}
	return &proto.DataStoreConnectionTestStep{
		Passed:  step.Passed,
		Status:  string(step.Status),
		Message: step.Message,
		Error:   errorMsg,
	}
}
