package resourcemanager

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type Verb string

type Client struct {
	client             *HTTPClient
	resourceName       string
	resourceNamePlural string
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

func (c HTTPClient) url(resourceName string, extra ...string) *url.URL {
	urlStr := c.baseURL + path.Join("/api", resourceName, strings.Join(extra, "/"))
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
	resourceName, resourceNamePlural string,
	opts ...option) Client {
	c := Client{
		client:             httpClient,
		resourceName:       resourceName,
		resourceNamePlural: resourceNamePlural,
	}

	for _, opt := range opts {
		opt(&c.options)
	}

	return c
}

type requestError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e requestError) Error() string {
	return e.Message
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
