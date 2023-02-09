package installer

import (
	"fmt"
	"strings"

	cliUI "github.com/kubeshop/tracetest/cli/ui"
	serverConfig "github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"gopkg.in/yaml.v3"
)

func configureDemoApp(conf configuration, ui cliUI.UI) configuration {
	conf.set("demo.enable.pokeshop", !conf.Bool("installer.only_tracetest"))
	conf.set("demo.enable.otel", false)

	switch conf.String("installer") {
	case "docker-compose":
		conf.set("demo.endpoint.pokeshop.http", "http://demo-api:8081")
		conf.set("demo.endpoint.pokeshop.grpc", "demo-api:8082")
		conf.set("demo.endpoint.otel.frontend", "http://otel-frontend:8084")
		conf.set("demo.endpoint.otel.product_catalog", "otel-productcatalogservice:3550")
		conf.set("demo.endpoint.otel.cart", "otel-cartservice:7070")
		conf.set("demo.endpoint.otel.checkout", "otel-checkoutservice:5050")
	case "kubernetes":
		conf.set("demo.endpoint.pokeshop.http", "http://demo-pokemon-api.demo")
		conf.set("demo.endpoint.pokeshop.grpc", "demo-pokemon-api.demo:8082")
		conf.set("demo.endpoint.otel.frontend", "http://otel-frontend.otel-demo:8084")
		conf.set("demo.endpoint.otel.product_catalog", "otel-productcatalogservice.otel-demo:3550")
		conf.set("demo.endpoint.otel.cart", "otel-cartservice.otel-demo:7070")
		conf.set("demo.endpoint.otel.checkout", "otel-checkoutservice.otel-demo:5050")
	}

	return conf
}

func configureTracetest(conf configuration, ui cliUI.UI) configuration {
	conf = configureBackend(conf, ui)
	conf.set("tracetest.analytics", true)

	return conf
}

func configureBackend(conf configuration, ui cliUI.UI) configuration {
	installBackend := !conf.Bool("installer.only_tracetest")
	conf.set("tracetest.backend.install", installBackend)

	if installBackend {

		// default values
		switch conf.String("installer") {
		case "docker-compose":
			conf.set("tracetest.backend.type", "otlp")
			conf.set("tracetest.backend.tls.insecure", true)
			conf.set("tracetest.backend.endpoint.collector", "http://otel-collector:4317")
			conf.set("tracetest.backend.endpoint", "tracetest:21321")
		case "kubernetes":
			conf.set("tracetest.backend.type", "otlp")
			conf.set("tracetest.backend.tls.insecure", true)
			conf.set("tracetest.backend.endpoint.collector", "http://otel-collector.tracetest:4317")
			conf.set("tracetest.backend.endpoint", "tracetest:21321")

		default:
			conf.set("tracetest.backend.type", "")
		}

		return conf
	}

	conf.set("tracetest.backend.type", "")
	return conf
}

func getTracetestConfigFileContents(psql string, ui cliUI.UI, config configuration) []byte {
	sc := serverConfig.Config{
		PostgresConnString: psql,
		PoolingConfig: serverConfig.PoolingConfig{
			MaxWaitTimeForTrace: "2m",
			RetryDelay:          "3s",
		},
		GA: serverConfig.GoogleAnalytics{
			Enabled: config.Bool("tracetest.analytics"),
		},
	}

	sc.Telemetry = telemetryConfig(ui, config)

	if config.Bool("tracetest.backend.install") {
		sc.Server = serverConfig.ServerConfig{
			Telemetry: serverConfig.ServerTelemetryConfig{
				DataStore: config.String("tracetest.backend.type"),
			},
		}
	}

	enabledDemos := []string{}
	if config.Bool("demo.enable.pokeshop") {
		enabledDemos = append(enabledDemos, "pokeshop")
	}
	if config.Bool("demo.enable.otel") {
		enabledDemos = append(enabledDemos, "otel")
	}

	sc.Demo = serverConfig.Demo{
		Enabled: enabledDemos,
		Endpoints: serverConfig.DemoEndpoints{
			PokeshopHttp:       config.String("demo.endpoint.pokeshop.http"),
			PokeshopGrpc:       config.String("demo.endpoint.pokeshop.grpc"),
			OtelFrontend:       config.String("demo.endpoint.otel.frontend"),
			OtelProductCatalog: config.String("demo.endpoint.otel.product_catalog"),
			OtelCart:           config.String("demo.endpoint.otel.cart"),
			OtelCheckout:       config.String("demo.endpoint.otel.checkout"),
		},
	}

	out, err := yaml.Marshal(sc)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot marshal tracetest config file: %w", err).Error())
	}

	if config.Bool("tracetest.backend.install") {
		out, err = fixConfigs(out)
		if err != nil {
			ui.Exit(fmt.Errorf("cannot fix tracertest config: %w", err).Error())
		}
	}

	return out
}

