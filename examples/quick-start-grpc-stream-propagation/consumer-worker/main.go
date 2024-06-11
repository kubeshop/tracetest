package main

import (
	"context"
	"io"
	"log"
	"os"

	pb "github.com/kubeshop/tracetest/quick-start-grpc-stream-propagation/consumer-worker/proto"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	producerAPIAddress := os.Getenv("PRODUCER_API_ADDRESS")
	if producerAPIAddress == "" {
		producerAPIAddress = "localhost:8080"
	}

	grpcClient, err := grpc.Dial(
		producerAPIAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("could not connect to producer API: %v", err)
	}

	log.Printf("Connected to producer API at %s", producerAPIAddress)

	client := pb.NewPaymentReceiverClient(grpcClient)

	stream, err := client.NotifyPayment(context.Background(), &pb.Empty{}, grpc.WaitForReady(true))
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

		log.Printf("Received payment notification: %v", notification)
	}
}
