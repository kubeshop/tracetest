package resourcemanager

import (
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

// NewClient creates a new client for a resource managed by the resourceamanger.
// The tableConfig parameter configures how the table view should be rendered.
// This configuration work both for a single resource from a Get, or a ResourceList from a List
func NewClient(httpClient *HTTPClient, resourceName, resourceNamePlural string, tableConfig TableConfig) client {
	return client{
		client:             httpClient,
		resourceName:       resourceName,
		resourceNamePlural: resourceNamePlural,
		tableConfig:        tableConfig,
	}
}

type requestError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e requestError) Error() string {
	return e.Message
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
