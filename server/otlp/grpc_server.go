package otlp

import (
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedTraceServiceServer

	addr     string
	ingester Ingester

	gServer *grpc.Server
	tracer  trace.Tracer
}

func NewGrpcServer(addr string, ingester Ingester, tracer trace.Tracer) *grpcServer {
	return &grpcServer{
		addr:     addr,
		ingester: ingester,
		tracer:   tracer,
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
	ctx, span := s.tracer.Start(ctx, "Export trace")
	defer span.End()

	return s.ingester.Ingest(ctx, request, "gRPC")
}
