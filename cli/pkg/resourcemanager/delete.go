package resourcemanager

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
)

const VerbDelete Verb = "delete"

func (c client) Delete(ctx context.Context, id string, format Format) (string, error) {
	fmt.Println(c.deleteSuccessMsg)
	if c.deleteSuccessMsg == "" {
		return "", ErrNotSupportedResourceAction
	}

	url := c.client.url(c.resourceNamePlural, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot build Delete request: %w", err)
	}

	d, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(d))

	err = format.BuildRequest(req, VerbDelete)
	if err != nil {
		return "", fmt.Errorf("cannot build Delete request: %w", err)
	}

	resp, err := c.client.do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute Delete request: %w", err)
	}
	defer resp.Body.Close()

	if !isSuccessResponse(resp) {
		err := parseRequestError(resp, format)
		reqErr, ok := err.(requestError)
		if ok && reqErr.Code == http.StatusNotFound {
			return "", fmt.Errorf("Resource %s with ID %s not found", c.resourceName, id)
		}

		return "", fmt.Errorf("could not Delete resource: %w", err)
	}

	return c.deleteSuccessMsg, nil
}
