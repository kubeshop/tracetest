package trigger

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goware/urlx"
	"go.opentelemetry.io/otel/propagation"
)

func HTTP() Triggerer {
	return &httpTriggerer{}
}

type httpTriggerer struct{}

func httpClient(sslVerification bool) http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !sslVerification},
	}

	return http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

func (te *httpTriggerer) Trigger(ctx context.Context, triggerConfig Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: TriggerResult{
			Type: te.Type(),
		},
	}

	if triggerConfig.Type != TriggerTypeHTTP {
		return response, fmt.Errorf(`trigger type "%s" not supported by HTTP triggerer`, triggerConfig.Type)
	}

	client := httpClient(triggerConfig.HTTP.SSLVerification)

	ctx, cncl := context.WithTimeout(ctx, 30*time.Second)
	defer cncl()

	tReq := triggerConfig.HTTP
	var body io.Reader
	if tReq.Body != "" {
		body = bytes.NewBufferString(tReq.Body)
	}

	parsedUrl, err := urlx.Parse(tReq.URL)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(string(tReq.Method)), parsedUrl.String(), body)
	if err != nil {
		return response, err
	}
	for _, h := range tReq.Headers {
		req.Header.Add(h.Key, h.Value)
	}

	tReq.Authenticate(req)
	propagators().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return response, err
	}

	mapped := mapResp(resp)
	response.Result.HTTP = &mapped
	response.SpanAttributes = map[string]string{
		"tracetest.run.trigger.http.response_code": strconv.Itoa(resp.StatusCode),
	}

	return response, nil
}

func (t *httpTriggerer) Type() TriggerType {
	return TriggerTypeHTTP
}

func mapResp(resp *http.Response) HTTPResponse {
	var mappedHeaders []HTTPHeader
	for key, headers := range resp.Header {
		for _, val := range headers {
			val := HTTPHeader{
				Key:   key,
				Value: val,
			}
			mappedHeaders = append(mappedHeaders, val)
		}
	}
	var body string
	if b, err := io.ReadAll(resp.Body); err == nil {
		body = string(b)
	} else {
		fmt.Println(err)
	}

	return HTTPResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    mappedHeaders,
		Body:       body,
	}
}

const TriggerTypeHTTP TriggerType = "http"

type HTTPMethod string

var (
	HTTPMethodGET      HTTPMethod = "GET"
	HTTPMethodPUT      HTTPMethod = "PUT"
	HTTPMethodPOST     HTTPMethod = "POST"
	HTTPMethodPATCH    HTTPMethod = "PATCH"
	HTTPMethodDELETE   HTTPMethod = "DELETE"
	HTTPMethodCOPY     HTTPMethod = "COPY"
	HTTPMethodHEAD     HTTPMethod = "HEAD"
	HTTPMethodOPTIONS  HTTPMethod = "OPTIONS"
	HTTPMethodLINK     HTTPMethod = "LINK"
	HTTPMethodUNLINK   HTTPMethod = "UNLINK"
	HTTPMethodPURGE    HTTPMethod = "PURGE"
	HTTPMethodLOCK     HTTPMethod = "LOCK"
	HTTPMethodUNLOCK   HTTPMethod = "UNLOCK"
	HTTPMethodPROPFIND HTTPMethod = "PROPFIND"
	HTTPMethodVIEW     HTTPMethod = "VIEW"
)

type HTTPHeader struct {
	Key   string `expr_enabled:"true" json:"key,omitempty"`
	Value string `expr_enabled:"true" json:"value,omitempty"`
}

type HTTPRequest struct {
	Method          HTTPMethod         `expr_enabled:"true" json:"method,omitempty"`
	URL             string             `expr_enabled:"true" json:"url,omitempty"`
	Body            string             `expr_enabled:"true" json:"body,omitempty"`
	Headers         []HTTPHeader       `json:"headers,omitempty"`
	Auth            *HTTPAuthenticator `json:"auth,omitempty"`
	SSLVerification bool               `json:"sslVerification,omitempty"`
}

type request struct {
	Method          HTTPMethod         `expr_enabled:"true" json:"method,omitempty"`
	URL             string             `expr_enabled:"true" json:"url,omitempty"`
	Body            string             `expr_enabled:"true" json:"body,omitempty"`
	Headers         []HTTPHeader       `json:"headers,omitempty"`
	Auth            *HTTPAuthenticator `json:"auth,omitempty"`
	SSLVerification bool               `json:"sslVerification,omitempty"`
}

