package app

import (
	"net/http"

	"google.golang.org/grpc"
)

type executionContext struct {
	httpServer *http.Server
	otlpServer *grpc.Server
}
