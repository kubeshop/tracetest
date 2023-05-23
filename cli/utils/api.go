package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type ListArgs struct {
	Take          int32
	Skip          int32
	SortDirection string
	SortBy        string
	All           bool
}

func GetAPIClient(cliConfig config.Config) *openapi.APIClient {
	config := openapi.NewConfiguration()
	config.AddDefaultHeader("x-client-id", analytics.ClientID())
	config.Scheme = cliConfig.Scheme
	config.Host = strings.TrimSuffix(cliConfig.Endpoint, "/")
	if cliConfig.ServerPath != nil {
		config.Servers = []openapi.ServerConfiguration{
			{
				URL: *cliConfig.ServerPath,
			},
		}
	}

	return openapi.NewAPIClient(config)
}

type ResourceClient struct {
	Client       http.Client
	BaseUrl      string
	ResourceType string
}

func GetResourceAPIClient(resourceType string, cliConfig config.Config) ResourceClient {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second,
		}).DialContext,
	}

	client := http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	baseUrl := fmt.Sprintf("%s://%s/api/%s", cliConfig.Scheme, cliConfig.Endpoint, resourceType)

	return ResourceClient{
		Client:       client,
		BaseUrl:      baseUrl,
		ResourceType: resourceType,
	}
}

func (resourceClient ResourceClient) NewRequest(url, method, body, contentType string, augmented bool) (*http.Request, error) {
	var reqBody io.Reader
	if body != "" {
		reqBody = StringToIOReader(body)
	}

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if contentType == "" {
		contentType = "application/json"
	}

	request.Header.Set("x-client-id", analytics.ClientID())
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", contentType)
	if augmented {
		request.Header.Set("X-Tracetest-Augmented", "true")
	}

	return request, err
}

func (resourceClient ResourceClient) Update(ctx context.Context, file file.File, ID string) (*file.File, error) {
	url := fmt.Sprintf("%s/%s", resourceClient.BaseUrl, ID)
	request, err := resourceClient.NewRequest(url, http.MethodPut, file.Contents(), file.ContentType(), false)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := resourceClient.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not update %s: %w", resourceClient.ResourceType, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%s id doesn't exist on server. Remove it from the definition file and try again", resourceClient.ResourceType)
	}

	if resp.StatusCode == http.StatusUnprocessableEntity || resp.StatusCode == http.StatusBadRequest {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("could not send request: %w", err)
		}

		validationError := string(body)
		return nil, fmt.Errorf("invalid %s: %s", resourceClient.ResourceType, validationError)
	}

	responseContentType := resp.Header.Get("Content-type")
	if responseContentType == "" {
		responseContentType = "application/json"
	}

	file = file.SaveChanges(IOReadCloserToString(resp.Body))

	return &file, nil
}

func (resourceClient ResourceClient) Delete(ctx context.Context, ID string) error {
	url := fmt.Sprintf("%s/%s", resourceClient.BaseUrl, ID)
	request, err := resourceClient.NewRequest(url, http.MethodDelete, "", "", false)
	if err != nil {
		return fmt.Errorf("could not delete resource: %w", err)
	}

	_, err = resourceClient.Client.Do(request)
	return err
}

func (resourceClient ResourceClient) Get(ctx context.Context, id string) (*file.File, error) {
	augmented := false
	if ctx.Value("X-Tracetest-Augmented") != nil {
		augmented = true
	}

	request, err := resourceClient.NewRequest(fmt.Sprintf("%s/%s", resourceClient.BaseUrl, id), http.MethodGet, "", "", augmented)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := resourceClient.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not get %s: %w", resourceClient.ResourceType, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		validationError := string(body)
		return nil, fmt.Errorf("invalid %s: %s", resourceClient.ResourceType, validationError)
	}

	file, err := file.NewFromRaw(fmt.Sprintf("%s_%s.yaml", resourceClient.ResourceType, id), []byte(IOReadCloserToString(resp.Body)))
	if err != nil {
		return nil, err
	}

	return &file, nil
}

type BaseListResponse struct {
	Count int           `json:"count"`
	Items []interface{} `json:"items"`
}

func parseListResponse(body string) (BaseListResponse, error) {
	var response BaseListResponse

	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		return BaseListResponse{}, fmt.Errorf("could not parse response: %w", err)
	}

	return response, nil
}

func (resourceClient ResourceClient) List(ctx context.Context, listArgs ListArgs) (*file.File, error) {
	augmented := false
	if ctx.Value("X-Tracetest-Augmented") != nil {
		augmented = true
	}

	url := fmt.Sprintf("%s?skip=%d&take=%d&sortBy=%s&sortDirection=%s", resourceClient.BaseUrl, listArgs.Skip, listArgs.Take, listArgs.SortBy, listArgs.SortDirection)
	request, err := resourceClient.NewRequest(url, http.MethodGet, "", "", augmented)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := resourceClient.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	body := IOReadCloserToString(resp.Body)
	if listArgs.All {
		baseListResponse, err := parseListResponse(body)
		if err != nil {
			return nil, err
		}

		if baseListResponse.Count > len(baseListResponse.Items) {
			return resourceClient.List(ctx, ListArgs{
				Skip:          0,
				Take:          int32(baseListResponse.Count) + 1,
				SortBy:        listArgs.SortBy,
				SortDirection: listArgs.SortDirection,
				All:           false,
			})
		}
	}

	file, err := file.NewFromRaw(fmt.Sprintf("%s_list.yaml", resourceClient.ResourceType), []byte(body))
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (resourceClient ResourceClient) Create(ctx context.Context, file file.File) (*file.File, error) {
	request, err := resourceClient.NewRequest(resourceClient.BaseUrl, http.MethodPost, file.Contents(), file.ContentType(), false)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := resourceClient.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("could not send request: %w", err)
		}

		validationError := string(body)
		return nil, fmt.Errorf("invalid %s: %s", resourceClient.ResourceType, validationError)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", resourceClient.ResourceType, err)
	}

	file = file.SaveChanges(IOReadCloserToString(resp.Body))
	return &file, nil
}
