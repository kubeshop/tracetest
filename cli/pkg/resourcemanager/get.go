package resourcemanager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"go.uber.org/zap"
)

const VerbGet Verb = "get"

func (c Client) Get(ctx context.Context, id string, format Format) (string, error) {
	url := c.client.url(c.resourceNamePlural, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot build Get request: %w", err)
	}

	err = format.BuildRequest(req, VerbGet)
	if err != nil {
		return "", fmt.Errorf("cannot build Get request: %w", err)
	}
	d, _ := httputil.DumpRequestOut(req, true)
	c.logger.Debug("get request",
		zap.String("request", string(d)),
	)

	resp, err := c.client.do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute Get request: %w", err)
	}
	defer resp.Body.Close()

	d, _ = httputil.DumpResponse(resp, true)
	c.logger.Debug("apply response",
		zap.String("response", string(d)),
	)

	if !isSuccessResponse(resp) {
		err := parseRequestError(resp, format)
		if errors.Is(err, ErrNotFound) {
			return fmt.Sprintf("Resource %s with ID %s not found", c.resourceName, id), ErrNotFound
		}

		return "", fmt.Errorf("could not Get resource: %w", err)

	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read Get response: %w", err)
	}

	return format.Format(string(body), c.options.tableConfig)
}
