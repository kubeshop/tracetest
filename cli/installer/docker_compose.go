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
	cliConfig "github.com/kubeshop/tracetest/cli/config"
	cliUI "github.com/kubeshop/tracetest/cli/ui"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var localSystemDockerComposeCommand = "docker compose"

var dockerCompose = installer{
	name: "docker-compose",
	preChecks: []preChecker{
		chocolateyForWindowsChecker,
		wslChecker,
		dockerChecker,
		dockerReadyChecker,
		dockerComposeChecker,
	},
	configs: []configurator{
		configureDockerCompose,
		configureTracetest,
		configureDemoApp,
		configureDockerComposeOutput,
	},
	installFn: dockerComposeInstaller,
}

func configureDockerCompose(conf configuration, ui cliUI.UI) configuration {
	conf.set("project.docker-compose.filename", dockerComposeFilename)
	conf.set("project.docker-compose.create", true)

	return conf
}

func configureDockerComposeOutput(conf configuration, ui cliUI.UI) configuration {
	conf.set("output.dir", "tracetest/")

	return conf
}

const (
	dockerComposeFilename       = "docker-compose.yaml"
	tracetestConfigFilename     = "tracetest.yaml"
	tracetestProvisionFilename  = "tracetest-provision.yaml"
	otelCollectorConfigFilename = "collector.config.yaml"
)

func dockerComposeInstaller(config configuration, ui cliUI.UI) {
	dir := config.String("output.dir")

	err := os.RemoveAll(dir)
	if err != nil {
		ui.Exit(err.Error())
	}

	tracetestConfigFile := getTracetestConfigFileContents("postgres", "postgres", "postgres", ui, config)

	dockerComposeFile := getDockerComposeFileContents(ui, config)
	dockerComposeFName := filepath.Join(dir, dockerComposeFilename)

	dockerCmd := fmt.Sprintf(
		"%s -f %s up -d",
		localSystemDockerComposeCommand,
		dockerComposeFName,
	)

	createDir(ui, dir)
	saveFile(ui, dockerComposeFName, dockerComposeFile)
	saveFile(ui, filepath.Join(dir, tracetestConfigFilename), tracetestConfigFile)

	tracetestProvisionFile := getTracetestProvisionFileContents(ui, config)
	saveFile(ui, filepath.Join(dir, tracetestProvisionFilename), tracetestProvisionFile)

	if !config.Bool("installer.only_tracetest") {
		collectorConfigFile := getCollectorConfigFileContents(ui, config)
		saveFile(ui, filepath.Join(dir, otelCollectorConfigFilename), collectorConfigFile)
	}

	ui.Success("Install successful!")
	ui.Println(fmt.Sprintf(`
To start tracetest:

	%s

Then, use your browser to navigate to:

  http://localhost:11633

Happy TraceTesting =)
`, dockerCmd))

}

func getDockerComposeFileContents(ui cliUI.UI, config configuration) []byte {
	project := getCompleteProject(ui, config)
	include := []string{"tracetest", "postgres"}

	if config.Bool("demo.enable.pokeshop") {
		include = append(include, "cache", "queue", "stream", "demo-api", "demo-worker", "demo-rpc", "demo-streaming-worker", "otel-collector")
	}

	// filter and update project
	filterAndFixContainers(ui, project, include)

	// set version for tracetest container
	if err := fixTracetestContainer(config, project, cliConfig.Version); err != nil {
		ui.Exit(fmt.Sprintf("cannot configure tracetest container: %s", err.Error()))
	}

	//remove provision parameters if we will not install a tracing backend
	if !config.Bool("tracetest.backend.install") {
		removeProvisioningInfo(ui, project)
	}

	output, err := yaml.Marshal(project)
	if err != nil {
		ui.Exit(fmt.Sprintf("cannot encode docker-compose file: %s", err.Error()))
	}

	sout := fixPortConfig(string(output))
	sout = strings.ReplaceAll(sout, "$", "$$")
	sout = strings.ReplaceAll(sout, "$${TRACETEST_DEV}", "${TRACETEST_DEV}")

	return []byte(sout)
}

func fixPortConfig(output string) string {
	publishedPortRegex := regexp.MustCompile(`published: "\d+"`)
	return publishedPortRegex.ReplaceAllStringFunc(output, func(s string) string {
		return strings.ReplaceAll(s, `"`, ``)
	})
}

