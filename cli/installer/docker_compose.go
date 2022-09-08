package installer

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/kubeshop/tracetest/cli/analytics"
	cliConfig "github.com/kubeshop/tracetest/cli/config"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var dockerCompose = installer{
	preChecks: []preChecker{
		dockerChecker,
		dockerReadyChecker,
		dockerComposeChecker,
	},
	configs: []configurator{
		configureTracetest,
		configureDemoApp,
		configureDockerCompose,
	},
	installFn: dockerComposeInstaller,
}

func configureDockerCompose(conf configuration, ui UI) configuration {
	conf.set(
		"output.dir",
		ui.TextInput("Tracetest output directory", "tracetest/"),
	)

	return conf
}

const (
	dockerComposeFilename       = "docker-compose.yaml"
	tracetestConfigFilename     = "tracetest.yaml"
	otelCollectorConfigFilename = "otel-collector.yaml"
)

func dockerComposeInstaller(config configuration, ui UI) {
	analytics.Track("Apply", "installer", map[string]string{
		"type":                    "docker-compose",
		"install_backend":         fmt.Sprintf("%t", config.Bool("tracetest.backend.install")),
		"install_collector":       fmt.Sprintf("%t", config.Bool("tracetest.collector.install")),
		"install_demo":            fmt.Sprintf("%t", config.Bool("demo.enable")),
		"enable_server_analytics": fmt.Sprintf("%t", config.Bool("tracetest.analytics")),
		"backend_type":            config.String("tracetest.backend.type"),
	})

	dir := config.String("output.dir")

	if Force {
		err := os.RemoveAll(dir)
		if err != nil {
			ui.Exit(err.Error())
		}
	}

	collectorConfigFile := getCollectorConfigFileContents(ui, config)
	tracetestConfigFile := getTracetestConfigFileContents(ui, config)
	dockerComposeFile := getDockerComposeFileContents(ui, config)
	dockerComposeFName := filepath.Join(dir, dockerComposeFilename)

	dockerCmd := fmt.Sprintf("docker compose -f %s up -d", dockerComposeFName)
	if f := getCWDDockerComposeIfExists(); f != "" {
		dockerCmd = fmt.Sprintf("docker compose -f %s -f %s up -d", f, dockerComposeFName)
	}

	createDir(ui, dir)
	saveFile(ui, dockerComposeFName, dockerComposeFile)
	saveFile(ui, filepath.Join(dir, tracetestConfigFilename), tracetestConfigFile)
	saveFile(ui, filepath.Join(dir, otelCollectorConfigFilename), collectorConfigFile)

	ui.Success("Install successful!")
	ui.Println(fmt.Sprintf(`
To start tracetest:

	%s

Then, use your browser to navigate to:

  http://localhost:8080

Happy TraceTesting =)
`, dockerCmd))

}

func getCWDDockerComposeIfExists() string {
	opts := []string{
		"./docker-compose.yaml",
		"./docker-compose.yml",
	}
	for _, f := range opts {
		if fileExists(f) {
			return f
		}
	}

	return ""
}

func getDockerComposeFileContents(ui UI, config configuration) []byte {
	project := getCompleteProject(ui)
	include := []string{"tracetest", "postgres"}

	if config.Bool("tracetest.backend.install") {
		include = append(include, "jaeger")
	}

	if config.Bool("tracetest.collector.install") {
		include = append(include, "otel-collector")
	}

	if config.Bool("demo.enable") {
		include = append(include, "cache", "queue", "demo-api", "demo-worker", "demo-rpc")
	}

	filterContainers(ui, project, include)

	if err := fixTracetestContainer(project, cliConfig.Version); err != nil {
		ui.Exit(fmt.Sprintf("cannot configure tracetest container: %s", err.Error()))
	}

	if config.Bool("tracetest.collector.install") {
		if err := fixOtelCollectorContainer(project); err != nil {
			ui.Exit(fmt.Sprintf("cannot configure otel-collector container: %s", err.Error()))
		}
	}

	output, err := yaml.Marshal(project)
	if err != nil {
		ui.Exit(fmt.Sprintf("cannot encode docker-compose file: %s", err.Error()))
	}

	sout := strings.ReplaceAll(string(output), "$", "$$")

	return []byte(sout)
}

