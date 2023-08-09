package datasource

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configcompression"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	name     string
	config   *configgrpc.GRPCClientSettings
	conn     *grpc.ClientConn
	callback GrpcCallback
}

func convertDomainConfigToOpenTelemetryConfig(config *datastore.GRPCClientSettings) *configgrpc.GRPCClientSettings {
	// manual map domain fields to OTel fields
	otelConfig := &configgrpc.GRPCClientSettings{
		Endpoint:        config.Endpoint,
		ReadBufferSize:  config.ReadBufferSize,
		WriteBufferSize: config.WriteBufferSize,
		WaitForReady:    config.WaitForReady,
		Headers:         config.Headers,
		BalancerName:    config.BalancerName,

		Compression: configcompression.CompressionType(config.Compression),
	}

	if config.TLS == nil {
		return otelConfig
	}

	otelConfig.TLSSetting = configtls.TLSClientSetting{
		Insecure:           config.TLS.Insecure,
		InsecureSkipVerify: config.TLS.InsecureSkipVerify,
		ServerName:         config.TLS.ServerName,
	}

	if config.TLS.Settings == nil {
		return otelConfig
	}

	otelConfig.TLSSetting.TLSSetting = configtls.TLSSetting{
		CAFile:     config.TLS.Settings.CAFile,
		CertFile:   config.TLS.Settings.CertFile,
		KeyFile:    config.TLS.Settings.KeyFile,
		MinVersion: config.TLS.Settings.MinVersion,
		MaxVersion: config.TLS.Settings.MaxVersion,
	}

	return otelConfig
}

func NewGrpcClient(name string, config *datastore.GRPCClientSettings, callback GrpcCallback) DataSource {
	otelConfig := convertDomainConfigToOpenTelemetryConfig(config)

	return &GrpcClient{
		name:     name,
		config:   otelConfig,
		callback: callback,
	}
}

func (client *GrpcClient) Ready() bool {
	return client.conn != nil
}

func (client *GrpcClient) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	return client.callback(ctx, traceID, client.conn)
}

func (client *GrpcClient) Endpoint() string {
	return client.config.Endpoint
}

func (client *GrpcClient) Connect(ctx context.Context) error {
	opts, err := client.config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return errors.Wrap(connection.ErrInvalidConfiguration, err.Error())
	}

	conn, err := grpc.DialContext(ctx, client.config.Endpoint, opts...)
	if err != nil {
		return errors.Wrap(connection.ErrConnectionFailed, err.Error())
	}

	client.conn = conn
	return nil
}

func (client *GrpcClient) TestConnection(ctx context.Context) model.ConnectionTestStep {
	connectionTestResult := model.ConnectionTestStep{
		Message: fmt.Sprintf(`Tracetest connected to "%s"`, client.config.Endpoint),
	}

	err := connection.CheckReachability(client.config.Endpoint, model.ProtocolGRPC)
	if err != nil {
		return model.ConnectionTestStep{
			Message: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, client.config.Endpoint),
			Error:   err,
		}
	}

	err = client.Connect(ctx)
	wrappedErr := errors.Unwrap(err)
	if errors.Is(wrappedErr, connection.ErrConnectionFailed) {
		return model.ConnectionTestStep{
			Message: fmt.Sprintf(`Tracetest tried to open a gRPC connection against "%s" and failed`, client.config.Endpoint),
			Error:   err,
		}
	}

	return connectionTestResult
}

func (client *GrpcClient) Close() error {
	err := client.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}

	return nil
}
