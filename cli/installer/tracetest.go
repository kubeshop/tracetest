package installer

import (
	"bytes"
	_ "embed"
	"html/template"

	"fmt"

	cliUI "github.com/kubeshop/tracetest/cli/ui"
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

	if !installBackend {
		conf.set("tracetest.backend.type", "")
		return conf
	}

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

//go:embed templates/config.yaml.tpl
var configTemplate string

func getTracetestConfigFileContents(pHost, pUser, pPasswd string, ui cliUI.UI, config configuration) []byte {
	vals := map[string]string{
		"pHost":   pHost,
		"pUser":   pUser,
		"pPasswd": pPasswd,
	}

	tpl, err := template.New("page").Parse(configTemplate)
	if err != nil {
		ui.Panic(fmt.Errorf("cannot parse config template: %w", err))
	}

	out := &bytes.Buffer{}
	tpl.Execute(out, vals)

	return out.Bytes()
}

//go:embed templates/provision.yaml.tpl
var provisionTemplate string

func getTracetestProvisionFileContents(ui cliUI.UI, config configuration) []byte {
	vals := map[string]string{
		"installBackend":   fmt.Sprintf("%t", config.Bool("tracetest.backend.install")),
		"backendType":      config.String("tracetest.backend.type"),
		"backendEndpoint":  config.String("tracetest.backend.endpoint.query"),
		"backendInsecure":  config.String("tracetest.backend.tls.insecure"),
		"backendAddresses": config.String("tracetest.backend.addresses"),
		"backendIndex":     config.String("tracetest.backend.index"),
		"backendToken":     config.String("tracetest.backend.token"),
		"backendRealm":     config.String("tracetest.backend.realm"),

		"analyticsEnabled": fmt.Sprintf("%t", config.Bool("tracetest.analytics")),

		"enablePokeshopDemo": fmt.Sprintf("%t", config.Bool("demo.enable.pokeshop")),
		"enableOtelDemo":     fmt.Sprintf("%t", config.Bool("demo.enable.otel")),
		"pokeshopHttp":       config.String("demo.endpoint.pokeshop.http"),
		"pokeshopGrpc":       config.String("demo.endpoint.pokeshop.grpc"),
		"otelFrontend":       config.String("demo.endpoint.otel.frontend"),
		"otelProductCatalog": config.String("demo.endpoint.otel.product_catalog"),
		"otelCart":           config.String("demo.endpoint.otel.cart"),
		"otelCheckout":       config.String("demo.endpoint.otel.checkout"),
	}

	tpl, err := template.New("page").Parse(provisionTemplate)
	if err != nil {
		ui.Panic(fmt.Errorf("cannot parse config template: %w", err))
	}

	out := &bytes.Buffer{}
	tpl.Execute(out, vals)

	return out.Bytes()
}
