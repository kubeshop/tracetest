package installer

import (
	"fmt"
	"strings"

	serverConfig "github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"gopkg.in/yaml.v3"
)

func configureDemoApp(conf configuration, ui UI) configuration {
	conf.set(
		"demo.enable",
		ui.Confirm("Do you want to enable the demo app?", true),
	)
	return conf
}

func configureTracetest(conf configuration, ui UI) configuration {
	conf = configureBackend(conf, ui)
	conf = configureCollector(conf, ui)

	conf.set(
		"tracetest.analytics",
		ui.Confirm("Do you want to help improve tracetest by providing anonymous analytics information?", true),
	)

	return conf
}
func configureCollector(conf configuration, ui UI) configuration {
	installCollector := false

	hasCollector := ui.Confirm("Do you have an OpenTelemetry Collector?", false)
	if hasCollector {
		conf.set("tracetest.collector.endpoint", ui.TextInput("Endpoint", "otel-collector:4317"))
	} else {
		if !ui.Confirm("Do you want me to set up one?", true) {
			ui.Exit(`
TraceTest requires OpenTelemetry Collector to work. You can rerun this installer and let me set it up for you,
or you can set one up manually. See https://opentelemetry.io/docs/collector/
`)
		}
		installCollector = true

		// default values
		conf.set("tracetest.collector.endpoint", "otel-collector:4317")
	}

	conf.set("tracetest.collector.install", installCollector)
	return conf
}

func configureBackend(conf configuration, ui UI) configuration {
	installBackend := false

	hasBackend := ui.Confirm("Do you have a supported tracing backend you want to use? (Jaeger, Tempo, OpenSearch, SignalFX)", false)
	if hasBackend {
		conf = configureBackendOptions(conf, ui)
	} else {
		if !ui.Confirm("Do you want me to set up Jaeger?", true) {
			ui.Exit(`
TraceTest requires a supported tracing backend to work. I only know how to install Jaeger.
If you want to use other option, check the supported backends and manually install one.
See https://kubeshop.github.io/tracetest/supported-backends/
			`)
		}
		installBackend = true

		// default values
		conf.set("tracetest.backend.type", "jaeger")
		conf.set("tracetest.backend.endpoint", "jaeger:16685")
		conf.set("tracetest.backend.tls.insecure", true)
	}

	conf.set("tracetest.backend.install", installBackend)

	return conf
}

func configureBackendOptions(conf configuration, ui UI) configuration {
	option := ui.Select("Which tracing backend do you want to use?", []option{
		{"Jaeger", func(ui UI) {
			conf.set("tracetest.backend.type", "jaeger")
			conf.set("tracetest.backend.endpoint", ui.TextInput("Endpoint", "jaeger-query:16685"))
			conf.set("tracetest.backend.tls.insecure", ui.Confirm("TLS/Insecure", true))
		}},
		{"Tempo", func(ui UI) {
			conf.set("tracetest.backend.type", "tempo")
			conf.set("tracetest.backend.endpoint", ui.TextInput("Endpoint", "tempo:9095"))
			conf.set("tracetest.backend.tls.insecure", ui.Confirm("Insecure", true))
		}},
		{"OpenSearch", func(ui UI) {
			conf.set("tracetest.backend.type", "opensearch")
			conf.set("tracetest.backend.addresses", ui.TextInput("Addresses (comma separated list)", "http://opensearch:9200"))
			conf.set("tracetest.backend.index", ui.TextInput("Index", "traces"))
			conf.set("tracetest.backend.data-prepper", ui.TextInput("Data Prepper", "data-prepper:21890"))
		}},
		{"SignalFX", func(ui UI) {
			conf.set("tracetest.backend.type", "signalfx")
			conf.set("tracetest.backend.token", ui.TextInput("Token", ""))
			conf.set("tracetest.backend.realm", ui.TextInput("Realm", "us1"))
		}},
	})

	option.fn(ui)

	return conf
}

func getTracetestConfigFileContents(ui UI, config configuration) []byte {
	sc := serverConfig.Config{
		PostgresConnString: "host=postgres user=postgres password=postgres port=5432 sslmode=disable",
		PoolingConfig: serverConfig.PoolingConfig{
			MaxWaitTimeForTrace: "2m",
			RetryDelay:          "1s",
		},
		GA: serverConfig.GoogleAnalytics{
			Enabled: config.Bool("tracetest.analytics"),
		},
	}

	sc.Telemetry = telemetryConfig(ui, config)
	sc.Server = serverConfig.ServerConfig{
		Telemetry: serverConfig.ServerTelemetryConfig{
			Exporter:            "collector",
			ApplicationExporter: "collector",
			DataStore:           config.String("tracetest.backend.type"),
		},
	}

	out, err := yaml.Marshal(sc)
	if err != nil {
		ui.Exit(fmt.Errorf("Cannot marshal tracetest config file: %w", err).Error())
	}

	return out
}

func telemetryConfig(ui UI, conf configuration) serverConfig.Telemetry {
	return serverConfig.Telemetry{
		DataStores: dataStoreConfig(ui, conf),
		Exporters:  exportersConfig(ui, conf),
	}
}

func exportersConfig(ui UI, conf configuration) map[string]serverConfig.TelemetryExporterOption {
	return map[string]serverConfig.TelemetryExporterOption{
		"collector": {
			ServiceName: "tracetest",
			Sampling:    100,
			Exporter: serverConfig.ExporterConfig{
				Type: "collector",
				CollectorConfiguration: serverConfig.OTELCollectorConfig{
					Endpoint: conf.String("tracetest.collector.endpoint"),
				},
			},
		},
	}

}
func dataStoreConfig(ui UI, conf configuration) map[string]serverConfig.TracingBackendDataStoreConfig {
	dstype := conf.String("tracetest.backend.type")
	var c serverConfig.TracingBackendDataStoreConfig
	switch dstype {
	case "jaeger":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			Jaeger: configgrpc.GRPCClientSettings{
				Endpoint: conf.String("tracetest.backend.endpoint"),
				TLSSetting: configtls.TLSClientSetting{
					Insecure: conf.Bool("tracetest.backend.tls.insecure"),
				},
			},
		}
	case "tempo":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			Tempo: configgrpc.GRPCClientSettings{
				Endpoint: conf.String("tracetest.backend.endpoint"),
				TLSSetting: configtls.TLSClientSetting{
					Insecure: conf.Bool("tracetest.backend.tls.insecure"),
				},
			},
		}
	case "opensearch":
		c = serverConfig.TracingBackendDataStoreConfig{
			Type: dstype,
			OpenSearch: serverConfig.OpensearchDataStoreConfig{
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
	default:
		ui.Panic(fmt.Errorf("unsupported dataStore type %s", dstype))
	}

	return map[string]serverConfig.TracingBackendDataStoreConfig{
		dstype: c,
	}

}
