package resourcemanager

import (
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

func (c HTTPClient) do(req *http.Request) (*http.Response, error) {
	for k, v := range c.extraHeaders {
		req.Header[k] = v
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

var ErrNotFound = requestError{
	Code:    http.StatusNotFound,
	Message: "Resource not found",
}

type requestError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e requestError) Error() string {
	return e.Message
}

func (e requestError) Is(target error) bool {
	t, ok := target.(requestError)
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
		return requestError{
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
	}
	var reqErr requestError
	err = format.Unmarshal(body, &reqErr)
	if err != nil {
		return fmt.Errorf("cannot parse response body: %w", err)
	}

	return reqErr
}
