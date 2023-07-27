package collector

import (
	"context"
	"fmt"
	"net"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedTraceServiceServer

	addr     string
	ingester ingester

	gServer *grpc.Server
}

func newGrpcServer(addr string, ingester ingester) *grpcServer {
	return &grpcServer{
		addr:     addr,
		ingester: ingester,
	}
}

func (s *grpcServer) Start() error {
	s.gServer = grpc.NewServer()
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("cannot listen on address %s: %w", s.addr, err)
	}
	pb.RegisterTraceServiceServer(s.gServer, s)
	return s.gServer.Serve(listener)
}

func (s *grpcServer) Stop() {
	s.gServer.Stop()
}

func (s grpcServer) Export(ctx context.Context, request *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceResponse, error) {
	return s.ingester.Ingest(ctx, request, "gRPC")
}