func filterContainers(ui UI, project *types.Project, included []string) {
	containers := make(types.Services, 0, len(included))
	if err := project.ForServices(included); err != nil {
		ui.Exit(err.Error())
	}

	for _, sc := range project.Services {
		if !slices.Contains(included, sc.Name) {
			continue
		}

		depMap := types.DependsOnConfig{}
		for sn, sv := range sc.DependsOn {
			if !slices.Contains(included, sn) {
				continue
			}

			depMap[sn] = sv
		}
		sc.DependsOn = depMap

		containers = append(containers, sc)
	}

	project.Services = containers
}

func getCollectorConfigFileContents(ui UI, config configuration) []byte {
	contents, err := getFileContentsForVersion("local-config/collector.config.yaml", cliConfig.Version)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot get collector config file: %w", err).Error())
	}

	type msa = map[string]any
	otelConfig := msa{}

	err = yaml.Unmarshal(contents, &otelConfig)
	if err != nil {
		ui.Exit(err.Error())
	}

	exporters := msa{}
	exporter := ""

	switch config.String("tracetest.backend.type") {
	case "jaeger":
		exporter = "jaeger"
		exporters["jaeger"] = msa{
			"endpoint": config.String("tracetest.backend.endpoint"),
			"tls": msa{
				"insecure": config.Bool("tracetest.backend.tls.insecure"),
			},
		}
	case "tempo":
		exporter = "otlp/2"
		exporters["otlp/2"] = msa{
			"endpoint": config.String("tracetest.backend.endpoint"),
			"tls": msa{
				"insecure": config.Bool("tracetest.backend.tls.insecure"),
			},
		}
	case "opensearch":
		exporter = "otlp/2"
		exporters["otlp/2"] = msa{
			"endpoint": config.String("tracetest.backend.data-prepper.endpoint"),
			"tls": msa{
				"insecure":             config.Bool("tracetest.backend.data-prepper.insecure"),
				"insecure_skip_verify": true,
			},
		}
	case "signalfx":
		exporter = "sapm"
		exporters["sapm"] = msa{
			"access_token":             config.String("tracetest.backend.token"),
			"access_token_passthrough": true,
			"endpoint":                 fmt.Sprintf("https://ingest.%s.signalfx.com/v2/trace", config.String("tracetest.backend.realm")),
			"max_connections":          100,
			"num_workers":              8,
		}
	}

	otelConfig["exporters"] = exporters
	otelConfig["service"].(msa)["pipelines"].(msa)["traces"].(msa)["exporters"] = []string{exporter}

	updated, err := yaml.Marshal(otelConfig)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot get collector config file: %w", err).Error())
	}

	return updated
}

func createDir(ui UI, name string) {
	err := os.Mkdir(name, 0755)
	if err != nil {
		ui.Exit(err.Error())
	}
}

func saveFile(ui UI, fname string, contents []byte) {
	if fileExists(fname) {
		ui.Warning(fmt.Sprintf(`file "%s" already exists.`, fname))
		if !ui.Confirm("Do you want to overwrite it?", false) {
			ui.Exit(fmt.Sprintf(
				`You choose NOT to overwrite "%s". Installation did not succeed`,
				fname,
			))
		}
	}

	err := os.WriteFile(fname, contents, 0644)
	if err != nil {
		ui.Exit(err.Error())
	}
}

func fixTracetestContainer(project *types.Project, version string) error {
	const serviceName = "tracetest"
	tts, err := project.GetService(serviceName)
	if err != nil {
		return err
	}

	if version == "dev" {
		version = "latest"
	}

	tts.Image = "kubeshop/tracetest:" + version
	tts.Build = nil
	tts.Volumes[0].Source = tracetestConfigFilename

	replaceService(project, serviceName, tts)

	return nil
}

func fixOtelCollectorContainer(project *types.Project) error {
	ocs, err := project.GetService("otel-collector")
	if err != nil {
		return err
	}

	ocs.Volumes[0].Source = otelCollectorConfigFilename

	return nil
}

