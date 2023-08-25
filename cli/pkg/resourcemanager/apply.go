package resourcemanager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/Jeffail/gabs/v2"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"go.uber.org/zap"
)

const VerbApply Verb = "apply"

type applyPreProcessorFn func(context.Context, fileutil.File) (fileutil.File, error)

func (c Client) validType(inputFile fileutil.File) error {
	c.logger.Debug("Validating resource type", zap.String("inputFile", inputFile.AbsPath()))

	var raw any
	err := (yamlFormat{}).Unmarshal(inputFile.Contents(), &raw)
	if err != nil {
		return fmt.Errorf("cannot unmarshal yaml: %w", err)
	}
	c.logger.Debug("Unmarshaled yaml", zap.Any("raw", raw))

	parsed := gabs.Wrap(raw)
	rawType := parsed.Path("type").Data()
	if rawType == nil {
		return errors.New("cannot find type in yaml")
	}
	c.logger.Debug("Found type", zap.String("type", fmt.Sprintf("%T", rawType)))
	t, ok := rawType.(string)
	if !ok {
		return fmt.Errorf("cannot parse type from yaml: %w", err)
	}
	c.logger.Debug("Parsed type", zap.String("type", t))

	if t != c.resourceType() && t != c.options.deprecatedAlias {
		return fmt.Errorf("cannot apply %s to %s resource", t, c.resourceType())
	}

	c.logger.Debug("resource type is valid")

	return nil

}

func (c Client) Apply(ctx context.Context, inputFile fileutil.File, requestedFormat Format) (string, error) {
	originalInputFile := inputFile

	if err := c.validType(inputFile); err != nil {
		return "", err
	}

	c.logger.Debug("Applying resource",
		zap.String("format", requestedFormat.String()),
		zap.String("resource", c.resourceName),
		zap.String("inputFile", inputFile.AbsPath()),
		zap.String("contents", string(inputFile.Contents())),
	)

	if c.options.applyPreProcessor != nil {
		var err error
		inputFile, err = c.options.applyPreProcessor(ctx, inputFile)
		if err != nil {
			return "", fmt.Errorf("cannot preprocess Apply request: %w", err)
		}
	}

	c.logger.Debug("preprocessed",
		zap.String("inputFile", inputFile.AbsPath()),
		zap.String("contents", string(inputFile.Contents())),
	)

	url := c.client.url(c.resourceNamePlural, "")
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), inputFile.Reader())
	if err != nil {
		return "", fmt.Errorf("cannot build Apply request: %w", err)
	}

	// we want the response inthe user's requested format
	err = requestedFormat.BuildRequest(req, VerbApply)
	if err != nil {
		return "", fmt.Errorf("cannot build Apply request: %w", err)
	}

	// the files must be in yaml format, so we can safely force the content type,
	// even if it doesn't matcht he user's requested format
	req.Header.Set("Content-Type", (yamlFormat{}).ContentType())

	// final request looks like this:
	// PUT {server}/{resourceNamePlural}
	// Content-Type: text/yaml
	// Accept: {requestedFormat.contentType}
	//
	// {yamlFileContent}
	//
	// This means that we'll send the request body as YAML (read from the user provided file)
	// and we'll get the reponse in the users's requrested format.
	d, _ := httputil.DumpRequestOut(req, true)
	c.logger.Debug("apply request",
		zap.String("request", string(d)),
	)

	resp, err := c.client.do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute Apply request: %w", err)
	}
	defer resp.Body.Close()

	d, _ = httputil.DumpResponse(resp, true)
	c.logger.Debug("apply response",
		zap.String("response", string(d)),
	)

	if !isSuccessResponse(resp) {
		err := parseRequestError(resp, requestedFormat)

		return "", fmt.Errorf("could not Apply resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read Apply response: %w", err)
	}

	c.logger.Debug("file has id?", zap.Bool("hasID", originalInputFile.HasID()))
	// if the original file doesn't have an ID, we need to get the server generated ID from the response
	// and write it to the original file
	if !originalInputFile.HasID() {
		jsonBody, err := requestedFormat.ToJSON(body)
		if err != nil {
			return "", fmt.Errorf("cannot convert response body to JSON format: %w", err)
		}

		parsed, err := gabs.ParseJSON(jsonBody)
		if err != nil {
			return "", fmt.Errorf("cannot parse Apply response: %w", err)
		}

		id, ok := parsed.Path("spec.id").Data().(string)
		if !ok {
			return "", fmt.Errorf("cannot get ID from Apply response")
		}

		c.logger.Debug("New ID", zap.String("id", id))

		originalInputFile, err = originalInputFile.SetID(id)
		if err != nil {
			return "", fmt.Errorf("cannot set ID on input file: %w", err)
		}

		_, err = originalInputFile.Write()
		if err != nil {
			return "", fmt.Errorf("cannot write updated input file: %w", err)
		}
	}

	return requestedFormat.Format(string(body), c.options.tableConfig, c.options.listPath)
}
