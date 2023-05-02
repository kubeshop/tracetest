package installer

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/analytics"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

var (
	Force            = false
	RunEnvironment   = NoneRunEnvironmentType
	InstallationMode = NotChosenInstallationModeType
)

const createIssueMsg = "If you need help, please create an issue: https://github.com/kubeshop/tracetest/issues/new/choose"

func Start() {
	analytics.Track("Start", "installer", map[string]string{})
	ui := cliUI.DefaultUI

	ui.Banner()

	ui.Println(`
Hi! Welcome to the TraceTest server installer. I'll help you set up your TraceTest server by asking you a few questions
and configuring your system with all the requirements, so you can start TraceTesting right away!

To get more info about TraceTest, you can check our docs at https://kubeshop.github.io/tracetest/

If you have any issues, please let us know by creating an issue (https://github.com/kubeshop/tracetest/issues/new/choose)
or reach us on Discord https://discord.gg/6zupCZFQbe

`)

	if RunEnvironment == DockerRunEnvironmentType { // check if docker was previously chosen as a CLI arg
		ui.Println("How do you want to run TraceTest?")
		ui.Println("  > Using Docker Compose")
		dockerCompose.Install(ui)
		return
	}

	if RunEnvironment == KubernetesRunEnvironmentType { // check if kubernetes was previously chosen as a CLI arg
		ui.Println("How do you want to run TraceTest?")
		ui.Println("  > Using Kubernetes")
		kubernetes.Install(ui)
		return
	}

	option := ui.Select("How do you want to run TraceTest?", []cliUI.Option{
		{Text: "Using Docker Compose", Fn: dockerCompose.Install},
		{Text: "Using Kubernetes", Fn: kubernetes.Install},
	}, 0)

	option.Fn(ui)
}

type installer struct {
	name      string
	preChecks []preChecker
	configs   []configurator
	installFn func(config configuration, ui cliUI.UI)
}

func (i installer) PreCheck(ui cliUI.UI) {
	ui.Title("Let's check if your system has everything we need")
	for _, pc := range i.preChecks {
		pc(ui)
	}

	ui.Title("Your system is ready! Now, let's configure TraceTest")
}

func (i installer) Configure(ui cliUI.UI) configuration {
	config := newConfiguration(ui)
	config.set("installer", i.name)

	setInstallationType(ui, config)

	for _, confFn := range i.configs {
		config = confFn(config, ui)
	}

	return config
}

func (i installer) Install(ui cliUI.UI) {
	analytics.Track("PreCheck", "installer", map[string]string{})
	i.PreCheck(ui)

	analytics.Track("Configure", "installer", map[string]string{})
	conf := i.Configure(ui)

	ui.Title("Thanks! We are ready to install TraceTest now")

	i.installFn(conf, ui)
}

type preChecker func(ui cliUI.UI)

func setInstallationType(ui cliUI.UI, config configuration) {
	if InstallationMode == WithoutDemoInstallationModeType { // check if it was previously chosen
		ui.Println("Do you have OpenTelemetry based tracing already set up, or would you like us to install a demo tracing environment and app?")
		ui.Println("  > I have a tracing environment already. Just install Tracetest")
		config.set("installer.only_tracetest", true)
		return
	}

	if InstallationMode == WithDemoInstallationModeType { // check if it was previously chosen
		ui.Println("Do you have OpenTelemetry based tracing already set up, or would you like us to install a demo tracing environment and app?")
		ui.Println("  > Just learning tracing! Install Tracetest, OpenTelemetry Collector and the sample app.")
		config.set("installer.only_tracetest", false)
		return
	}

	option := ui.Select("Do you have OpenTelemetry based tracing already set up, or would you like us to install a demo tracing environment and app?", []cliUI.Option{
		{Text: "I have a tracing environment already. Just install Tracetest", Fn: func(ui cliUI.UI) {
			config.set("installer.only_tracetest", true)
		}},
		{Text: "Just learning tracing! Install Tracetest, OpenTelemetry Collector and the sample app.", Fn: func(ui cliUI.UI) {
			config.set("installer.only_tracetest", false)
		}},
	}, 0)

	option.Fn(ui)
}

func trackInstall(name string, config configuration, extra map[string]string) {
	props := map[string]string{
		"type":                    name,
		"install_backend":         fmt.Sprintf("%t", config.Bool("tracetest.backend.install")),
		"install_demo_pokeshop":   fmt.Sprintf("%t", config.Bool("demo.enable.pokeshop")),
		"install_demo_otel":       fmt.Sprintf("%t", config.Bool("demo.enable.otel")),
		"enable_server_analytics": fmt.Sprintf("%t", config.Bool("tracetest.analytics")),
		"backend_type":            config.String("tracetest.backend.type"),
	}

	for k, v := range extra {
		props[k] = v
	}

	analytics.Track("Apply", "installer", props)
}
