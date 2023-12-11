package otlp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	protoBufContentType = "application/x-protobuf"
	jsonContentType     = "application/json"
)

type httpServer struct {
	addr     string
	ingester Ingester
	logger   *zap.Logger

	hServer *http.Server
}

func NewHttpServer(addr string, ingester Ingester) *httpServer {
	return &httpServer{
		addr:     addr,
		ingester: ingester,
		logger:   zap.NewNop(),
	}
}

func (s *httpServer) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

func (s *httpServer) Start() error {
	s.logger.Debug("Starting HTTP server", zap.String("addr", s.addr))
	r := mux.NewRouter()
	r.HandleFunc("/v1/traces", s.Export).Methods("POST")

	s.hServer = &http.Server{
		Addr:    s.addr,
		Handler: handlers.CompressHandler(decompressBodyHandler(s.logger, handlers.ContentTypeHandler(r, protoBufContentType, jsonContentType))),
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.logger.Error("cannot listen on address", zap.String("addr", s.addr), zap.Error(err))
		return fmt.Errorf("cannot listen on address %s: %w", s.addr, err)
	}

	go s.hServer.Serve(listener)
	s.logger.Debug("HTTP server started", zap.String("addr", s.addr))
	return nil
}

func (s *httpServer) Stop() {
	s.hServer.Close()
}

func (s httpServer) Export(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("content-type")
	response := newHttpResponse(w, contentType)
	r.Header.Set("Access-Control-Allow-Origin", "*")

	s.logger.Debug("Received ExportTraceServiceRequest", zap.String("content-type", contentType))

	request, err := s.parseBody(r.Body, contentType)
	if err != nil {
		s.logger.Error("Could not parse request body", zap.Error(err))
		sendErr := response.sendError(http.StatusUnprocessableEntity, status.Errorf(codes.InvalidArgument, "Could not parse request body %s", err.Error()))
		if sendErr != nil {
			s.logger.Error("Could not send error response", zap.Error(sendErr))
		}
		return
	}
	s.logger.Debug("Parsed ExportTraceServiceRequest", zap.Any("request", request))

	result, err := s.ingester.Ingest(r.Context(), request, "HTTP")
	if err != nil {
		s.logger.Error("Error when ingesting spans", zap.Error(err))
		sendErr := response.sendError(http.StatusInternalServerError, status.Errorf(codes.InvalidArgument, "Error when ingesting spans %s", err.Error()))
		if sendErr != nil {
			s.logger.Error("Could not send error response", zap.Error(sendErr))
		}
		return
	}
	s.logger.Debug("Ingested spans", zap.Any("result", result))

	err = response.send(http.StatusOK, result)
	if err != nil {
		s.logger.Error("Error when sending response", zap.Error(err))
	}

	s.logger.Debug("Sent ExportTraceServiceResponse")
}

func (s httpServer) parseProtoBuf(body []byte) (*pb.ExportTraceServiceRequest, error) {
	request := pb.ExportTraceServiceRequest{}

	err := proto.Unmarshal(body, &request)
	if err != nil {
		return nil, fmt.Errorf("cannot parse protobuf: %w", err)
	}

	return &request, nil
}

func (s httpServer) parseJson(body []byte) (*pb.ExportTraceServiceRequest, error) {
	exportRequest := ptraceotlp.NewExportRequest()

	err := exportRequest.UnmarshalJSON(body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse json: %w", err)
	}

	protoBody, err := exportRequest.MarshalProto()
	if err != nil {
		return nil, fmt.Errorf("cannot marshalProto: %w", err)
	}

	return s.parseProtoBuf(protoBody)
}

func (s httpServer) parseBody(reqBody io.ReadCloser, contentType string) (*pb.ExportTraceServiceRequest, error) {
	body, err := io.ReadAll(reqBody)
	if err != nil {
		return nil, fmt.Errorf("cannot read request body: %w", err)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("empty body")
	}

	if contentType == protoBufContentType {
		return s.parseProtoBuf(body)
	}

	return s.parseJson(body)
}

type httpResponse struct {
	w           http.ResponseWriter
	contentType string
}

func newHttpResponse(w http.ResponseWriter, contentType string) httpResponse {
	return httpResponse{
		w:           w,
		contentType: contentType,
	}
}

func (r httpResponse) send(statusCode int, message proto.Message) error {
	body, err := r.paseResponseBody(message)
	if err != nil {
		fmt.Println("could not attach body to response: %w", err)
		return err
	}

	r.w.WriteHeader(statusCode)
	r.w.Write(body)

	return nil
}

func (r httpResponse) sendError(code int, err error) error {
	rpcError, _ := status.FromError(err)

	return r.send(code, rpcError.Proto())
}

func (r httpResponse) paseResponseBody(data proto.Message) ([]byte, error) {
	if r.contentType == protoBufContentType {
		return proto.Marshal(data)
	}

	return json.Marshal(data)
}

func decompressBodyHandler(logger *zap.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("content-encoding"), "gzip") {
			logger.Debug("Decompressing request body")
			compressedBody, err := decompressBody(r.Body)
			if err != nil {
				logger.Error("Could not decompress request body", zap.Error(err))
				response := newHttpResponse(w, r.Header.Get("content-type"))
				sendErr := response.sendError(http.StatusUnprocessableEntity, status.Errorf(codes.InvalidArgument, "Could not parse request body %s", err.Error()))
				if sendErr != nil {
					logger.Error("Could not send error response", zap.Error(sendErr))
				}
				return
			}

			r.Body = compressedBody
			r.Header.Set("accept-encoding", "gzip")
		}

		h.ServeHTTP(w, r)
	})
}

func decompressBody(reqBody io.ReadCloser) (io.ReadCloser, error) {
	var body []byte
	if b, err := io.ReadAll(reqBody); err == nil {
		body = b
	} else {
		return nil, err
	}

	reader := bytes.NewReader(body)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	output, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(output)), nil
}
