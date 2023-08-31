package datasource

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/goware/urlx"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
)

type HttpClient struct {
	name     string
	config   *http.Request
	client   *http.Client
	callback HttpCallback
}

func NewHttpClient(name string, config *datastore.HttpClientConfig, callback HttpCallback) DataSource {
	endpoint, _ := urlx.Parse(config.Url)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: getTlsConfig(config.TLS),
		},
	}

	header := http.Header{}
	for key, value := range config.Headers {
		header.Set(key, value)
	}

	request := &http.Request{
		URL:    endpoint,
		Header: header,
	}

	return &HttpClient{
		name:     name,
		config:   request,
		client:   client,
		callback: callback,
	}
}

func (client *HttpClient) Ready() bool {
	return client.callback != nil
}

func (client *HttpClient) Close() error {
	return nil
}

func (client *HttpClient) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	return client.callback(ctx, traceID, client)
}

func (client *HttpClient) Endpoint() string {
	return client.config.URL.String()
}

func (client *HttpClient) Connect(ctx context.Context) error {
	_, err := client.client.Transport.RoundTrip(client.config)

	return err
}

func (client *HttpClient) TestConnection(ctx context.Context) model.ConnectionTestStep {
	connectionTestResult := model.ConnectionTestStep{
		Message: fmt.Sprintf(`Tracetest connected to "%s"`, client.config.URL.String()),
	}

	err := connection.CheckReachability(client.config.URL.String(), model.ProtocolHTTP)
	if err != nil {
		return model.ConnectionTestStep{
			Message: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, client.config.URL.String()),
			Error:   err,
		}
	}

	err = client.Connect(ctx)
	wrappedErr := errors.Unwrap(err)
	if errors.Is(wrappedErr, connection.ErrConnectionFailed) {
		return model.ConnectionTestStep{
			Message: fmt.Sprintf(`Tracetest tried to open a connection against "%s" and failed`, client.config.URL.String()),
			Error:   err,
		}
	}

	return connectionTestResult
}

func (client *HttpClient) Request(ctx context.Context, path, method, body string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", client.config.URL.String(), path)
	var readerBody io.Reader
	if body != "" {
		readerBody = bytes.NewBufferString(body)
	}

	request, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), url, readerBody)
	if err != nil {
		return nil, err
	}

	request.Header = client.config.Header
	response, err := client.client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func getTlsConfig(dataStoreTls *datastore.TLS) *tls.Config {
	tlsConfig := tls.Config{}

	if dataStoreTls == nil {
		return &tlsConfig
	}

	if dataStoreTls.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	if dataStoreTls.Settings == nil {
		return &tlsConfig
	}

	caCertFile := dataStoreTls.Settings.CAFile

	if caCertFile != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCertFile))
		tlsConfig.RootCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}

	return &tlsConfig
}
