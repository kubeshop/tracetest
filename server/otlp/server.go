package otlp

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedTraceServiceServer
	db model.Repository
}

func StartServer(port int, db model.Repository) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	server := Server{
		db: db,
	}

	s := grpc.NewServer()
	pb.RegisterTraceServiceServer(s, &server)
	if err := s.Serve(listener); err != nil {
		return fmt.Errorf("could not serve: %w", err)
	}

	return nil
}

func (s Server) Export(ctx context.Context, request *pb.ExportTraceServiceRequest) (*pb.ExportTraceServiceResponse, error) {
	if len(request.ResourceSpans) == 0 {
		return &pb.ExportTraceServiceResponse{}, nil
	}

	spansByTrace := s.getSpansByTrace(request)

	fmt.Println(spansByTrace)

	for traceID, spans := range spansByTrace {
		s.saveSpansIntoTest(ctx, traceID, spans)
	}

	return &pb.ExportTraceServiceResponse{}, nil
}

func (s Server) getSpansByTrace(request *pb.ExportTraceServiceRequest) map[trace.TraceID][]traces.Span {
	otelSpans := make([]*v1.Span, 0)
	for _, resourceSpan := range request.ResourceSpans {
		for _, spans := range resourceSpan.ScopeSpans {
			otelSpans = append(otelSpans, spans.Spans...)
		}
	}

	spansByTrace := make(map[trace.TraceID][]traces.Span)

	for _, span := range otelSpans {
		traceID := traces.CreateTraceID(span.TraceId)
		var existingArray []traces.Span
		if spansArray, ok := spansByTrace[traceID]; !ok {
			existingArray = spansArray
		} else {
			existingArray = make([]traces.Span, 0)
		}

		existingArray = append(existingArray, *traces.ConvertOtelSpanIntoSpan(span))
		spansByTrace[traceID] = existingArray
	}

	return spansByTrace
}

func (s Server) saveSpansIntoTest(ctx context.Context, traceID trace.TraceID, spans []traces.Span) error {
	run, err := s.db.GetRunByTraceID(ctx, traceID)
	if strings.Contains(err.Error(), "record not found") {
		// span is not part of any known test run. So it will be ignored
		return nil
	}

	if err != nil {
		return fmt.Errorf("could not find test run with traceID %s: %w", traceID.String(), err)
	}

	existingSpans := run.Trace.Spans()
	newSpans := append(existingSpans, spans...)
	newTrace := traces.New(traceID.String(), newSpans)

	run.Trace = &newTrace

	err = s.db.UpdateRun(ctx, run)
	if err != nil {
		return fmt.Errorf("could not update run: %w", err)
	}

	return nil
}
