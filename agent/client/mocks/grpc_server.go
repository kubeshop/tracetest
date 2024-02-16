package mocks

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/avast/retry-go"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/server/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

type GrpcServerMock struct {
	proto.UnimplementedOrchestratorServer
	port                      int
	triggerChannel            chan Message[*proto.TriggerRequest]
	pollingChannel            chan Message[*proto.PollingRequest]
	otlpConnectionTestChannel chan Message[*proto.OTLPConnectionTestRequest]
	terminationChannel        chan Message[*proto.ShutdownRequest]
	dataStoreTestChannel      chan Message[*proto.DataStoreConnectionTestRequest]

	lastTriggerResponse             Message[*proto.TriggerResponse]
	lastPollingResponse             Message[*proto.PollingResponse]
	lastOtlpConnectionResponse      Message[*proto.OTLPConnectionTestResponse]
	lastDataStoreConnectionResponse Message[*proto.DataStoreConnectionTestResponse]

	server *grpc.Server
}

type Message[T any] struct {
	Context context.Context
	Data    T
}

func NewGrpcServer() *GrpcServerMock {
	server := &GrpcServerMock{
		triggerChannel:            make(chan Message[*proto.TriggerRequest]),
		pollingChannel:            make(chan Message[*proto.PollingRequest]),
		terminationChannel:        make(chan Message[*proto.ShutdownRequest]),
		dataStoreTestChannel:      make(chan Message[*proto.DataStoreConnectionTestRequest]),
		otlpConnectionTestChannel: make(chan Message[*proto.OTLPConnectionTestRequest]),
	}
	var wg sync.WaitGroup
	wg.Add(1)

	err := retry.Do(func() error {
		return server.start(&wg, 0)
	}, retry.Attempts(client.ReconnectRetryAttempts), retry.Delay(client.ReconnectRetryAttemptDelay))
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	return server
}

func (s *GrpcServerMock) Addr() string {
	return fmt.Sprintf("localhost:%d", s.port)
}

func (s *GrpcServerMock) start(wg *sync.WaitGroup, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.port = lis.Addr().(*net.TCPAddr).Port

	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
			otelgrpc.WithPropagators(
				propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
			),
		)),
	)
	proto.RegisterOrchestratorServer(server, s)

	s.server = server

	wg.Done()

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal("failed to serve: %w", err)
		}
	}()

	return nil
}

func (s *GrpcServerMock) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.AgentConfiguration, error) {
	return &proto.AgentConfiguration{
		Configuration: &proto.SessionConfiguration{
			BatchTimeout: 1000,
		},
		Identification: &proto.AgentIdentification{
			Token: "token",
		},
	}, nil
}

func (s *GrpcServerMock) RegisterStopRequestAgent(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterStopRequestAgentServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	return nil
}

func (s *GrpcServerMock) RegisterTriggerAgent(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterTriggerAgentServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	for {
		triggerRequest := <-s.triggerChannel
		err := telemetry.InjectContextIntoStream(triggerRequest.Context, stream)
		if err != nil {
			log.Println(err.Error())
		}

		err = stream.Send(triggerRequest.Data)
		if err != nil {
			log.Println("could not send trigger request to agent: %w", err)
		}

	}
}

