package main

import (
	"context"
	"io"
	"log"

	pb "github.com/kubeshop/tracetest/quick-start-grpc-stream-propagation/consumer-worker/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	producerAPIAddress := getEnvVar("PRODUCER_API_ADDRESS", "localhost:8080")
	otelExporterEndpoint := getEnvVar("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317")
	otelServiceName := getEnvVar("OTEL_SERVICE_NAME", "producer-api")

	tracer, err := setupOpenTelemetry(context.Background(), otelExporterEndpoint, otelServiceName)
	if err != nil {
		log.Fatalf("failed to initialize OpenTelemetry: %v", err)
		return
	}

	grpcClient, err := grpc.NewClient(
		producerAPIAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("could not connect to producer API: %v", err)
	}

	log.Printf("Connected to producer API at %s", producerAPIAddress)

	client := pb.NewPaymentReceiverClient(grpcClient)

	stream, err := client.NotifyPayment(ctx, &pb.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		log.Fatalf("could not receive payment notifications: %v", err)
	}

	log.Printf("Listening for payment notifications...")

	for {
		notification, err := stream.Recv()
		if err == io.EOF {
			log.Printf("No more payment notifications")
			return
		}
		if err != nil {
			log.Fatalf("could not receive payment notification: %v", err)
		}

		processPaymentNotification(tracer, notification)
	}
}

func processPaymentNotification(tracer trace.Tracer, notification *pb.PaymentNotification) {
	messageProcessingCtx := injectMetadataIntoContext(context.Background(), notification.Metadata)
	_, span := tracer.Start(messageProcessingCtx, "ProcessPaymentNotification")
	defer span.End()

	log.Printf("Received payment notification: %v", notification)
}
