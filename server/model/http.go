package model

import "net/http"

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
	Auth    HTTPAuthenticator
}

type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    []HTTPHeader
	Body       string
}

type HTTPAuthenticator interface {
	Authenticate(*http.Request)
}

type APIKeyPosition string

const (
	APIKeyPositionHeader APIKeyPosition = "header"
	APIKeyPositionQuery  APIKeyPosition = "query"
)

type APIKeyAuthenticator struct {
	Key   string
	Value string
	In    APIKeyPosition
}

func (a APIKeyAuthenticator) Authenticate(req *http.Request) {
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

func (a BasicAuthenticator) Authenticate(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}

type BearerAuthenticator struct {
	Bearer string
}

func (a BearerAuthenticator) Authenticate(req *http.Request) {
	req.Header.Add("Authorization", a.Bearer)
}
