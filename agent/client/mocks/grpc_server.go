package mocks

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/kubeshop/tracetest/agent/proto"
	"google.golang.org/grpc"
)

type GrpcServerMock struct {
	proto.UnimplementedOrchestratorServer
	port                 int
	triggerChannel       chan *proto.TriggerRequest
	pollingChannel       chan *proto.PollingRequest
	terminationChannel   chan *proto.ShutdownRequest
	dataStoreTestChannel chan *proto.DataStoreConnectionTestRequest

	lastTriggerResponse *proto.TriggerResponse
	lastPollingResponse *proto.PollingResponse

	server *grpc.Server
}

func NewGrpcServer() *GrpcServerMock {
	server := &GrpcServerMock{
		triggerChannel:       make(chan *proto.TriggerRequest),
		pollingChannel:       make(chan *proto.PollingRequest),
		terminationChannel:   make(chan *proto.ShutdownRequest),
		dataStoreTestChannel: make(chan *proto.DataStoreConnectionTestRequest),
	}
	var wg sync.WaitGroup
	wg.Add(1)

	err := retry.Do(func() error {
		return server.start(&wg, 0)
	}, retry.Attempts(3), retry.Delay(time.Second))
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

	server := grpc.NewServer()
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

func (s *GrpcServerMock) RegisterTriggerAgent(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterTriggerAgentServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	for {
		triggerRequest := <-s.triggerChannel
		err := stream.Send(triggerRequest)
		if err != nil {
			log.Println("could not send trigger request to agent: %w", err)
		}

	}
}

func (s *GrpcServerMock) SendTriggerResult(ctx context.Context, result *proto.TriggerResponse) (*proto.Empty, error) {
	if result.AgentIdentification == nil || result.AgentIdentification.Token != "token" {
		return nil, fmt.Errorf("could not validate token")
	}

	s.lastTriggerResponse = result
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) RegisterPollerAgent(id *proto.AgentIdentification, stream proto.Orchestrator_RegisterPollerAgentServer) error {
	if id.Token != "token" {
		return fmt.Errorf("could not validate token")
	}

	for {
		pollerRequest := <-s.pollingChannel
		err := stream.Send(pollerRequest)
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
		err := stream.Send(dsTestRequest)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

func (s *GrpcServerMock) SendPolledSpans(ctx context.Context, result *proto.PollingResponse) (*proto.Empty, error) {
	if result.AgentIdentification == nil || result.AgentIdentification.Token != "token" {
		return nil, fmt.Errorf("could not validate token")
	}

	s.lastPollingResponse = result
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) RegisterShutdownListener(_ *proto.AgentIdentification, stream proto.Orchestrator_RegisterShutdownListenerServer) error {
	for {
		shutdownRequest := <-s.terminationChannel
		err := stream.Send(shutdownRequest)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

// Test methods

func (s *GrpcServerMock) SendTriggerRequest(request *proto.TriggerRequest) {
	s.triggerChannel <- request
}

func (s *GrpcServerMock) SendPollingRequest(request *proto.PollingRequest) {
	s.pollingChannel <- request
}

func (s *GrpcServerMock) GetLastTriggerResponse() *proto.TriggerResponse {
	return s.lastTriggerResponse
}

func (s *GrpcServerMock) GetLastPollingResponse() *proto.PollingResponse {
	return s.lastPollingResponse
}

func (s *GrpcServerMock) TerminateConnection(reason string) {
	s.terminationChannel <- &proto.ShutdownRequest{
		Reason: reason,
	}
}

func (s *GrpcServerMock) Stop() {
	s.server.Stop()
}
