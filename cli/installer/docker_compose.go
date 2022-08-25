package installer

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	cliConfig "github.com/kubeshop/tracetest/cli/config"
	"gopkg.in/yaml.v3"
)

var dockerCompose = installer{
	preChecks: []preChecker{
		dockerChecker,
		dockerReadyChecker,
		dockerComposeChecker,
	},
	configs: []configurator{
		configureDemoApp,
		configureDockerCompose,
	},
	installFn: dockerComposeInstaller,
}

func configureDockerCompose(conf configuration, ui UI) configuration {
	conf["docker-compose.filename"] = ui.TextInput("Docker Compose output file name", "docker-compose.yaml")
	conf["tracetest-config.filename"] = ui.TextInput("TraceTest Config output file name", "tracetest.yaml")
	conf["collector-config.filename"] = ui.TextInput("OTel-Collector Config output file name", "collector.yaml")

	return conf
}

func dockerComposeInstaller(config configuration, ui UI) {

	collectorConfigFile := getCollectorConfigFileContents(ui, config)
	collectorConfigFName, err := config.String("collector-config.filename")
	if err != nil {
		ui.Exit(err.Error())
	}

	tracetestConfigFile := getTracetestConfigFileContents(ui, config)
	tracetestConfigFName, err := config.String("tracetest-config.filename")
	if err != nil {
		ui.Exit(err.Error())
	}

	dockerComposeFile := getDockerComposeFileContents(ui, config)
	dockerComposeFName, err := config.String("docker-compose.filename")
	if err != nil {
		ui.Exit(err.Error())
	}

	saveFile(ui, collectorConfigFName, collectorConfigFile)
	saveFile(ui, tracetestConfigFName, tracetestConfigFile)
	saveFile(ui, dockerComposeFName, dockerComposeFile)

	ui.Success("Install successful!")
	ui.Println(fmt.Sprintf(`
To start tracetest:

  docker compose -f %s up -d

Then, use your browser to navigate to:

  http://localhost:8080

Happy TraceTesting =)
`, dockerComposeFName))

}

func getDockerComposeFileContents(ui UI, config configuration) []byte {
	project := getCompleteProject(ui)
	include := []string{"tracetest", "postgres", "otel-collector", "jaeger"}

	if includeDemo, err := config.Bool("demo.enable"); err != nil {
		ui.Exit(err.Error())
	} else if includeDemo {
		include = append(include, "cache", "queue", "demo-api", "demo-worker", "demo-rpc")
	}

	if err := project.ForServices(include); err != nil {
		ui.Exit(err.Error())
	}

	ttfile, err := config.String("tracetest-config.filename")
	if err != nil {
		ui.Exit(err.Error())
	}
	if err := fixTracetestContainer(project, ttfile, cliConfig.Version); err != nil {
		ui.Exit(err.Error())
	}

	ocfile, err := config.String("collector-config.filename")
	if err != nil {
		ui.Exit(err.Error())
	}
	if err := fixOtelCollectorContainer(project, ocfile); err != nil {
		ui.Exit(err.Error())
	}

	output, err := yaml.Marshal(project)
	if err != nil {
		ui.Exit(err.Error())
	}

	return output
}

func getCollectorConfigFileContents(ui UI, config configuration) []byte {
	contents, err := getFileContentsForVersion("local-config/collector.config.yaml", cliConfig.Version)
	if err != nil {
		ui.Exit(fmt.Errorf("Cannot get docker-compose file: %w", err).Error())
	}

	return contents
}

func getTracetestConfigFileContents(ui UI, config configuration) []byte {
	contents, err := getFileContentsForVersion("local-config/config.tests.yaml", cliConfig.Version)
	if err != nil {
		ui.Exit(fmt.Errorf("Cannot get docker-compose file: %w", err).Error())
	}

	return contents
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

func fixTracetestContainer(project *types.Project, tracetestConfigFile, version string) error {
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
	tts.Volumes[0].Source = tracetestConfigFile

	replaceService(project, serviceName, tts)

	return nil
}

func fixOtelCollectorContainer(project *types.Project, collectorConfigFile string) error {
	ocs, err := project.GetService("otel-collector")
	if err != nil {
		return err
	}

	ocs.Volumes[0].Source = collectorConfigFile

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
		ui.Exit(fmt.Errorf("Cannot get docker-compose file: %w", err).Error())
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
		ui.Exit(fmt.Errorf("Cannot parse docker-compose file: %w", err).Error())
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
			"Check the docker install docks on https://docs.docker.com/get-docker/",
		)},
	})

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
			"Check the docker compose install docks on https://docs.docker.com/compose/install/",
		)},
	})

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
		windows:            "Check the install docks: https://docs.docker.com/compose/install/",
		other:              "Check the install docks: https://docs.docker.com/compose/install/",
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
		args: map[string]string{
			"DockerVersion": dockerVersion(ui),
			"Architecture":  detectArchitecture(),
		},
		apt: `
			DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')

			# cleanup. see https://docs.docker.com/engine/install/$DISTRO/#uninstall-old-versions
			sudo apt-get -y remove docker docker-engine docker.io containerd runc

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
		macIntelChipManual: "TODO", // todo. see https://apple.stackexchange.com/a/73931
		macAppleChipManual: "TODO", // todo. see https://apple.stackexchange.com/a/73931
		windows:            "Check the install docks: https://docs.docker.com/engine/install/",
		other:              "Check the install docks: https://docs.docker.com/engine/install/",
	}).exec(ui)
}

func installDockerDesktop(ui UI) {
	(cmd{
		sudo:          true,
		notConfirmMsg: "No worries. You can try installing Docker Desktop manually. See https://docs.docker.com/desktop/install/",
		args: map[string]string{
			"DockerVersion": dockerVersion(ui),
			"Architecture":  detectArchitecture(),
		},
		apt: `
			DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')
			# cleanup old installs. See https://docs.docker.com/desktop/install/$DISTRO/#prerequisites
			rm -rf $HOME/.docker/desktop
			sudo rm -f /usr/local/bin/com.docker.cli
			sudo apt-get -y purge docker-desktop

			# add docker repo. See https://docs.docker.com/engine/install/$DISTRO/#set-up-the-repository
			sudo apt-get -y update
			sudo apt-get -y install \
				ca-certificates \
				curl \
				gnupg \
				lsb-release
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
		windows:  "Check the install docks: https://docs.docker.com/desktop/install/windows-install/",
		other:    "Check the install docks: https://docs.docker.com/desktop/#download-and-install",
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
