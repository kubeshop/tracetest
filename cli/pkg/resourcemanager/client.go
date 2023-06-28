package resourcemanager

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Verb string

type client struct {
	client             *HTTPClient
	resourceName       string
	resourceNamePlural string
	deleteSuccessMsg   string
	tableConfig        TableConfig
}

type HTTPClient struct {
	client http.Client

	baseURL      string
	extraHeaders http.Header
}

func NewHTTPClient(baseURL string, extraHeaders http.Header) *HTTPClient {
	return &HTTPClient{
		client:       http.Client{},
		baseURL:      baseURL,
		extraHeaders: extraHeaders,
	}
}

func (c HTTPClient) url(resourceName string, extra ...string) *url.URL {
	url, _ := url.Parse(fmt.Sprintf("%s/api/%s/%s", c.baseURL, resourceName, strings.Join(extra, "/")))
	return url
}

func (c HTTPClient) do(req *http.Request) (*http.Response, error) {
	for k, v := range c.extraHeaders {
		req.Header[k] = v
	}

	return c.client.Do(req)
}

type options func(c *client)

func WithDeleteEnabled(deleteSuccessMssg string) options {
	return func(c *client) {
		c.deleteSuccessMsg = deleteSuccessMssg
	}
}

func WithTableConfig(tableConfig TableConfig) options {
	return func(c *client) {
		c.tableConfig = tableConfig
	}
}

var ErrNotSupportedResourceAction = errors.New("the specified resource type doesn't support the action")

// NewClient creates a new client for a resource managed by the resourceamanger.
// The tableConfig parameter configures how the table view should be rendered.
// This configuration work both for a single resource from a Get, or a ResourceList from a List
func NewClient(
	httpClient *HTTPClient,
	resourceName, resourceNamePlural string,
	opts ...options) client {
	c := client{
		client:             httpClient,
		resourceName:       resourceName,
		resourceNamePlural: resourceNamePlural,
	}

	for _, opt := range opts {
		opt(&c)
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
