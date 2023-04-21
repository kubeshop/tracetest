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

	ptraceotlp "go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const protoBufContentType = "application/x-protobuf"
const jsonContentType = "application/json"
const gzipEncoding = "gzip"

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
	contentType, contentEncoding, err := s.validateRequest(r)
	shouldDecompress := contentEncoding == gzipEncoding
	shouldEncode := strings.Contains(r.Header.Get("Accept-Encoding"), gzipEncoding)

	response := newHttpResponse(w, contentType, shouldEncode)
	if err != nil {
		response.sendError(http.StatusBadRequest, status.Errorf(codes.InvalidArgument, "Error when validating request %s", err.Error()))
		return
	}

	request, err := s.parseBody(r.Body, contentType, shouldDecompress)
	if err != nil {
		response.sendError(http.StatusUnprocessableEntity, status.Errorf(codes.InvalidArgument, "Could not parse request body %s", err.Error()))
		return
	}

	result, err := s.exporter.Ingest(r.Context(), request)
	if err != nil {
		response.sendError(http.StatusInternalServerError, status.Errorf(codes.InvalidArgument, "Error when ingesting spans %s", err.Error()))
		return
	}
	response.send(http.StatusOK, result)
}

func (s HttpServer) parseProtoBuf(body []byte) (*pb.ExportTraceServiceRequest, error) {
	request := pb.ExportTraceServiceRequest{}

	err := proto.Unmarshal(body, &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s HttpServer) parseJson(body []byte) (*pb.ExportTraceServiceRequest, error) {
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

func (s HttpServer) parseBody(reqBody io.ReadCloser, contentType string, shouldDecompress bool) (*pb.ExportTraceServiceRequest, error) {
	var body []byte
	if b, err := io.ReadAll(reqBody); err == nil {
		body = b
	} else {
		return nil, err
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("empty body")
	}

	if shouldDecompress {
		decompressedBody, err := s.decompressBody(body)
		if err != nil {
			return nil, err
		}

		body = decompressedBody
	}

	switch contentType {
	case protoBufContentType:
		return s.parseProtoBuf(body)
	case jsonContentType:
		return s.parseJson(body)
	}

	return nil, fmt.Errorf("content-type %s not supported", contentType)
}

func (s HttpServer) validateRequest(r *http.Request) (string, string, error) {
	contentEncoding := r.Header.Get("content-encoding")
	contentType := r.Header.Get("content-type")

	if r.URL.Path != "/v1/traces" {
		return contentType, contentEncoding, fmt.Errorf("path %s not supported", r.URL.Path)
	}
	if r.Method != http.MethodPost {
		return contentType, contentEncoding, fmt.Errorf("method %s not supported", r.Method)
	}
	if contentType != "" && contentType != protoBufContentType && contentType != jsonContentType {
		return contentType, contentEncoding, fmt.Errorf("content-type %s not supported", contentType)
	}

	return contentType, contentEncoding, nil
}

func (s HttpServer) decompressBody(body []byte) ([]byte, error) {
	reader := bytes.NewReader(body)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return output, nil
}

type httpResponse struct {
	w            http.ResponseWriter
	shouldEncode bool
	contentType  string
}

func newHttpResponse(w http.ResponseWriter, contentType string, shouldEncode bool) httpResponse {

	return httpResponse{
		w:            w,
		shouldEncode: shouldEncode,
		contentType:  contentType,
	}
}

func (r httpResponse) send(statusCode int, message proto.Message) error {
	body, err := r.paseBody(message)
	if err != nil {
		fmt.Println("Could not attach body to response", err.Error())
		return err
	}

	r.w.Header().Set("Content-Type", r.contentType)
	if r.shouldEncode {
		r.w.Header().Set("Content-Encoding", gzipEncoding)
		r.w.WriteHeader(statusCode)
		r.w.Write(r.compressBody(body))
	} else {
		r.w.WriteHeader(statusCode)
		r.w.Write(body)
	}

	return nil
}

func (r httpResponse) sendError(code int, err error) {
	rpcError, _ := status.FromError(err)

	r.send(code, rpcError.Proto())
}

func (r httpResponse) paseBody(data proto.Message) ([]byte, error) {
	switch r.contentType {
	case protoBufContentType:
		return proto.Marshal(data)
	case jsonContentType:
		return json.Marshal(data)
	}

	return []byte(fmt.Sprintf("content-type <%s> not supported", r.contentType)), nil
}

func (r httpResponse) compressBody(body []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(body)
	w.Close()

	return b.Bytes()
}
