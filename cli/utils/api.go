package utils

import (
	"context"
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
	BaseHeader   http.Header
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
	baseHeader := http.Header{
		"x-client-id":  []string{analytics.ClientID()},
		"Content-Type": []string{"text/yaml"},
	}

	return ResourceClient{
		Client:       client,
		BaseUrl:      baseUrl,
		BaseHeader:   baseHeader,
		ResourceType: resourceType,
	}
}

func (resourceClient ResourceClient) NewRequest(url string, method string, body string) (*http.Request, error) {
	var reqBody io.Reader
	if body != "" {
		reqBody = StringToIOReader(body)
	}

	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	request.Header = resourceClient.BaseHeader
	return request, err
}

func (resourceClient ResourceClient) Update(ctx context.Context, file file.File, ID string) (*file.File, error) {
	url := fmt.Sprintf("%s/%s", resourceClient.BaseUrl, ID)
	request, err := resourceClient.NewRequest(url, http.MethodPut, file.Contents())
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

	file = file.SaveChanges(IOReadCloserToString(resp.Body))

	return &file, nil
}

func (resourceClient ResourceClient) Delete(ctx context.Context, ID string) error {
	url := fmt.Sprintf("%s/%s", resourceClient.BaseUrl, ID)
	request, err := resourceClient.NewRequest(url, http.MethodDelete, "")
	if err != nil {
		return fmt.Errorf("could not delete resource: %w", err)
	}

	_, err = resourceClient.Client.Do(request)
	return err
}

func (resourceClient ResourceClient) Get(ctx context.Context, id string) (*file.File, error) {
	request, err := resourceClient.NewRequest(fmt.Sprintf("%s/%s", resourceClient.BaseUrl, id), http.MethodGet, "")
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

func (resourceClient ResourceClient) List(ctx context.Context, listArgs ListArgs) (*file.File, error) {
	url := fmt.Sprintf("%s?skip=%d&take=%d&sortBy=%s&sortDirection=%s", resourceClient.BaseUrl, listArgs.Skip, listArgs.Take, listArgs.SortBy, listArgs.SortDirection)
	request, err := resourceClient.NewRequest(url, http.MethodGet, "")
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := resourceClient.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	file, err := file.NewFromRaw(fmt.Sprintf("%s_list.yaml", resourceClient.ResourceType), []byte(IOReadCloserToString(resp.Body)))
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (resourceClient ResourceClient) Create(ctx context.Context, file file.File) (*file.File, error) {
	request, err := resourceClient.NewRequest(resourceClient.BaseUrl, http.MethodPost, file.Contents())
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
