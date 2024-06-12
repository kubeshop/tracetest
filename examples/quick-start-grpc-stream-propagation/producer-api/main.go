package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/kubeshop/tracetest/quick-start-grpc-stream-propagation/producer-api/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	grpc "google.golang.org/grpc"
)

// Implement the PaymentReceiverServer interface
type serverImpl struct {
	pb.PaymentReceiverServer
	tracer trace.Tracer
}

type paymentWithMetadata struct {
	payment  *pb.Payment
	metadata map[string]string
}

var _ pb.PaymentReceiverServer = &serverImpl{}
var (
	paymentChannel = make(chan *paymentWithMetadata)
)

func (s *serverImpl) ReceivePayment(ctx context.Context, payment *pb.Payment) (*pb.ReceivePaymentResponse, error) {
	go func() {
		ctx, span := s.tracer.Start(ctx, "EnqueuePayment")
		defer span.End()

		message := &paymentWithMetadata{
			payment:  payment,
			metadata: extractMetadataFromContext(ctx),
		}

		// handle channel as in-memory queue
		paymentChannel <- message
	}()

	return &pb.ReceivePaymentResponse{Received: true}, nil
}

func (s *serverImpl) NotifyPayment(_ *pb.Empty, stream pb.PaymentReceiver_NotifyPaymentServer) error {
	for {
		message, ok := <-paymentChannel
		if !ok {
			return nil
		}

		ctx := injectMetadataIntoContext(context.Background(), message.metadata)
		ctx, span := s.tracer.Start(ctx, "SendPaymentNotification")

		payment := message.payment
		highValuePayment := payment.Amount > 10_000

		notification := &pb.PaymentNotification{
			Payment:          payment,
			HighValuePayment: highValuePayment,
		}

		// extract OTel data from context and add it to the notification
		notification.Metadata = extractMetadataFromContext(ctx)

		if err := stream.Send(notification); err != nil {
			return err
		}

		span.End()
	}
}

func main() {
	port := getEnvVar("PORT", "8080")
	otelExporterEndpoint := getEnvVar("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	otelServiceName := getEnvVar("OTEL_SERVICE_NAME", "producer-api")

	tracer, err := setupOpenTelemetry(context.Background(), otelExporterEndpoint, otelServiceName)
	if err != nil {
		log.Fatalf("failed to initialize OpenTelemetry: %v", err)
		return
	}

	address := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	log.Printf("server listening at %s", lis.Addr())

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	paymentReceiverServer := &serverImpl{
		tracer: tracer,
	}

	pb.RegisterPaymentReceiverServer(grpcServer, paymentReceiverServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
