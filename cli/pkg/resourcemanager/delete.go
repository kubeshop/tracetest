package resourcemanager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

const VerbDelete Verb = "delete"

func (c Client) Delete(ctx context.Context, id string, format Format) (string, error) {
	prefix := ""
	if c.options.prefixGetterFn != nil {
		prefix = c.options.prefixGetterFn()
	}

	url := c.client.url(c.resourceNamePlural, prefix, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot build Delete request: %w", err)
	}

	err = format.BuildRequest(req, VerbDelete)
	if err != nil {
		return "", fmt.Errorf("cannot build Delete request: %w", err)
	}

	resp, err := c.client.do(req)
	c.logger.Debug("Resource Delete", zap.String("request", spew.Sdump(req)))
	if err != nil {
		return "", fmt.Errorf("cannot execute Delete request: %w", err)
	}
	defer resp.Body.Close()
	c.logger.Debug("Resource Delete", zap.String("response.status", resp.Status))

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

	c.logger.Debug("Resource Delete", zap.String("response.msg", msg))
	return msg, nil
}
