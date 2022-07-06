package model

import (
	"net/http"
)

const HTTPTriggerType TriggerType = "http"

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
	Key, Value string
}

type HTTPRequest struct {
	Method  HTTPMethod
	URL     string
	Headers []HTTPHeader
	Body    string
	Auth    *HTTPAuthenticator
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
	Type  string
	Props map[string]string
}

func (a HTTPAuthenticator) AuthenticateGRPC() {}
func (a HTTPAuthenticator) AuthenticateHTTP(req *http.Request) {
	var auth authenticator
	switch a.Type {
	case "apiKey":
		auth = APIKeyAuthenticator{
			Key:   a.Props["key"],
			Value: a.Props["value"],
			In:    APIKeyPosition(a.Props["in"]),
		}
	case "basic":
		auth = BasicAuthenticator{
			Username: a.Props["username"],
			Password: a.Props["password"],
		}
	case "bearer":
		auth = BearerAuthenticator{
			Bearer: a.Props["token"],
		}
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
	Key   string
	Value string
	In    APIKeyPosition
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
	Username string
	Password string
}

func (a BasicAuthenticator) AuthenticateGRPC() {}
func (a BasicAuthenticator) AuthenticateHTTP(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}

type BearerAuthenticator struct {
	Bearer string
}

func (a BearerAuthenticator) AuthenticateGRPC() {}
func (a BearerAuthenticator) AuthenticateHTTP(req *http.Request) {
	req.Header.Add("Authorization", a.Bearer)
}
