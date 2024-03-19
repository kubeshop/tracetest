package otlp

import (
	"context"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type RequestType string

var (
	RequestTypeHTTP RequestType = "HTTP"
	RequestTypeGRPC RequestType = "gRPC"
)

type Ingester interface {
	Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType RequestType) (*pb.ExportTraceServiceResponse, error)
}
