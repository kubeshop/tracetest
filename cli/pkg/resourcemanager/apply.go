package resourcemanager

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Jeffail/gabs/v2"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
)

const VerbApply Verb = "apply"

type applyPreProcessorFn func(context.Context, fileutil.File) (fileutil.File, error)

func (c Client) Apply(ctx context.Context, inputFile fileutil.File, requestedFormat Format) (string, error) {
	originalInputFile := inputFile

	if c.options.applyPreProcessor != nil {
		var err error
		inputFile, err = c.options.applyPreProcessor(ctx, inputFile)
		if err != nil {
			return "", fmt.Errorf("cannot preprocess Apply request: %w", err)
		}
	}

	url := c.client.url(c.resourceNamePlural)
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
	yamlFormat, err := Formats.Get(FormatYAML)
	if err != nil {
		return "", fmt.Errorf("cannot get json format: %w", err)
	}
	req.Header.Set("Content-Type", yamlFormat.ContentType())

	// final request looks like this:
	// PUT {server}/{resourceNamePlural}
	// Content-Type: text/yaml
	// Accept: {requestedFormat.contentType}
	//
	// {yamlFileContent}
	//
	// This means that we'll send the request body as YAML (read from the user provided file)
	// and we'll get the reponse in the users's requrested format.
	resp, err := c.client.do(req)
	if err != nil {
		return "", fmt.Errorf("cannot execute Apply request: %w", err)
	}
	defer resp.Body.Close()

	if !isSuccessResponse(resp) {
		err := parseRequestError(resp, requestedFormat)

		return "", fmt.Errorf("could not Apply resource: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read Apply response: %w", err)
	}

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

		originalInputFile, err = originalInputFile.SetID(id)
		if err != nil {
			return "", fmt.Errorf("cannot set ID on input file: %w", err)
		}

		_, err = originalInputFile.Write()
		if err != nil {
			return "", fmt.Errorf("cannot write updated input file: %w", err)
		}

	}

	return requestedFormat.Format(string(body), c.options.tableConfig)
}
