package resourcemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Verb string

const (
	VerbList Verb = "list"
)

type client struct {
	client             http.Client
	baseURL            string
	resourceName       string
	resourceNamePlural string
	mappings           []TableCellConfig
}

// NewClient creates a new client for a resource managed by the resourceamanger.
// The mappings parameter is a list of column names to the json path of the value.
// The path is relative to the Resource root, so the client does not need to consider the ResourceList wrapper struct.
// Example:
//
//	mappings := []TableCellConfig{
//		{Header: "ID", Path: "id"},
//		{Header: "Name", Path: "metadata.name"},
//		{Header: "Namespace", Path: "metadata.namespace"},
//		{Header: "Age", Path: "metadata.creationTimestamp"},
//	}
//
// this example would work both for a single resource from a Get, or a ResourceList from a List
func NewClient(baseURL, resourceName, resourceNamePlural string, mappings []TableCellConfig) client {
	return client{
		baseURL:            baseURL,
		resourceName:       resourceName,
		resourceNamePlural: resourceNamePlural,
		mappings:           mappings,
	}
}

func (c client) url(extra ...string) string {
	return fmt.Sprintf("%s/api/%s/%s", c.baseURL, c.resourceNamePlural, strings.Join(extra, "/"))
}

type ListOption struct {
	Take          int32
	Skip          int32
	SortDirection string
	SortBy        string
	All           bool
}

func (p *ListOption) Validate(args []string) []error {
	errors := make([]error, 0)

	if p.Take < 0 {
		errors = append(errors, fmt.Errorf("[take] must be greater than 0"))
	}

	if p.Skip < 0 {
		errors = append(errors, fmt.Errorf("[skip] must be greater than 0"))
	}

	if p.SortDirection != "asc" && p.SortDirection != "desc" {
		errors = append(errors, fmt.Errorf("[sortDirection] must be either asc or desc"))
	}

	return errors
}

func (c client) List(ctx context.Context, opt ListOption, format Format) (string, error) {
	url := c.url()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot build List request: %w", err)
	}

	err = format.BuildRequest(req, VerbList)
	if err != nil {
		return "", fmt.Errorf("cannot build List request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute List request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read List response: %w", err)
	}

	return format.Format(string(body), c.mappings)
}