func (r *HTTPRequest) UnmarshalJSON(data []byte) error {
	request := request{}
	err := json.Unmarshal(data, &request)
	if err != nil {
		return err
	}

	r.Method = request.Method
	r.URL = request.URL
	r.Body = request.Body
	r.Headers = request.Headers
	r.SSLVerification = request.SSLVerification
	if request.Auth != nil && request.Auth.Type != "" {
		r.Auth = request.Auth
	}

	return nil
}

func (a HTTPRequest) Authenticate(req *http.Request) {
	if a.Auth == nil {
		return
	}

	a.Auth.AuthenticateHTTP(req)
}

type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    []HTTPHeader
	Body       string
}

type HTTPAuthenticator struct {
	Type   string               `json:"type,omitempty" expr_enabled:"true"`
	APIKey *APIKeyAuthenticator `json:"apiKey,omitempty"`
	Basic  *BasicAuthenticator  `json:"basic,omitempty"`
	Bearer *BearerAuthenticator `json:"bearer,omitempty"`
}

type auth struct {
	Type   string               `json:"type,omitempty" expr_enabled:"true"`
	APIKey *APIKeyAuthenticator `json:"apiKey,omitempty"`
	Basic  *BasicAuthenticator  `json:"basic,omitempty"`
	Bearer *BearerAuthenticator `json:"bearer,omitempty"`
}

func (a *HTTPAuthenticator) UnmarshalJSON(data []byte) error {
	auth := auth{}
	err := json.Unmarshal(data, &auth)
	if err != nil {
		return err
	}

	a.Type = auth.Type
	if auth.Type == "apiKey" {
		a.APIKey = auth.APIKey
	}

	if auth.Type == "basic" {
		a.Basic = auth.Basic
	}

	if auth.Type == "bearer" {
		a.Bearer = auth.Bearer
	}

	return nil
}

func (a HTTPAuthenticator) Map(mapFn func(current string) (string, error)) (HTTPAuthenticator, error) {
	var err error
	switch a.Type {
	case "apiKey":
		in := string(a.APIKey.In)
		in, err = mapFn(in)
		if err != nil {
			return a, err
		}
		a.APIKey.In = APIKeyPosition(in)
		a.APIKey.Key, err = mapFn(a.APIKey.Key)
		if err != nil {
			return a, err
		}
		a.APIKey.Value, err = mapFn(a.APIKey.Value)
		if err != nil {
			return a, err
		}
	case "basic":
		a.Basic.Username, err = mapFn(a.Basic.Username)
		if err != nil {
			return a, err
		}
		a.Basic.Password, err = mapFn(a.Basic.Password)
		if err != nil {
			return a, err
		}
	case "bearer":
		a.Bearer.Bearer, err = mapFn(a.Bearer.Bearer)
		if err != nil {
			return a, err
		}
	}
	return a, nil
}

func (a HTTPAuthenticator) AuthenticateGRPC() {}
func (a HTTPAuthenticator) AuthenticateHTTP(req *http.Request) {
	var auth authenticator
	switch a.Type {
	case "apiKey":
		auth = a.APIKey
	case "basic":
		auth = a.Basic
	case "bearer":
		auth = a.Bearer
	default:
		return
	}

	auth.AuthenticateHTTP(req)
}

type APIKeyPosition string

const (
	APIKeyPositionHeader APIKeyPosition = "header"
	APIKeyPositionQuery  APIKeyPosition = "query"
)

type authenticator interface {
	AuthenticateHTTP(req *http.Request)
	AuthenticateGRPC()
}

type APIKeyAuthenticator struct {
	Key   string         `json:"key,omitempty" expr_enabled:"true"`
	Value string         `json:"value,omitempty" expr_enabled:"true"`
	In    APIKeyPosition `json:"in,omitempty" expr_enabled:"true"`
}

func (a APIKeyAuthenticator) AuthenticateGRPC() {}
func (a APIKeyAuthenticator) AuthenticateHTTP(req *http.Request) {
	switch a.In {
	case APIKeyPositionHeader:
		req.Header.Set(a.Key, a.Value)
	case APIKeyPositionQuery:
		q := req.URL.Query()
		q.Add(a.Key, a.Value)
		req.URL.RawQuery = q.Encode()
	}
}

type BasicAuthenticator struct {
	Username string `json:"username,omitempty" expr_enabled:"true"`
	Password string `json:"password,omitempty" expr_enabled:"true"`
}

func (a BasicAuthenticator) AuthenticateGRPC() {}
func (a BasicAuthenticator) AuthenticateHTTP(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}

type BearerAuthenticator struct {
	Bearer string `json:"bearer,omitempty" expr_enabled:"true"`
}

func (a BearerAuthenticator) AuthenticateGRPC() {}
func (a BearerAuthenticator) AuthenticateHTTP(req *http.Request) {
	req.Header.Add("Authorization", a.Bearer)
}
