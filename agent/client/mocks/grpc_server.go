package mocks

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/kubeshop/tracetest/agent/proto"
	"google.golang.org/grpc"
)

type GrpcServerMock struct {
	proto.UnimplementedOrchestratorServer
	port           int
	triggerChannel chan *proto.TriggerRequest
	pollingChannel chan *proto.PollingRequest

	lastTriggerResponse *proto.TriggerResponse
	lastPollingResponse *proto.PollingResponse
}

func NewGrpcServer() *GrpcServerMock {
	server := &GrpcServerMock{
		triggerChannel: make(chan *proto.TriggerRequest),
		pollingChannel: make(chan *proto.PollingRequest),
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go server.start(&wg)

	wg.Wait()

	return server
}

func (s *GrpcServerMock) Addr() string {
	return fmt.Sprintf("localhost:%d", s.port)
}

func (s *GrpcServerMock) start(wg *sync.WaitGroup) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.port = lis.Addr().(*net.TCPAddr).Port

	server := grpc.NewServer()
	proto.RegisterOrchestratorServer(server, s)

	wg.Done()
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *GrpcServerMock) Connect(ctx context.Context, req *proto.ConnectRequest) (*proto.ConnectResponse, error) {
	return &proto.ConnectResponse{Configuration: &proto.SessionConfiguration{
		BatchTimeout: 1000,
	}}, nil
}

func (s *GrpcServerMock) RegisterTriggerAgent(_ *proto.ConnectRequest, stream proto.Orchestrator_RegisterTriggerAgentServer) error {
	for {
		triggerRequest := <-s.triggerChannel
		err := stream.Send(triggerRequest)
		if err != nil {
			log.Println("could not send trigger request to agent: %w", err)
		}
	}
}

func (s *GrpcServerMock) SendTriggerResult(ctx context.Context, result *proto.TriggerResponse) (*proto.Empty, error) {
	s.lastTriggerResponse = result
	return &proto.Empty{}, nil
}

func (s *GrpcServerMock) RegisterPollerAgent(_ *proto.ConnectRequest, stream proto.Orchestrator_RegisterPollerAgentServer) error {
	for {
		pollerRequest := <-s.pollingChannel
		err := stream.Send(pollerRequest)
		if err != nil {
			log.Println("could not send polling request to agent: %w", err)
		}
	}
}

func (s *GrpcServerMock) SendPolledSpans(ctx context.Context, result *proto.PollingResponse) (*proto.Empty, error) {
	s.lastPollingResponse = result
	return &proto.Empty{}, nil
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

func (s *GrpcServerMock) GetLastSpans() *proto.PollingResponse {
	return s.lastPollingResponse
}
