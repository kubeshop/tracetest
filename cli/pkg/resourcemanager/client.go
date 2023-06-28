package resourcemanager

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Verb string

const (
	VerbList Verb = "list"
)

type client struct {
	client             HTTPClient
	resourceName       string
	resourceNamePlural string
	tableConfig        TableConfig
}

type HTTPClient struct {
	client http.Client

	baseURL      string
	extraHeaders http.Header
}

func NewHTTPClient(baseURL string, extraHeaders http.Header) HTTPClient {
	return HTTPClient{
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
func NewClient(httpClient HTTPClient, resourceName, resourceNamePlural string, tableConfig TableConfig) client {
	return client{
		client:             httpClient,
		resourceName:       resourceName,
		resourceNamePlural: resourceNamePlural,
		tableConfig:        tableConfig,
	}
}

type ListOption struct {
	Take          int32
	Skip          int32
	SortDirection string
	SortBy        string
	All           bool
}

func (c client) List(ctx context.Context, opt ListOption, format Format) (string, error) {
	url := c.client.url(c.resourceNamePlural)

	q := url.Query()
	q.Add("skip", fmt.Sprintf("%d", opt.Skip))
	q.Add("take", fmt.Sprintf("%d", opt.Take))
	q.Add("sortBy", opt.SortBy)
	q.Add("sortDirection", opt.SortDirection)

	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot build List request: %w", err)
	}

	err = format.BuildRequest(req, VerbList)
	if err != nil {
		return "", fmt.Errorf("cannot build List request: %w", err)
	}

	resp, err := c.client.do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute List request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := parseRequestError(resp, format)

		return "", fmt.Errorf("could not list resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read List response: %w", err)
	}

	return format.Format(string(body), c.tableConfig)
}

type requestError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func (e requestError) Error() string {
	return e.Message
}

func parseRequestError(resp *http.Response, format Format) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	fmt.Println(string(body))

	var reqErr requestError
	err = format.Unmarshal(body, &reqErr)
	if err != nil {
		return fmt.Errorf("cannot parse response body: %w", err)
	}

	return reqErr
}
