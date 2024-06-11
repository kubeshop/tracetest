package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/kubeshop/tracetest/quick-start-grpc-stream-propagation/producer-api/proto"
	grpc "google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPaymentReceiverServer
}

var _ pb.PaymentReceiverServer = &server{}
var paymentChannel = make(chan *pb.Payment)

func (s *server) ReceivePayment(ctx context.Context, payment *pb.Payment) (*pb.ReceivePaymentResponse, error) {
	go func() {
		paymentChannel <- payment
	}()

	return &pb.ReceivePaymentResponse{Received: true}, nil
}

func (s *server) NotifyPayment(_ *pb.Empty, stream pb.PaymentReceiver_NotifyPaymentServer) error {
	for {
		payment, ok := <-paymentChannel
		if !ok {
			return nil
		}

		highValuePayment := payment.Amount > 10_000

		notification := &pb.PaymentNotification{
			Payment:          payment,
			HighValuePayment: highValuePayment,
		}

		if err := stream.Send(notification); err != nil {
			return err
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	log.Printf("server listening at %s", lis.Addr())

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentReceiverServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
