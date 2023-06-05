package trigger

import (
	"net/http"
)

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
	Key   string `expr_enabled:"true"`
	Value string `expr_enabled:"true"`
}

type HTTPRequest struct {
	Method          HTTPMethod         `expr_enabled:"true" json:"method,omitempty"`
	URL             string             `expr_enabled:"true" json:"url"`
	Body            string             `expr_enabled:"true" json:"body,omitempty"`
	Headers         []HTTPHeader       `json:"headers,omitempty"`
	Auth            *HTTPAuthenticator `json:"auth,omitempty"`
	SSLVerification bool               `json:"sslVerification,omitempty"`
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
	Type   string `expr_enabled:"true"`
	APIKey APIKeyAuthenticator
	Basic  BasicAuthenticator
	Bearer BearerAuthenticator
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
	Key   string         `expr_enabled:"true"`
	Value string         `expr_enabled:"true"`
	In    APIKeyPosition `expr_enabled:"true"`
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
	Username string `expr_enabled:"true"`
	Password string `expr_enabled:"true"`
}

func (a BasicAuthenticator) AuthenticateGRPC() {}
func (a BasicAuthenticator) AuthenticateHTTP(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}

type BearerAuthenticator struct {
	Bearer string `expr_enabled:"true"`
}

func (a BearerAuthenticator) AuthenticateGRPC() {}
func (a BearerAuthenticator) AuthenticateHTTP(req *http.Request) {
	req.Header.Add("Authorization", a.Bearer)
}
