package otlp

import (
	"context"
	"fmt"
	"net"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	pb.UnimplementedTraceServiceServer

	addr     string
	exporter IExporter

	gServer *grpc.Server
}

func NewGrpcServer(addr string, exporter IExporter) *GrpcServer {
	return &GrpcServer{
		addr:     addr,
		exporter: exporter,
	}
}

func (s *GrpcServer) Start() error {
	s.gServer = grpc.NewServer()
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("cannot listen on address %s: %w", s.addr, err)
	}
	pb.RegisterTraceServiceServer(s.gServer, s)
	return s.gServer.Serve(listener)
}

func (s *GrpcServer) Stop() {
	s.gServer.Stop()
}

func (s GrpcServer) Export(ctx context.Context, request *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceResponse, error) {
	return s.exporter.Export(ctx, request)
}