func filterAndFixContainers(ui cliUI.UI, project *types.Project, included []string) {
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

func removeProvisioningInfo(ui cliUI.UI, project *types.Project) {
	const serviceName = "tracetest"

	tts, err := project.GetService(serviceName)
	if err != nil {
		ui.Exit(err.Error())
	}

	// set command as empty
	tts.Command = []string{}

	// remove provisioning volume
	tts.Volumes = []types.ServiceVolumeConfig{tts.Volumes[0]}
	replaceService(project, serviceName, tts)
}

type msa = map[string]any

func getCollectorConfigFileContents(ui cliUI.UI, config configuration) []byte {
	contents, err := getFileContentsForVersion("examples/collector/collector.config.yaml", cliConfig.Version)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot get collector config file: %w", err).Error())
	}

	otelConfig := msa{}

	err = yaml.Unmarshal(contents, &otelConfig)
	if err != nil {
		ui.Exit(err.Error())
	}

	exporter := "otlp/1"
	exporters := msa{
		"otlp/1": msa{
			"endpoint": config.String("tracetest.backend.endpoint"),
			"tls": msa{
				"insecure": config.Bool("tracetest.backend.tls.insecure"),
			},
		},
	}

	otelConfig["exporters"] = exporters

	otelConfig["service"].(msa)["pipelines"].(msa)["traces/1"].(msa)["exporters"] = []string{exporter}

	updated, err := yaml.Marshal(otelConfig)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot get collector config file: %w", err).Error())
	}

	return updated
}

func createDir(ui cliUI.UI, name string) {
	err := os.Mkdir(name, 0755)
	if err != nil {
		ui.Exit(err.Error())
	}
}

func saveFile(ui cliUI.UI, fname string, contents []byte) {
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

func fixTracetestContainer(config configuration, project *types.Project, version string) error {
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
	tracetestDevEnv := "${TRACETEST_DEV}"
	tts.Environment["TRACETEST_DEV"] = &tracetestDevEnv

	replaceService(project, serviceName, tts)

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

func getCompleteProject(ui cliUI.UI, config configuration) *types.Project {
	tracetestDCContents, err := getFileContentsForVersion("examples/collector/docker-compose.yml", cliConfig.Version)
	if err != nil {
		ui.Exit(fmt.Errorf("cannot get docker-compose file: %w", err).Error())
	}

	configFiles := []types.ConfigFile{
		{Filename: "docker-compose.yaml", Content: tracetestDCContents},
	}

	if config.Bool("demo.enable.pokeshop") {
		demoDCContents, err := getFileContentsForVersion("examples/docker-compose.demo.yaml", cliConfig.Version)
		if err != nil {
			ui.Exit(fmt.Errorf("cannot get docker-compose file: %w", err).Error())
		}
		configFiles = append(configFiles, types.ConfigFile{Filename: "docker-compose.yaml", Content: demoDCContents})
	}

	workingDir, err := os.Getwd()
	if err != nil {
		ui.Panic(err)
	}

	project, err := loader.Load(types.ConfigDetails{
		WorkingDir:  workingDir,
		ConfigFiles: configFiles,
		Environment: map[string]string{
			"TRACETEST_DEV": "",
		},
	})
	if err != nil {
		ui.Exit(fmt.Errorf("cannot parse docker-compose file: %w", err).Error())
	}

	return project
}

func dockerChecker(ui cliUI.UI) {
	if commandExists("docker") {
		ui.Println(ui.Green("✔ docker already installed"))
		return
	}

	ui.Warning("I didn't find docker in your system")
	exitOption(
		"Check the docker install docs on https://docs.docker.com/get-docker/",
	)(ui)
}

func dockerReadyChecker(ui cliUI.UI) {
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

func dockerComposeChecker(ui cliUI.UI) {
	if commandSuccess("docker-compose") {
		localSystemDockerComposeCommand = "docker-compose"
		ui.Println(ui.Green("✔ docker-compose already installed"))
		return
	}

	if commandSuccess("docker compose") {
		localSystemDockerComposeCommand = "docker compose"
		ui.Println(ui.Green("✔ docker compose already installed"))
		return
	}

	ui.Warning("I didn't find docker compose in your system")
	exitOption(
		"Check the docker compose install docs on https://docs.docker.com/compose/install/",
	)(ui)
}
