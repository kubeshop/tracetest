package mocks

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
)

type OTLPIngestionServer struct {
	otlpServer *grpcServer
}

func NewOTLPIngestionServer() (*OTLPIngestionServer, error) {
	grpcServer := newGrpcServer()
	grpcServer.Start()

	return &OTLPIngestionServer{
		otlpServer: grpcServer,
	}, nil
}

func (s *OTLPIngestionServer) Addr() string {
	return fmt.Sprintf("localhost:%d", s.otlpServer.port)
}

func (s *OTLPIngestionServer) ReceivedSpans() []*v1.Span {
	return s.otlpServer.receivedSpans
}

type grpcServer struct {
	pb.UnimplementedTraceServiceServer
	port int

	gServer       *grpc.Server
	receivedSpans []*v1.Span
}

func newGrpcServer() *grpcServer {
	return &grpcServer{}
}

func (s *grpcServer) Start() error {
	s.gServer = grpc.NewServer()
	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		return fmt.Errorf("cannot listen on random tcp port: %w", err)
	}
	s.port = listener.Addr().(*net.TCPAddr).Port

	pb.RegisterTraceServiceServer(s.gServer, s)
	go func() {
		log.Fatal(s.gServer.Serve(listener))
	}()
	return nil
}

func (s *grpcServer) Stop() {
	s.gServer.Stop()
}

func (s *grpcServer) Export(ctx context.Context, request *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceResponse, error) {
	spans := make([]*v1.Span, 0)
	for _, resourceSpan := range request.ResourceSpans {
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			spans = append(spans, scopeSpan.Spans...)
		}
	}

	s.receivedSpans = append(s.receivedSpans, spans...)
	return &pb.ExportTraceServiceResponse{
		PartialSuccess: &pb.ExportTracePartialSuccess{
			RejectedSpans: 0,
		},
	}, nil
}
