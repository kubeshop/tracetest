package resourcemanager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const VerbDelete Verb = "delete"

func (c Client) Delete(ctx context.Context, id string, format Format) (string, error) {
	url := c.client.url(c.resourceNamePlural, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot build Delete request: %w", err)
	}

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
		if errors.Is(err, ErrNotFound) {
			return fmt.Sprintf("Resource %s with ID %s not found", c.resourceName, id), ErrNotFound
		}

		return "", fmt.Errorf("could not Delete resource: %w", err)
	}

	msg := ""
	if c.options.deleteSuccessMsg != "" {
		msg = c.options.deleteSuccessMsg
	} else {
		ucfirst := strings.ToUpper(string(c.resourceName[0])) + c.resourceName[1:]
		msg = fmt.Sprintf("%s successfully deleted", ucfirst)
	}

	return msg, nil
}