func replaceService(project *types.Project, service string, sc types.ServiceConfig) {
	for i, s := range project.Services {
		if s.Name != service {
			continue
		}

		project.Services[i] = sc
		break
	}
}

func getFileContentsForVersion(path, version string) ([]byte, error) {
	if version == "dev" {
		version = "main"
	}
	url := fmt.Sprintf("https://raw.githubusercontent.com/kubeshop/tracetest/%s/%s", version, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot download file: %w", err)
	}

	contents, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("cannot download file: %w", err)
	}

	return contents, nil

}

func getCompleteProject(ui UI) *types.Project {
	contents, err := getFileContentsForVersion("docker-compose.yaml", cliConfig.Version)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot get docker-compose file: %w", err).Error())
	}

	workingDir, err := os.Getwd()
	if err != nil {
		ui.Panic(err)
	}

	project, err := loader.Load(types.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []types.ConfigFile{
			{Filename: "docker-compose.yaml", Content: contents},
		},
	})
	if err != nil {
		ui.Exit(fmt.Errorf("cannot parse docker-compose file: %w", err).Error())
	}

	return project
}

func dockerChecker(ui UI) {
	if commandExists("docker") {
		ui.Println(ui.Green("✔ docker already installed"))
		return
	}

	ui.Warning("I didn't find docker in your system")
	option := ui.Select("What do you want to do?", []option{
		{"Install Docker Engine", installDockerEngine},
		{"Install Docker Desktop", installDockerDesktop},
		{"Fix manually", exitOption(
			"Check the docker install docs on https://docs.docker.com/get-docker/",
		)},
	}, 0)

	option.fn(ui)

	if commandExists("docker") {
		ui.Println(ui.Green("✔ docker was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ docker could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func dockerReadyChecker(ui UI) {
	if commandSuccess("docker ps") {
		ui.Println(ui.Green("✔ docker is ready"))
		return
	}

	ui.Exit(`
		Docker doesn't seem to be responding correctly. Please make sure the service is correctly installed and started.
		You can try to connect to docker using:

			docker ps

		` + createIssueMsg,
	)
}

func dockerComposeChecker(ui UI) {
	if commandSuccess("docker compose") {
		ui.Println(ui.Green("✔ docker compose already installed"))
		return
	}

	ui.Warning("I didn't find docker compose in your system")
	option := ui.Select("What do you want to do?", []option{
		{"Install Docker Compose", installDockerCompose},
		{"Fix manually", exitOption(
			"Check the docker compose install docs on https://docs.docker.com/compose/install/",
		)},
	}, 0)

	option.fn(ui)

	if commandSuccess("docker compose") {
		ui.Println(ui.Green("✔ docker compose was successfully installed"))
	} else {
		ui.Exit(ui.Red("✘ docker compose could not be installed. Check output for errors. " + createIssueMsg))
	}
}

func installDockerCompose(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries. You can try installing Docker Compose manually. See https://docs.docker.com/compose/install/",
		args: map[string]string{
			"DockerVersion": dockerVersion(ui),
			"Architecture":  detectArchitecture(),
		},
		apt: `
			DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')

			# Repo install. see https://docs.docker.com/engine/install/$DISTRO/#install-using-the-repository
			# setup repo
			sudo apt-get -y update
			sudo apt-get -y install \
				ca-certificates \
				curl \
				gnupg \
				lsb-release
			sudo mkdir -p /etc/apt/keyrings
			curl -fsSL https://download.docker.com/linux/$DISTRO/gpg | sudo gpg --batch --yes --dearmor -o /etc/apt/keyrings/docker.gpg

			sudo chmod a+r /etc/apt/keyrings/docker.gpg

			echo \
				"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$DISTRO \
				$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

			# install
			sudo apt-get -y update
			sudo apt-get -y install docker-compose-plugin
			`,
		dnf: `
			# Repo install. see https://docs.docker.com/engine/install/fedora/#install-using-the-repository
			# setup repo
			sudo dnf -y install dnf-plugins-core
			sudo dnf -y config-manager \
				--add-repo \
				https://download.docker.com/linux/fedora/docker-ce.repo

			# install
			sudo dnf -y install docker-compose-plugin
			`,
		yum: `
			# Repo install. see https://docs.docker.com/engine/install/centos/#install-using-the-repository
			# setup repo
			sudo yum install -y yum-utils
			sudo yum-config-manager \
				--add-repo \
				https://download.docker.com/linux/centos/docker-ce.repo

			# install
			sudo yum -y install docker-compose-plugin
			`,
		homebrew:           "brew install docker-compose",
		macIntelChipManual: "TODO", // todo. see https://apple.stackexchange.com/a/73931
		macAppleChipManual: "TODO", // todo. see https://apple.stackexchange.com/a/73931
		windows:            "Check the install docs: https://docs.docker.com/compose/install/",
		other:              "Check the install docs: https://docs.docker.com/compose/install/",
	}).exec(ui)
}

func installDockerEngine(ui UI) {
	post := `
			# post-install. not neccesary for root
			if [ "$(id -u)" != "0" ]; then
				sudo groupadd docker
				sudo usermod -aG docker $USER
				newgrp docker
			fi
			`
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries. You can try installing Docker Engine manually. See https://docs.docker.com/engine/install/",
		installDocs:   "https://docs.docker.com/engine/install/",
		args: map[string]string{
			"DockerVersion": dockerVersion(ui),
			"Architecture":  detectArchitecture(),
		},
		apt: `
			# Repo install. see https://docs.docker.com/engine/install/$DISTRO/#install-using-the-repository
			# setup repo
			sudo apt-get -y update
			sudo apt-get -y install \
				ca-certificates \
				curl \
				gnupg \
				lsb-release

			DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')

			# cleanup. see https://docs.docker.com/engine/install/$DISTRO/#uninstall-old-versions
			sudo apt-get -y remove docker docker-engine docker.io containerd runc

			sudo mkdir -p /etc/apt/keyrings
			curl -fsSL https://download.docker.com/linux/$DISTRO/gpg | sudo gpg --batch --yes --dearmor -o /etc/apt/keyrings/docker.gpg

			sudo chmod a+r /etc/apt/keyrings/docker.gpg

			echo \
				"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$DISTRO \
				$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

			# install
			sudo apt-get -y update
			sudo apt-get -y install docker-ce docker-ce-cli containerd.io docker-compose-plugin
			` + post,
		dnf: `
			# cleanup. see https://docs.docker.com/engine/install/fedora/#uninstall-old-versions
			sudo dnf -y remove docker \
					docker-client \
					docker-client-latest \
					docker-common \
					docker-latest \
					docker-latest-logrotate \
					docker-logrotate \
					docker-selinux \
					docker-engine-selinux \
					docker-engine


			# Repo install. see https://docs.docker.com/engine/install/fedora/#install-using-the-repository
			# setup repo
			sudo dnf -y install dnf-plugins-core
			sudo dnf -y config-manager \
				--add-repo \
				https://download.docker.com/linux/fedora/docker-ce.repo

			# install
			sudo dnf -y install docker-ce docker-ce-cli containerd.io docker-compose-plugin
			sudo systemctl start docker
			` + post,
		yum: `
			# cleanup. see https://docs.docker.com/engine/install/centos/#uninstall-old-versions
			sudo yum -y remove docker \
					docker-client \
					docker-client-latest \
					docker-common \
					docker-latest \
					docker-latest-logrotate \
					docker-logrotate \
					docker-engine \
					podman \
					runc

			# Repo install. see https://docs.docker.com/engine/install/centos/#install-using-the-repository
			# setup repo
			sudo yum install -y yum-utils
			sudo yum-config-manager \
				--add-repo \
				https://download.docker.com/linux/centos/docker-ce.repo

			# install
			sudo yum -y install docker-ce docker-ce-cli containerd.io docker-compose-plugin
			sudo systemctl start docker
			` + post,
		homebrew:           "brew install docker",
		macIntelChipManual: "", // empty means not supported
		macAppleChipManual: "", // empty means not supported
		windows:            "", // empty means not supported
		other:              "", // empty means not supported
	}).exec(ui)
}

