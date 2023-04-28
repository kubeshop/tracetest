package installer

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	cliConfig "github.com/kubeshop/tracetest/cli/config"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

var kubernetes = installer{
	name: "kubernetes",
	preChecks: []preChecker{
		windowsGnuToolsChecker,
		kubectlChecker,
		helmChecker,
	},
	configs: []configurator{
		configureKubernetes,
		configureTracetest,
		configureIngress,
		configureDemoApp,
	},
	installFn: kubernetesInstaller,
}

func windowsGnuToolsChecker(ui cliUI.UI) {
	if !isWindows() {
		return
	}

	if commandExists("sed") {
		ui.Println(ui.Green("✔ sed already installed"))
		return
	}

	ui.Warning("I didn't find sed in your system")
	option := ui.Select("What do you want to do?", []cliUI.Option{
		{Text: "Install sed", Fn: installSed},
		{Text: "Fix it manually", Fn: exitOption(
			"Check the helm install docs on https://community.chocolatey.org/packages/sed",
		)},
	}, 0)

	option.Fn(ui)

	if commandExists("sed") {
		ui.Println(ui.Green("✔ sed was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ sed could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func installSed(ui cliUI.UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing sed manually. See https://community.chocolatey.org/packages/sed",
		installDocs:   "https://community.chocolatey.org/packages/sed",
		windows:       "choco install sed",
	}).exec(ui)
}

func kubernetesInstaller(config configuration, ui cliUI.UI) {
	trackInstall("kubernetes", config, nil)

	execCmdIgnoreErrors(kubectlCmd(config, "create namespace "+config.String("k8s.namespace")))

	if !config.Bool("installer.only_tracetest") {
		installCollector(config, ui)
	}
	installTracetest(config, ui)
}

func installCollector(config configuration, ui cliUI.UI) {
	execCmd(
		kubectlNamespaceCmd(config, "apply -f "+collectorYaml),
		ui,
	)

	ui.Println(ui.Green("✔ collector ready"))
}

func installTracetest(conf configuration, ui cliUI.UI) {
	setupHelmRepo(conf, ui)

	installTracetestChart(conf, ui)
	fixTracetestConfiguration(conf, ui)

	if !conf.Bool("installer.only_tracetest") {
		installOtelCollector(conf, ui)
	}

	execCmd(kubectlNamespaceCmd(conf, "delete pods -l app.kubernetes.io/name=tracetest"), ui)

	if !conf.Bool("installer.only_tracetest") {
		installDemo(conf, ui)
	}

	ui.Success("Install successful!")
	ui.Println(fmt.Sprintf(`
To access tracetest:

	%s

Then, use your browser to navigate to:

  http://localhost:11633

Happy TraceTesting =)
`, kubectlNamespaceCmd(conf, "port-forward svc/tracetest 11633")))

}

func installDemo(conf configuration, ui cliUI.UI) {
	helm := helmCmd(conf, "")
	script := strings.ReplaceAll(demoScript, "#helm#", helm)
	script = fmt.Sprintf(script, conf.String("tracetest.backend.endpoint.collector"))

	execCmd(script, ui)
}

func installOtelCollector(conf configuration, ui cliUI.UI) {
	cc := createTmpFile("collector-config", string(getCollectorConfigFileContents(ui, conf)), ui)
	defer os.Remove(cc.Name())

	cmdString := kubectlNamespaceCmd(conf,
		"create configmap collector-config --from-file="+cc.Name()+" -o yaml --dry-run=client",
		"| sed 's#"+path.Base(cc.Name())+"#collector.yaml#' |",
		kubectlNamespaceCmd(conf, "replace -f -"),
	)

	execCmd(
		cmdString,
		ui,
	)
	execCmd(kubectlNamespaceCmd(conf, "delete pods -l app.kubernetes.io/name=otel-collector"), ui)
}

func fixTracetestConfiguration(conf configuration, ui cliUI.UI) {
	c := getTracetestConfigFileContents("tracetest-postgresql", "tracetest", "not-secure-database-password", ui, conf)
	ttc := createTmpFile("tracetest-config", string(c), ui)
	defer os.Remove(ttc.Name())

	p := getTracetestProvisionFileContents(ui, conf)
	ttp := createTmpFile("tracetest-provisioning", string(p), ui)
	defer os.Remove(ttp.Name())

	execCmd(
		kubectlNamespaceCmd(conf,
			"create configmap tracetest --from-file="+ttc.Name()+" --from-file="+ttp.Name()+" -o yaml --dry-run=client",
			"| sed 's#"+path.Base(ttc.Name())+"#config.yaml#'",
			"| sed 's#"+path.Base(ttp.Name())+"#provisioning.yaml#' |",
			kubectlNamespaceCmd(conf, "replace -f -"),
		),
		ui,
	)
}

func installTracetestChart(conf configuration, ui cliUI.UI) {
	cmd := []string{
		"upgrade --install tracetest kubeshop/tracetest",
		"--namespace " + conf.String("k8s.namespace") + " --create-namespace",
	}

	if cliConfig.Version == "dev" {
		cmd = append(cmd, "--set image.tag=latest")
	}

	if os.Getenv("TRACETEST_DEV") != "" {
		cmd = append(cmd, "--set env.tracetestDev=true")
	}

	execCmd(helmCmd(conf, cmd...), ui)
}

func setupHelmRepo(conf configuration, ui cliUI.UI) {
	execCmd(
		helmCmd(conf, "repo add --force-update kubeshop https://kubeshop.github.io/helm-charts"),
		ui,
	)
	execCmd(
		helmCmd(conf, "repo update"),
		ui,
	)
}

func helmCmd(config configuration, cmd ...string) string {
	return fmt.Sprintf(
		"helm --kubeconfig %s --kube-context %s %s",
		config.String("k8s.kubeconfig"),
		config.String("k8s.context"),
		strings.Join(cmd, " "),
	)
}

const (
	collectorYaml = "https://raw.githubusercontent.com/kubeshop/tracetest/main/k8s/collector.yml"

	demoScript = `
tmpdir=$(mktemp -d)
curl -L https://github.com/kubeshop/pokeshop/tarball/master | tar -xz --strip-components 1 -C  $tmpdir
cd $tmpdir/helm-chart
#helm# dependency update

#helm# upgrade --install demo . \
  --namespace demo --create-namespace \
  -f values.yaml \
  --set image.tag=latest \
  --set image.pullPolicy=Always \
  --set postgres.auth.username=ashketchum,postgres.auth.password=squirtle123,postgres.auth.database=pokeshop \
  --set rabbitmq.auth.username=guest,rabbitmq.auth.password=guest,rabbitmq.auth.erlangCookie=secretcookie \
  --set 'env[4].value=%s'
`
)

func createTmpFile(name, contents string, ui cliUI.UI) *os.File {
	f, err := os.CreateTemp("", name)
	if err != nil {
		ui.Exit(fmt.Sprintf("Cannot create temp %s file: %s", name, err))
	}

	if _, err := f.Write([]byte(contents)); err != nil {
		ui.Exit(fmt.Sprintf("Cannot write temp %s file: %s", name, err))
	}

	if err := f.Close(); err != nil {
		ui.Exit(fmt.Sprintf("Cannot close temp %s file: %s", name, err))
	}

	return f
}

func kubectlNamespaceCmd(config configuration, cmd ...string) string {
	ns := "--namespace " + config.String("k8s.namespace")

	return kubectlCmd(config, append([]string{ns}, cmd...)...)
}

func kubectlCmd(config configuration, cmd ...string) string {
	return fmt.Sprintf(
		"kubectl --kubeconfig %s --context %s %s",
		config.String("k8s.kubeconfig"),
		config.String("k8s.context"),
		strings.Join(cmd, " "),
	)
}

type k8sContext struct {
	name     string
	selected bool
}

func getK8sContexts(conf configuration, ui cliUI.UI) []k8sContext {
	records, err := getKubernetesContextArray(conf.String("k8s.kubeconfig"))
	if err != nil {
		ui.Exit(fmt.Sprintf("cannot get kubectl contexts: %s", err.Error()))
	}

	results := []k8sContext{}
	for _, r := range records {
		results = append(results, k8sContext{
			name:     r[1],
			selected: r[0] == "*",
		})
	}

	return results
}

func getKubernetesContextArray(kubeconfig string) ([][]string, error) {
	output := getCmdOutput(fmt.Sprintf(
		`kubectl --kubeconfig %s config get-contexts --no-headers`,
		kubeconfig,
	))

	// replace spaces with comma
	spaceRegex := regexp.MustCompile(`[ ]+`)
	newStringBytes := spaceRegex.ReplaceAll([]byte(output), []byte(","))
	output = string(newStringBytes)

	records, err := csv.NewReader(strings.NewReader(output)).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func configureKubernetes(conf configuration, ui cliUI.UI) configuration {
	conf.set("k8s.kubeconfig", "${HOME}/.kube/config")

	contexts := getK8sContexts(conf, ui)
	if len(contexts) == 0 {
		ui.Exit(
			"We didn't detect any kubectl contexts available. " +
				"Make sure your kubectl tool is correctly configured and try again. \n" +
				createIssueMsg,
		)
	}

	if len(contexts) == 1 {
		conf.set("k8s.context", contexts[0].name)
	}

	if len(contexts) > 1 {
		options := []cliUI.Option{}
		defaultIndex := 0
		for i, c := range contexts {
			if c.selected {
				defaultIndex = i
			}
			options = append(options, cliUI.Option{Text: c.name, Fn: func(ui cliUI.UI) {}})
		}

		selected := ui.Select("Kubectl context", options, defaultIndex)
		conf.set("k8s.context", selected.Text)
	}

	conf.set("k8s.namespace", "tracetest")
	return conf
}

func configureIngress(conf configuration, ui cliUI.UI) configuration {
	conf.set("k8s.ingress-host", "tracetest")
	return conf
}

func helmChecker(ui cliUI.UI) {
	if commandExists("helm") {
		ui.Println(ui.Green("✔ helm already installed"))
		return
	}

	ui.Exit("I didn't find helm in your system. Check the helm install docs on https://helm.sh/docs/intro/install/")
}

func kubectlChecker(ui cliUI.UI) {
	if commandExists("kubectl") {
		ui.Println(ui.Green("✔ kubectl already installed"))
		return
	}

	ui.Exit("I didn't find kubectl in your system")
}
