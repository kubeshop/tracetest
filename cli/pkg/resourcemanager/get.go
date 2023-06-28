package resourcemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const VerbGet Verb = "get"

func (c client) Get(ctx context.Context, id string, format Format) (string, error) {
	url := c.client.url(c.resourceNamePlural, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot build Get request: %w", err)
	}

	err = format.BuildRequest(req, VerbGet)
	if err != nil {
		return "", fmt.Errorf("cannot build Get request: %w", err)
	}

	resp, err := c.client.do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute Get request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := parseRequestError(resp, format)
		reqErr, ok := err.(requestError)
		if ok && reqErr.Code == http.StatusNotFound {
			return fmt.Sprintf("Resource %s with ID %s not found", c.resourceName, id), nil
		}

		return "", fmt.Errorf("could not Get resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read Get response: %w", err)
	}

	return format.Format(string(body), c.tableConfig)
}
