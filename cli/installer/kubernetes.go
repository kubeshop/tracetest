package installer

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	cliUI "github.com/kubeshop/tracetest/cli/ui"
)

var kubernetes = installer{
	name: "kubernetes",
	preChecks: []preChecker{
		windowsGnuToolsChecker,
		kubectlChecker,
		helmChecker,
		localEnvironmentChecker,
	},
	configs: []configurator{
		configureKubernetes,
		configureTracetest,
		configureKubernetesInstalls,
		configureIngress,
		configureDemoApp,
	},
	installFn: kubernetesInstaller,
}

func windowsGnuToolsChecker(ui UI) {
	if !isWindows() {
		return
	}

	if commandExists("sed") {
		ui.Println(ui.Green("✔ sed already installed"))
		return
	}

	ui.Warning("I didn't find sed in your system")
	option := ui.Select("What do you want to do?", []option{
		{"Install sed", installSed},
		{"Fix it manually", exitOption(
			"Check the helm install docs on https://community.chocolatey.org/packages/sed",
		)},
	}, 0)

	option.fn(ui)

	if commandExists("sed") {
		ui.Println(ui.Green("✔ sed was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ sed could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func installSed(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing sed manually. See https://community.chocolatey.org/packages/sed",
		installDocs:   "https://community.chocolatey.org/packages/sed",
		windows:       "choco install sed",
	}).exec(ui)
}

func kubernetesInstaller(config configuration, ui cliUI.UI) {
	trackInstall("kubernetes", config, nil)

	installCertManager(config, ui)
	installJaegerOperator(config, ui)

	execCmdIgnoreErrors(kubectlCmd(config, "create namespace "+config.String("k8s.namespace")))

	installJaeger(config, ui)
	installCollector(config, ui)

	installTracetest(config, ui)

}

func installTracetest(conf configuration, ui cliUI.UI) {
	setupHelmRepo(conf, ui)

	installTracetestChart(conf, ui)
	fixTracetestConfiguration(conf, ui)

	installOtelCollector(conf, ui)

	execCmd(kubectlNamespaceCmd(conf, "delete pods -l app.kubernetes.io/name=tracetest"), ui)

	installDemo(conf, ui)

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
	if !conf.Bool("demo.enable.pokeshop") {
		return
	}

	helm := helmCmd(conf, "")
	script := strings.ReplaceAll(demoScript, "#helm#", helm)
	script = fmt.Sprintf(script, conf.String("tracetest.collector.endpoint"))

	execCmd(script, ui)
}

func installOtelCollector(conf configuration, ui cliUI.UI) {
	if !conf.Bool("tracetest.collector.install") {
		return
	}

	cc := createTmpFile("collector-config", string(getCollectorConfigFileContents(ui, conf)), ui)
	defer os.Remove(cc.Name())

	execCmd(
		kubectlNamespaceCmd(conf,
			"create configmap collector-config --from-file="+cc.Name()+" -o yaml --dry-run=client",
			"| sed 's#"+path.Base(cc.Name())+"#collector.yaml#' |",
			kubectlNamespaceCmd(conf, "replace -f -"),
		),
		ui,
	)
	execCmd(kubectlNamespaceCmd(conf, "delete pods -l app.kubernetes.io/name=otel-collector"), ui)
}

func fixTracetestConfiguration(conf configuration, ui cliUI.UI) {
	psql := "host=tracetest-postgresql user=tracetest password=not-secure-database-password  port=5432 sslmode=disable"
	c := getTracetestConfigFileContents(psql, ui, conf)
	ttc := createTmpFile("tracetest-config", string(c), ui)
	defer os.Remove(ttc.Name())
	execCmd(
		kubectlNamespaceCmd(conf,
			"create configmap tracetest --from-file="+ttc.Name()+" -o yaml --dry-run=client",
			"| sed 's#"+path.Base(ttc.Name())+"#config.yaml#' |",
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

	if conf.Bool("k8s.expose") {
		cmd = append(cmd, []string{
			"--set ingress.enabled=true",
			"--set 'ingress.hosts[0].host=" + conf.String("k8s.ingress-host") +
				",ingress.hosts[0].paths[0].path=/,ingress.hosts[0].paths[0].pathType=Prefix'",
		}...)
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
	certManagerYaml          = "https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml"
	certManagerClusterIssuer = `
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned
spec:
  selfSigned: {}
`
	jaegerOperatorYaml = "https://github.com/jaegertracing/jaeger-operator/releases/download/v1.32.0/jaeger-operator.yaml"
	jaegerManifest     = `
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger
`
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

func installCollector(config configuration, ui cliUI.UI) {
	if !config.Bool("tracetest.collector.install") {
		return
	}

	execCmd(
		kubectlNamespaceCmd(config, "create -f "+collectorYaml),
		ui,
	)

	ui.Println(ui.Green("✔ collector ready"))

}

func installCertManager(config configuration, ui cliUI.UI) {
	if !config.Bool("k8s.cert-manager.install") {
		return
	}

	execCmd(
		kubectlCmd(config, "create -f "+certManagerYaml),
		ui,
	)

	execCmd(
		kubectlCmd(config, "--namespace cert-manager wait --for=condition=ready pod -l app=webhook --timeout 5m"),
		ui,
	)
	// give it a sec just in case
	time.Sleep(time.Second)

	f := createTmpFile("cert-manager", certManagerClusterIssuer, ui)
	defer os.Remove(f.Name())

	execCmd(
		kubectlCmd(config, "apply -f "+f.Name()),
		ui,
	)

	ui.Println(ui.Green("✔ cert-manager ready"))
}

func installJaegerOperator(config configuration, ui cliUI.UI) {
	if !config.Bool("k8s.jaeger-operator.install") {
		return
	}

	execCmdIgnoreError(
		kubectlCmd(config, "create namespace observability "),
		ui,
	)
	execCmd(
		kubectlCmd(config, "--namespace observability create -f "+jaegerOperatorYaml),
		ui,
	)

	execCmd(
		kubectlCmd(config, "--namespace observability wait --for=condition=ready pod -l name=jaeger-operator --timeout 5m"),
		ui,
	)

	ui.Println(ui.Green("✔ jaeger-operator ready"))

}

func installJaeger(config configuration, ui cliUI.UI) {
	if !config.Bool("tracetest.backend.install") {
		return
	}

	jmf := createTmpFile("tracetest-jaeger", jaegerManifest, ui)
	defer os.Remove(jmf.Name())

	execCmd(
		kubectlNamespaceCmd(config, "apply -f "+jmf.Name()),
		ui,
	)
	ui.Println(ui.Green("✔ jaeger instance ready"))

}

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

func crdExists(config configuration, crd string) bool {
	cmd := fmt.Sprintf(
		"%s > /dev/null 2>&1; echo $?",
		kubectlCmd(config, "get crd "+crd),
	)
	out := getCmdOutputClean(cmd)

	// if output is 0, it means the CRD exists
	return out == "0"
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
	conf.set("k8s.kubeconfig", ui.TextInput("Kubeconfig file", "${HOME}/.kube/config"))

	contexts := getK8sContexts(conf, ui)
	if len(contexts) == 0 {
		ui.Exit(
			"We didn't detect any kubectl contexts available. " +
				"Make sure your kubectl tool is correctly configured and try again. \n" +
				createIssueMsg,
		)
	}
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

	conf.set("k8s.namespace", ui.TextInput("Namespace", "tracetest"))

	return conf
}

func configureIngress(conf configuration, ui cliUI.UI) configuration {
	expose := ui.Confirm("Do you want to expose tracetest (via ingress)?", false)
	if expose {
		conf.set("k8s.ingress-host", ui.TextInput("Host", ""))
	}
	conf.set("k8s.expose", expose)
	return conf
}

func configureKubernetesInstalls(conf configuration, ui cliUI.UI) configuration {
	installCertManager := false
	installJaegerOperator := false
	if conf.Bool("tracetest.backend.install") {
		ui.Warning("I am going to install Jaeger in your cluster!")
		ui.Println("This requires the Jaeger Operator and cert-manager")

		if crdExists(conf, "certificates.cert-manager.io") {
			ui.Println(ui.Green("✔ cert-manager ready"))
		} else {
			ui.Println(ui.Red("✘ cert-manager not available"))
			installCertManager = ui.Confirm("Do you want me to install it?", false)
			if !installCertManager {
				ui.Exit("cert-manager is requried to run the Jaeger Operator. Check the docs at https://cert-manager.io/. " + createIssueMsg)
			}
		}

		if crdExists(conf, "jaegers.jaegertracing.io") {
			ui.Println(ui.Green("✔ jaeger-operator ready"))
		} else {
			ui.Println(ui.Red("✘ jaeger-operator not available"))
			installJaegerOperator = ui.Confirm("Do you want me to install it?", false)
			if !installJaegerOperator {
				ui.Exit("jaeger-operator is requried. Check the docs at https://www.jaegertracing.io/docs/latest/operator/. " + createIssueMsg)
			}
		}
	}
	conf.set("k8s.cert-manager.install", installCertManager)
	conf.set("k8s.jaeger-operator.install", installJaegerOperator)

	return conf
}

func helmChecker(ui cliUI.UI) {
	if commandExists("helm") {
		ui.Println(ui.Green("✔ helm already installed"))
		return
	}

	ui.Warning("I didn't find helm in your system")
	option := ui.Select("What do you want to do?", []cliUI.Option{
		{"Install Helm", installHelm},
		{"Fix it manually", exitOption(
			"Check the helm install docs on https://helm.sh/docs/intro/install/",
		)},
	}, 0)

	option.Fn(ui)

	if commandExists("helm") {
		ui.Println(ui.Green("✔ helm was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ helm could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func kubectlChecker(ui cliUI.UI) {
	if commandExists("kubectl") {
		ui.Println(ui.Green("✔ kubectl already installed"))
		return
	}

	ui.Warning("I didn't find kubectl in your system")
	option := ui.Select("What do you want to do?", []cliUI.Option{
		{"Install Kubectl", installKubectl},
		{"Fix it manually", exitOption(
			"Check the kubectl install docs on https://kubernetes.io/docs/tasks/tools/#kubectl",
		)},
	}, 0)

	option.Fn(ui)

	if commandExists("kubectl") {
		ui.Println(ui.Green("✔ kubectl was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ kubectl could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func localEnvironmentChecker(ui cliUI.UI) {

	option := ui.Select("Are you going to run it locally or in a remote cluster?", []cliUI.Option{
		{"Locally", confirmLocalK8sRunning},
		{"Remote cluster", func(ui cliUI.UI) {}},
	}, 0)

	option.Fn(ui)
}

func confirmLocalK8sRunning(ui cliUI.UI) {
	localK8sRunning := ui.Confirm("Do you have a local kubernentes running?", true)
	if !localK8sRunning {
		option := ui.Select("We can fix that:", []cliUI.Option{
			{"Install minikube", minikubeChecker},
			{"Fix manually", exitOption(
				"Check the minikube install docs on https://minikube.sigs.k8s.io/docs/start/",
			)},
		}, 0)

		option.Fn(ui)
	}
}

func minikubeChecker(ui cliUI.UI) {
	// before procceding, we must start minikube
	// starting in insecure mode to prevent certificate problems
	// TODO: configure the certificate instead to prevent the x509 error
	defer execCmd(`minikube start --insecure-registry="registry-1.docker.io"`, ui)

	if commandExists("minikube") {
		ui.Println(ui.Green("✔ minikube already installed"))
		return
	}

	installMinikube(ui)

	if commandExists("minikube") {
		ui.Println(ui.Green("✔ minikube was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ minikube could not be installed. Check output for errors. " + createIssueMsg))
	}

	if !commandExists("docker") {
		ui.Warning("Docker is required to run your minikube instance")
		dockerChecker(ui)
	}
}

func installHelm(ui cliUI.UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing helm manually. See https://helm.sh/docs/intro/install/",
		installDocs:   "https://helm.sh/docs/intro/install/",
		args:          map[string]string{},
		apt: `
			sudo apt update -y
			sudo apt install -y curl gpg

			curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
			sudo apt-get install apt-transport-https --yes
			echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
			sudo apt-get update
			sudo apt-get install helm
		`,
		yum:      "sudo yum -y install helm",
		dnf:      "sudo dnf -y install helm",
		homebrew: "brew install helm",
		macIntelChipManual: `
			curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
			chmod 700 get_helm.sh
			./get_helm.sh
		`,
		macAppleChipManual: `
			curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
			chmod 700 get_helm.sh
			./get_helm.sh
		`,
		windows: "choco install kubernetes-helm",
		other:   "",
	}).exec(ui)
}

func installKubectl(ui cliUI.UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing kubectl manually. See https://kubernetes.io/docs/tasks/tools/#kubectl",
		installDocs:   "https://kubernetes.io/docs/tasks/tools/#kubectl",
		args:          map[string]string{},
		apt: `
			sudo apt-get update
			sudo apt-get install -y ca-certificates curl apt-transport-https

			sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
			echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

			sudo apt-get update
			sudo apt-get install -y kubectl
		`,
		yum: `
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
sudo yum -y install kubectl
		`,
		dnf: `
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
sudo dnf -y install kubectl
		`,
		homebrew: "brew install helm",
		macIntelChipManual: `
			curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
			chmod +x ./kubectl
			sudo mv ./kubectl /usr/local/bin/kubectl
			sudo chown root: /usr/local/bin/kubectl
		`,
		macAppleChipManual: `
			curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/arm64/kubectl"
			chmod +x ./kubectl
			sudo mv ./kubectl /usr/local/bin/kubectl
			sudo chown root: /usr/local/bin/kubectl
		`,
		windows: "choco install kubernetes-cli",
		other:   "",
	}).exec(ui)
}

func installMinikube(ui cliUI.UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries, you can try installing minikube manually. See https://minikube.sigs.k8s.io/docs/start/",
		installDocs:   "https://minikube.sigs.k8s.io/docs/start/",
		args: map[string]string{
			"Architecture":    detectArchitecture(),
			"RpmArchitecture": strings.Replace(detectArchitecture(), "amd64", "x86_64", 1),
		},
		apt: `
			sudo apt update
			sudo apt install curl -y
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube_latest_{{.Architecture}}.deb
			sudo dpkg -i minikube_latest_{{.Architecture}}.deb
			rm minikube_latest_{{.Architecture}}.deb
		`,
		yum: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-latest.{{.RpmArchitecture}}.rpm
			sudo rpm -Uvh minikube-latest.{{.RpmArchitecture}}.rpm
			rm minikube-latest.{{.RpmArchitecture}}.rpm
		`,
		dnf: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-latest.{{.RpmArchitecture}}.rpm
			sudo rpm -Uvh minikube-latest.{{.RpmArchitecture}}.rpm
			rm minikube-latest.{{.RpmArchitecture}}.rpm
		`,
		homebrew: "brew install minikube",
		macIntelChipManual: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
			sudo install minikube-darwin-amd64 /usr/local/bin/minikube
		`,
		macAppleChipManual: `
			curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-arm64
			sudo install minikube-darwin-arm64 /usr/local/bin/minikube
		`,
		windows: "choco install minikube",
		other:   "",
	}).exec(ui)
}
