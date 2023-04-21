package otlp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	ptraceotlp "go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
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
	ingester ingester

	hServer *http.Server
}

func NewHttpServer(addr string, ingester ingester) *httpServer {
	return &httpServer{
		addr:     addr,
		ingester: ingester,
	}
}

func (s *httpServer) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/v1/traces", s.Export).Methods("POST")

	s.hServer = &http.Server{
		Addr:    s.addr,
		Handler: handlers.CompressHandler(decompressBodyHandler(handlers.ContentTypeHandler(r, protoBufContentType, jsonContentType))),
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("cannot listen on address %s: %w", s.addr, err)
	}

	return s.hServer.Serve(listener)
}

func (s *httpServer) Stop() {
	s.hServer.Close()
}

func (s httpServer) Export(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("content-type")
	response := newHttpResponse(w, contentType)

	request, err := s.parseBody(r.Body, contentType)
	if err != nil {
		response.sendError(http.StatusUnprocessableEntity, status.Errorf(codes.InvalidArgument, "Could not parse request body %s", err.Error()))
		return
	}

	result, err := s.ingester.Ingest(r.Context(), request)
	if err != nil {
		response.sendError(http.StatusInternalServerError, status.Errorf(codes.InvalidArgument, "Error when ingesting spans %s", err.Error()))
		return
	}

	response.send(http.StatusOK, result)
}

func (s httpServer) parseProtoBuf(body []byte) (*pb.ExportTraceServiceRequest, error) {
	request := pb.ExportTraceServiceRequest{}

	err := proto.Unmarshal(body, &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s httpServer) parseJson(body []byte) (*pb.ExportTraceServiceRequest, error) {
	exportRequest := ptraceotlp.NewRequest()

	err := exportRequest.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}

	protoBody, err := exportRequest.MarshalProto()
	if err != nil {
		return nil, err
	}

	return s.parseProtoBuf(protoBody)
}

func (s httpServer) parseBody(reqBody io.ReadCloser, contentType string) (*pb.ExportTraceServiceRequest, error) {
	var body []byte
	if b, err := io.ReadAll(reqBody); err == nil {
		body = b
	} else {
		return nil, err
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
		fmt.Println("Could not attach body to response", err.Error())
		return err
	}

	r.w.WriteHeader(statusCode)
	r.w.Write(body)

	return nil
}

func (r httpResponse) sendError(code int, err error) {
	rpcError, _ := status.FromError(err)

	r.send(code, rpcError.Proto())
}

func (r httpResponse) paseResponseBody(data proto.Message) ([]byte, error) {
	if r.contentType == protoBufContentType {
		return proto.Marshal(data)
	}

	return json.Marshal(data)
}

func decompressBodyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("content-encoding"), "gzip") {
			compressedBody, err := decompressBody(r.Body)
			if err != nil {
				response := newHttpResponse(w, r.Header.Get("content-type"))
				response.sendError(http.StatusUnprocessableEntity, status.Errorf(codes.InvalidArgument, "Could not parse request body %s", err.Error()))
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

	output, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return io.NopCloser(bytes.NewReader(output)), nil
}