func (s *GrpcServerMock) SendTriggerResult(ctx context.Context, result *proto.TriggerResponse) (*proto.Empty, error) {
	if result.AgentIdentification == nil || result.AgentIdentification.Token != "token" {
		return nil, fmt.Errorf("could not validate token")
	}

	s.lastTriggerResponse = Message[*proto.TriggerResponse]{Data: result, Context: ctx}
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) RegisterPollerAgent(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterPollerAgentServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	for {
		pollerRequest := <-s.pollingChannel
		err := stream.Send(pollerRequest.Data)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

func (s *GrpcServerMock) RegisterDataStoreConnectionTestAgent(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterDataStoreConnectionTestAgentServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	for {
		dsTestRequest := <-s.dataStoreTestChannel
		err := stream.Send(dsTestRequest.Data)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

func (s *GrpcServerMock) RegisterOTLPConnectionTestListener(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterOTLPConnectionTestListenerServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	for {
		testRequest := <-s.otlpConnectionTestChannel
		err := stream.Send(testRequest.Data)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

func (s *GrpcServerMock) SendOTLPConnectionTestResult(ctx context.Context, result *proto.OTLPConnectionTestResponse) (*proto.Empty, error) {
	if result.AgentIdentification == nil || result.AgentIdentification.Token != "token" {
		return nil, fmt.Errorf("could not validate token")
	}

	s.lastOtlpConnectionResponse = Message[*proto.OTLPConnectionTestResponse]{Data: result, Context: ctx}
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) SendDataStoreConnectionTestResult(ctx context.Context, result *proto.DataStoreConnectionTestResponse) (*proto.Empty, error) {
	if result.AgentIdentification == nil || result.AgentIdentification.Token != "token" {
		return nil, fmt.Errorf("could not validate token")
	}

	s.lastDataStoreConnectionResponse = Message[*proto.DataStoreConnectionTestResponse]{Data: result, Context: ctx}
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) SendPolledSpans(ctx context.Context, result *proto.PollingResponse) (*proto.Empty, error) {
	if result.AgentIdentification == nil || result.AgentIdentification.Token != "token" {
		return nil, fmt.Errorf("could not validate token")
	}

	s.lastPollingResponse = Message[*proto.PollingResponse]{Data: result, Context: ctx}
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) RegisterShutdownListener(_ *proto.AgentIdentification, stream proto.Orchestrator_RegisterShutdownListenerServer) error {
	for {
		shutdownRequest := <-s.terminationChannel
		err := stream.Send(shutdownRequest.Data)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

// Test methods

func (s *GrpcServerMock) SendTriggerRequest(ctx context.Context, request *proto.TriggerRequest) {
	s.triggerChannel <- Message[*proto.TriggerRequest]{Context: ctx, Data: request}
}

func (s *GrpcServerMock) SendPollingRequest(ctx context.Context, request *proto.PollingRequest) {
	s.pollingChannel <- Message[*proto.PollingRequest]{Context: ctx, Data: request}
}

func (s *GrpcServerMock) SendDataStoreConnectionTestRequest(ctx context.Context, request *proto.DataStoreConnectionTestRequest) {
	s.dataStoreTestChannel <- Message[*proto.DataStoreConnectionTestRequest]{Context: ctx, Data: request}
}

func (s *GrpcServerMock) SendOTLPConnectionTestRequest(ctx context.Context, request *proto.OTLPConnectionTestRequest) {
	s.otlpConnectionTestChannel <- Message[*proto.OTLPConnectionTestRequest]{Context: ctx, Data: request}
}

func (s *GrpcServerMock) GetLastTriggerResponse() Message[*proto.TriggerResponse] {
	return s.lastTriggerResponse
}

func (s *GrpcServerMock) GetLastPollingResponse() Message[*proto.PollingResponse] {
	return s.lastPollingResponse
}

func (s *GrpcServerMock) GetLastOTLPConnectionResponse() Message[*proto.OTLPConnectionTestResponse] {
	return s.lastOtlpConnectionResponse
}

func (s *GrpcServerMock) GetLastDataStoreConnectionResponse() Message[*proto.DataStoreConnectionTestResponse] {
	return s.lastDataStoreConnectionResponse
}

func (s *GrpcServerMock) TerminateConnection(ctx context.Context, reason string) {
	s.terminationChannel <- Message[*proto.ShutdownRequest]{
		Context: ctx,
		Data:    &proto.ShutdownRequest{Reason: reason},
	}
}

func (s *GrpcServerMock) Stop() {
	s.server.Stop()
}
