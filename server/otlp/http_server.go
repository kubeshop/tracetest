package otlp

import (
	"fmt"
	"io"
	"net"
	"net/http"

	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
)

type HttpServer struct {
	addr     string
	exporter IExporter

	hServer *http.Server
}

func NewHttpServer(addr string, exporter IExporter) *HttpServer {
	return &HttpServer{
		addr:     addr,
		exporter: exporter,
	}
}

func (s *HttpServer) Start() error {
	s.hServer = &http.Server{
		Addr:    s.addr,
		Handler: s,
	}

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("cannot listen on address %s: %w", s.addr, err)
	}

	return s.hServer.Serve(listener)
}

func (s *HttpServer) Stop() {}

func (s HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if b, err := io.ReadAll(r.Body); err == nil {
		body = b
	} else {
		fmt.Println(err)
	}

	request := pb.ExportTraceServiceRequest{}
	err := proto.Unmarshal(body, &request)
	if err != nil {
		fmt.Println("Error when parsing request", err.Error())
	}

	_, err = s.exporter.Export(r.Context(), &request)
	if err != nil {
		fmt.Println("Error when exporting request", err.Error())
	}

	w.WriteHeader(http.StatusOK)
}
