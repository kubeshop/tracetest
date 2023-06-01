package environment

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
	"github.com/kubeshop/tracetest/cli-e2etest/config"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var (
	mutex               = sync.Mutex{}
	envCounter    int64 = 0
	supportedEnvs       = []string{"jaeger"}
)

type Manager interface {
	Name() string
	Start(t *testing.T)
	Close(t *testing.T)
	GetCLIConfigPath(t *testing.T) string
	GetEnvironmentResourcePath(t *testing.T, resourceName string) string
	GetTestResourcePath(t *testing.T, resourceName string) string
}

type option func(*internalManager)

type internalManager struct {
	environmentType               string
	dockerComposeNoApiFilePath    string
	dockerComposePokeshopFilePath string
	dockerProjectName             string
	pokeshopEnabled               bool
	datastoreEnabled              bool
}

func CreateAndStart(t *testing.T, options ...option) Manager {
	mutex.Lock()
	defer mutex.Unlock()

	environmentName := config.GetConfigAsEnvVars().TestEnvironment

	if !slices.Contains(supportedEnvs, environmentName) {
		t.Fatalf("environment %s not registered", environmentName)
	}

	environment := GetManager(environmentName, options...)
	environment.Start(t)

	return environment
}

func WithPokeshop() option {
	return func(im *internalManager) {
		im.pokeshopEnabled = true
	}
}

func WithDataStoreEnabled() option {
	return func(im *internalManager) {
		im.datastoreEnabled = true
	}
}

func getExecutingDir() string {
	_, filename, _, _ := runtime.Caller(0) // get file of the getExecutingDir caller
	return path.Dir(filename)
}

// Today we are assuming that the internal manager only deals with docker-compose,
// but in the future we can rename it do "dockerManager" and create another Manager to handle kubernetes environments

// This module assumes that no test will be run in parallel
// if we change this decision in the future, we will need to update the docker compose usage
// to use something like github.com/testcontainers/testcontainers-go
// (github.com/testcontainers/testcontainers-go/modules/compose in specific)

func GetManager(environmentType string, options ...option) Manager {
	currentDir := getExecutingDir()
	dockerComposeNoApiFilepath := fmt.Sprintf("%s/%s/server-setup/docker-compose-no-api.yaml", currentDir, environmentType)
	dockerComposePokeshopFilepath := fmt.Sprintf("%s/%s/server-setup/docker-compose-pokeshop.yaml", currentDir, environmentType)

	atomic.AddInt64(&envCounter, 1)

	manager := &internalManager{
		environmentType:               environmentType,
		dockerComposeNoApiFilePath:    dockerComposeNoApiFilepath,
		dockerComposePokeshopFilePath: dockerComposePokeshopFilepath,
		dockerProjectName:             fmt.Sprintf("tracetest-env-%d", envCounter),
	}

	for _, option := range options {
		option(manager)
	}

	return manager
}

func (m *internalManager) Name() string {
	return m.environmentType
}

func (m *internalManager) Start(t *testing.T) {
	args := []string{
		"compose",
		"--file", m.dockerComposeNoApiFilePath, // choose docker compose relative to the chosen environment
		"--project-name", m.dockerProjectName, // create a project name to isolate this scenario
		"up", "--detach",
	}

	if m.pokeshopEnabled {
		args = []string{
			"compose",
			"--file", m.dockerComposeNoApiFilePath, // choose docker compose relative to the chosen environment
			"--file", m.dockerComposePokeshopFilePath, // choose docker compose relative to the chosen environment
			"--project-name", m.dockerProjectName, // create a project name to isolate this scenario
			"up", "--detach",
		}
	}

	result, err := command.Exec("docker", args...)

	require.NoError(t, err)
	helpers.RequireExitCodeEqual(t, result, 0)

	// give the system 1s to get everything ready
	time.Sleep(500 * time.Millisecond)

	// wait until tracetest port is ready
	waitForPort("11633")

	if m.pokeshopEnabled {
		// wait for pokeshp services
		waitForPort("8081")

		time.Sleep(500 * time.Millisecond)
		waitForPort("8082")
	}

	if m.datastoreEnabled {
		cliConfig := m.GetCLIConfigPath(t)
		dataStorePath := m.GetEnvironmentResourcePath(t, "data-store")

		result = tracetestcli.Exec(t, fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
	}
}

func waitForPort(port string) {
	// get this source file absolute path
	// cwd changes to the dir of the test being executed
	// we need a fixed reference to calculate path to script
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic(errors.New("unable to get the current filename"))
	}
	thisDir := filepath.Dir(filename)

	script := filepath.Join(thisDir, "../../../scripts/wait-for-port.sh")
	result, err := command.Exec(script, port)

	if err != nil {
		panic(fmt.Errorf("cannot execute script: %w", err))
	}

	if result.ExitCode != 0 {
		fmt.Println("ERROR!!")
		fmt.Println("exit code: ", result.ExitCode)
		fmt.Println("output: ", result.StdOut)
		panic("script error")
	}
}

func (m *internalManager) Close(t *testing.T) {
	args := []string{
		"compose",
		"--file", m.dockerComposeNoApiFilePath, // choose docker compose relative to the chosen environment
		"--project-name", m.dockerProjectName, // choose isolated project name
		"rm",
		"--force",   // bypass removal question
		"--volumes", // remove volumes attached to this project
		"--stop",    // force containers to stop
	}

	if m.pokeshopEnabled {
		args = []string{
			"compose",
			"--file", m.dockerComposeNoApiFilePath, // choose docker compose relative to the chosen environment
			"--file", m.dockerComposePokeshopFilePath, // choose docker compose relative to the chosen environment
			"--project-name", m.dockerProjectName, // choose isolated project name
			"rm",
			"--force",   // bypass removal question
			"--volumes", // remove volumes attached to this project
			"--stop",    // force containers to stop
		}
	}

	result, err := command.Exec("docker", args...)
	require.NoError(t, err)
	helpers.RequireExitCodeEqual(t, result, 0)
}

func (m *internalManager) GetCLIConfigPath(t *testing.T) string {
	currentDir := getExecutingDir()
	return fmt.Sprintf("%s/%s/cli-config.yaml", currentDir, m.environmentType)
}

func (m *internalManager) GetEnvironmentResourcePath(t *testing.T, resourceName string) string {
	currentDir := getExecutingDir()
	return fmt.Sprintf("%s/%s/resources/%s.yaml", currentDir, m.environmentType, resourceName)
}

func (m *internalManager) GetTestResourcePath(t *testing.T, resourceName string) string {
	_, filename, _, _ := runtime.Caller(1) // get file of the GetTestResourcePath caller
	testDir := path.Dir(filename)

	return fmt.Sprintf("%s/resources/%s.yaml", testDir, resourceName)
}
