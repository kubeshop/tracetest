package resourcemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ListOption struct {
	Take          int32
	Skip          int32
	SortDirection string
	SortBy        string
	All           bool
}

const VerbList Verb = "list"

func (c Client) List(ctx context.Context, opt ListOption, format Format) (string, error) {
	prefix := ""
	if c.options.prefixGetterFn != nil {
		prefix = c.options.prefixGetterFn()
	}
	url := c.client.url(c.resourceNamePlural, prefix)

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

	fmt.Println("@@@", resp.Status, resp.StatusCode, req.URL.String())
	defer resp.Body.Close()

	if !isSuccessResponse(resp) {
		err := parseRequestError(resp, format)

		return "", fmt.Errorf("could not list resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read List response: %w", err)
	}

	return format.Format(string(body), c.options.tableConfig, c.options.listPath)
}
