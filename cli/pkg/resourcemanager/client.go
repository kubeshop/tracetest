package resourcemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Verb string

type Client struct {
	client             *HTTPClient
	resourceName       string
	resourceNamePlural string
	logger             *zap.Logger
	options            options
}

type HTTPClient struct {
	client http.Client

	baseURL      string
	extraHeaders http.Header
}

func NewHTTPClient(baseURL string, extraHeaders http.Header) *HTTPClient {
	return &HTTPClient{
		client: http.Client{
			// this function avoids blindly followin redirects.
			// the problem with redirects is that they don't guarantee to preserve the method, body, headers, etc.
			// This can hide issues when developing, because the client will follow the redirect and the request
			// will succeed, but the server will not receive the request that the user intended to send.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		baseURL:      baseURL,
		extraHeaders: extraHeaders,
	}
}

func (c HTTPClient) url(resourceName, prefix string, extra ...string) *url.URL {
	urlStr := c.baseURL + path.Join("/", prefix, resourceName, strings.Join(extra, "/"))
	url, _ := url.Parse(urlStr)
	return url
}

type ContextOption string

const (
	ContextHeadersKey ContextOption = "headers"
)

func SetRequestContextHeaders(ctx context.Context, headers map[string]string) context.Context {
	httpHeaders := http.Header{}
	for k, v := range headers {
		httpHeaders.Set(k, v)
	}

	return context.WithValue(ctx, ContextHeadersKey, httpHeaders)
}

func (c HTTPClient) do(req *http.Request) (*http.Response, error) {
	for k, v := range c.extraHeaders {
		req.Header[k] = v
	}

	contextHeaders := req.Context().Value(ContextHeadersKey)
	if contextHeaders != nil {
		for k, v := range contextHeaders.(http.Header) {
			req.Header[k] = v
		}
	}

	return c.client.Do(req)
}

// NewClient creates a new client for a resource managed by the resourceamanger.
// The tableConfig parameter configures how the table view should be rendered.
// This configuration work both for a single resource from a Get, or a ResourceList from a List
func NewClient(
	httpClient *HTTPClient,
	logger *zap.Logger,
	resourceName, resourceNamePlural string,
	opts ...option) Client {
	c := Client{
		client:             httpClient,
		resourceName:       resourceName,
		resourceNamePlural: resourceNamePlural,
		logger:             logger,
	}

	for _, opt := range opts {
		opt(&c.options)
	}

	return c
}

func (c Client) WithHttpClient(HTTPClient *HTTPClient) Client {
	c.client = HTTPClient
	return c
}

func (c Client) WithOptions(opts ...option) Client {
	for _, opt := range opts {
		opt(&c.options)
	}

	return c
}

func (c Client) resourceType() string {
	if c.options.resourceType != "" {
		return c.options.resourceType
	}

	// language.Und means Undefined
	caser := cases.Title(language.Und, cases.NoLower)
	return caser.String(c.resourceName)
}

var ErrNotFound = RequestError{
	Code:    http.StatusNotFound,
	Message: "Resource not found",
}

type RequestError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

type alternateRequestError struct {
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (e RequestError) Error() string {
	return e.Message
}

func (e RequestError) Is(target error) bool {
	t, ok := target.(RequestError)
	return ok && t.Code == e.Code
}

func isSuccessResponse(resp *http.Response) bool {
	// successfull http status codes are 2xx
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func parseRequestError(resp *http.Response, format Format) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	if len(body) == 0 {
		return RequestError{
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
	}
	var reqErr RequestError
	err = format.Unmarshal(body, &reqErr)
	if err != nil {
		return fmt.Errorf("cannot parse response body: %w", err)
	}

	emptyRequestError := reqErr.Code == 0 && reqErr.Message == ""
	if !emptyRequestError {
		// Success, parsed error message
		return reqErr
	}

	// Fallback, try to parse message in other format

	var alternateReqError alternateRequestError
	err = format.Unmarshal(body, &alternateReqError)
	if err != nil {
		return fmt.Errorf("cannot parse response body: %w", err)
	}

	return RequestError{
		Code:    alternateReqError.Status,
		Message: alternateReqError.Detail,
	}
}
