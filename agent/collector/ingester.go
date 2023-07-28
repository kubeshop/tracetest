package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/server/otlp"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ingester interface {
	Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType otlp.RequestType) (*pb.ExportTraceServiceResponse, error)
	Stop()
}

func newForwardIngester(ctx context.Context, batchTimeout time.Duration, remoteIngesterConfig remoteIngesterConfig) (ingester, error) {
	ingester := &forwardIngester{
		BatchTimeout:   batchTimeout,
		RemoteIngester: remoteIngesterConfig,
		buffer:         &buffer{},
		done:           make(chan bool),
	}

	err := ingester.connectToRemoteServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not connect to remote server: %w", err)
	}

	go ingester.startBatchWorker()

	return ingester, nil
}

// forwardIngester forwards all incoming spans to a remote ingester. It also batches those
// spans to reduce network traffic.
type forwardIngester struct {
	BatchTimeout   time.Duration
	RemoteIngester remoteIngesterConfig
	client         pb.TraceServiceClient
	buffer         *buffer
	done           chan bool
}

type remoteIngesterConfig struct {
	URL   string
	Token string
}

type buffer struct {
	mutex sync.Mutex
	spans []*v1.ResourceSpans
}

func (i *forwardIngester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType otlp.RequestType) (*pb.ExportTraceServiceResponse, error) {
	i.buffer.mutex.Lock()
	i.buffer.spans = append(i.buffer.spans, request.ResourceSpans...)
	i.buffer.mutex.Unlock()

	return &pb.ExportTraceServiceResponse{
		PartialSuccess: &pb.ExportTracePartialSuccess{
			RejectedSpans: 0,
		},
	}, nil
}

func (i *forwardIngester) connectToRemoteServer(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, i.RemoteIngester.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("could not connect to remote server: %w", err)
	}

	i.client = pb.NewTraceServiceClient(conn)
	return nil
}

func (i *forwardIngester) startBatchWorker() {
	ticker := time.NewTicker(i.BatchTimeout)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := i.executeBatch(context.Background())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (i *forwardIngester) executeBatch(ctx context.Context) error {
	i.buffer.mutex.Lock()
	newSpans := i.buffer.spans
	i.buffer.spans = []*v1.ResourceSpans{}
	i.buffer.mutex.Unlock()

	if len(newSpans) == 0 {
		return nil
	}

	err := i.forwardSpans(ctx, newSpans)
	if err != nil {
		return err
	}

	return nil
}

func (i *forwardIngester) forwardSpans(ctx context.Context, spans []*v1.ResourceSpans) error {
	_, err := i.client.Export(ctx, &pb.ExportTraceServiceRequest{
		ResourceSpans: spans,
	})

	if err != nil {
		return fmt.Errorf("could not forward spans to remote server: %w", err)
	}

	return nil
}

func (i *forwardIngester) Stop() {
	i.done <- true
}
