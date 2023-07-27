package collector

import (
	"context"
	"time"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type requestType string

var (
	RequestTypeHTTP requestType = "HTTP"
	RequestTypeGRPC requestType = "gRPC"
)

type ingester interface {
	Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType requestType) (*pb.ExportTraceServiceResponse, error)
}

// forwardIngester forwards all incoming spans to a remote ingester. It also batches those
// spans to reduce network traffic.
type forwardIngester struct {
	BatchTimeout   time.Duration
	RemoteIngester remoteIngesterConfig
}

type remoteIngesterConfig struct {
	URL   string
	Token string
}

func (i *forwardIngester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType requestType) (*pb.ExportTraceServiceResponse, error) {
	return nil, nil
}