func installDockerDesktop(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries. You can try installing Docker Desktop manually. See https://docs.docker.com/desktop/install/",
		installDocs:   "https://docs.docker.com/desktop/",
		args: map[string]string{
			"DockerVersion": dockerVersion(ui),
			"Architecture":  detectArchitecture(),
		},
		apt: `
			# add docker repo. See https://docs.docker.com/engine/install/$DISTRO/#set-up-the-repository
			sudo apt-get -y update
			sudo apt-get -y install \
				ca-certificates \
				curl \
				gnupg \
				lsb-release

			DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')

			# cleanup old installs. See https://docs.docker.com/desktop/install/$DISTRO/#prerequisites
			rm -rf $HOME/.docker/desktop
			sudo rm -f /usr/local/bin/com.docker.cli
			sudo apt-get -y purge docker-desktop

			sudo mkdir -p /etc/apt/keyrings
			curl -fsSL https://download.docker.com/linux/$DISTRO/gpg | sudo gpg --yes --dearmor -o /etc/apt/keyrings/docker.gpg
			echo \
				"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/$DISTRO \
				$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
			sudo apt-get -y update

			# download and install docker desktop. See https://docs.docker.com/desktop/install/$DISTRO/#install-docker-desktop
			curl -LO https://desktop.docker.com/linux/main/{{.Architecture}}/docker-desktop-{{.DockerVersion}}-{{.Architecture}}.deb
			sudo apt-get -y install ./docker-desktop-{{.DockerVersion}}-{{.Architecture}}.deb
			rm ./docker-desktop-{{.DockerVersion}}-{{.Architecture}}.deb

			# enable docker desktop auto start
			systemctl --user enable docker-desktop
			systemctl start docker-desktop
			`,
		dnf: `
			# setup repo. See https://docs.docker.com/engine/install/fedora/#set-up-the-repository
			sudo dnf -y install dnf-plugins-core
			sudo dnf -y config-manager \
				--add-repo \
				https://download.docker.com/linux/fedora/docker-ce.repo

			# download and install docker desktop. See https://docs.docker.com/desktop/install/$DISTRO/#install-docker-desktop
			curl -LO https://desktop.docker.com/linux/main/{{.Architecture}}/docker-desktop-{{.DockerVersion}}-{{.Architecture}}.deb
			sudo dnf -y install ./docker-desktop-{{.DockerVersion}}-{{.Architecture}}.deb

			# enable docker desktop auto start
			systemctl --user enable docker-desktop
			`,
		homebrew: "brew install --cask docker",
		macIntelChipManual: `
			curl -LO https://desktop.docker.com/mac/main/amd64/Docker.dmg
			sudo hdiutil attach $(pwd)/Docker.dmg
			sudo cp -R /Volumes/Docker/Docker.app /Applications
			sudo hdiutil detach /Volumes/Docker
			rm -f Docker.dmg
			`,
		macAppleChipManual: `
			curl -LO https://desktop.docker.com/mac/main/arm64/Docker.dmg
			sudo hdiutil attach $(pwd)/Docker.dmg
			sudo cp -R /Volumes/Docker/Docker.app /Applications
			sudo hdiutil detach /Volumes/Docker
			rm -f Docker.dmg
		`,
		windows: "", // empty means not supported
		other:   "", // empty means not supported
	}).exec(ui)
}

func dockerVersion(ui UI) string {
	resp, err := http.Get("https://docs.docker.com/desktop/release-notes/")

	if err != nil {
		ui.Panic(fmt.Errorf("cannot get docker releases: %w", err))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ui.Panic(fmt.Errorf("cannot get docker releases: %w", err))
	}

	// matches "Docker Desktop 4.10.1", case insensitive
	var dockerVersionRegex = regexp.MustCompile(`(?i)docker desktop (\d+\.\d+\.\d+)`)

	matches := dockerVersionRegex.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		ui.Panic(fmt.Errorf("could not find a valid docker desktop version"))
	}
	return matches[1]
}
