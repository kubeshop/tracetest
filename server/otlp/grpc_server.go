package otlp

import (
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedTraceServiceServer

	addr     string
	ingester Ingester

	gServer *grpc.Server
	tracer  trace.Tracer
	logger  *zap.Logger
}

func NewGrpcServer(addr string, ingester Ingester, tracer trace.Tracer) *grpcServer {
	return &grpcServer{
		addr:     addr,
		ingester: ingester,
		tracer:   tracer,
		logger:   zap.NewNop(),
	}
}

func (s *grpcServer) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

func (s *grpcServer) Start() error {
	size := 1024 * 1024 * 50
	s.gServer = grpc.NewServer(
		grpc.MaxSendMsgSize(size),
		grpc.MaxRecvMsgSize(size),
	)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("cannot listen on address %s: %w", s.addr, err)
	}
	pb.RegisterTraceServiceServer(s.gServer, s)
	go s.gServer.Serve(listener)
	return nil
}

func (s *grpcServer) Stop() {
	s.gServer.Stop()
}

func (s grpcServer) Export(ctx context.Context, request *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceResponse, error) {
	ctx, span := s.tracer.Start(ctx, "Export trace")
	defer span.End()

	s.logger.Debug("Received ExportTraceServiceRequest", zap.Any("request", request))

	response, err := s.ingester.Ingest(ctx, request, RequestTypeGRPC)

	s.logger.Debug("Sending ExportTraceServiceResponse", zap.Any("response", response), zap.Error(err))

	return response, err
}