func fixConfigs(conf []byte) ([]byte, error) {
	encoded := msa{}

	err := yaml.Unmarshal(conf, &encoded)
	if err != nil {
		return nil, fmt.Errorf("cannot decode mapstructure: %w", err)
	}

	ds := encoded["telemetry"].(msa)["datastores"].(msa)

	var (
		target msa
		key    string
	)

	if d, ok := ds["jaeger"]; ok {
		target = d.(msa)["jaeger"].(msa)
		key = "jaeger"
	} else if d, ok := ds["tempo"]; ok {
		target = d.(msa)["tempo"].(msa)
		key = "tempo"
	} else {
		// we only need to fix tempo/jaeger
		return yaml.Marshal(encoded)
	}

	target["tls"] = target["tlssetting"]
	delete(target, "tlssetting")

	encoded["telemetry"].(msa)["datastores"].(msa)[key].(msa)[key] = target

	return yaml.Marshal(encoded)
}

func telemetryConfig(ui cliUI.UI, conf configuration) serverConfig.Telemetry {
	if conf.Bool("tracetest.backend.install") {
		return serverConfig.Telemetry{
			DataStores: dataStoreConfig(ui, conf),
			Exporters:  map[string]serverConfig.TelemetryExporterOption{},
		}
	}

	return serverConfig.Telemetry{
		DataStores: map[string]serverConfig.TracingBackendDataStoreConfig{},
		Exporters:  map[string]serverConfig.TelemetryExporterOption{},
	}
}

func dataStoreConfig(ui cliUI.UI, conf configuration) map[string]serverConfig.TracingBackendDataStoreConfig {
	dstype := conf.String("tracetest.backend.type")
	var c serverConfig.TracingBackendDataStoreConfig
	switch dstype {
	case "jaeger":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			Jaeger: configgrpc.GRPCClientSettings{
				Endpoint: conf.String("tracetest.backend.endpoint.query"),
				TLSSetting: configtls.TLSClientSetting{
					Insecure: conf.Bool("tracetest.backend.tls.insecure"),
				},
			},
		}
	case "tempo":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			Tempo: serverConfig.BaseClientConfig{
				Grpc: configgrpc.GRPCClientSettings{
					Endpoint: conf.String("tracetest.backend.endpoint"),
					TLSSetting: configtls.TLSClientSetting{
						Insecure: conf.Bool("tracetest.backend.tls.insecure"),
					},
				},
				Type: "grpc",
			},
		}
	case "opensearch":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			OpenSearch: serverConfig.ElasticSearchDataStoreConfig{
				Addresses: strings.Split(conf.String("tracetest.backend.addresses"), ","),
				Index:     conf.String("tracetest.backend.index"),
			},
		}
	case "signalfx":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			SignalFX: serverConfig.SignalFXDataStoreConfig{
				Token: conf.String("tracetest.backend.token"),
				Realm: conf.String("tracetest.backend.realm"),
			},
		}
	case "otlp":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
		}
	default:
		ui.Panic(fmt.Errorf("unsupported dataStore type %s", dstype))
	}

	return map[string]serverConfig.TracingBackendDataStoreConfig{
		dstype: c,
	}

}
