package model

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

type HTTPAuth struct {
	// TODO: model this
}

type HTTPRequest struct {
	Method  HTTPMethod
	URL     string
	Headers []HTTPHeader
	Body    string
	Auth    HTTPAuth
}

type HTTPResponse struct {
	Status     string
	StatusCode int
	Headers    []HTTPHeader
	Body       string
}
