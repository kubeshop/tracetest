package environment

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var (
	mutex               = sync.Mutex{}
	envCounter    int64 = 0
	defaultEnv          = "jaeger"
	supportedEnvs       = []string{"jaeger"}
)

type Manager interface {
	Name() string
	Start(t *testing.T)
	Close(t *testing.T)
	GetCLIConfigPath(t *testing.T) string
	GetManisfestResourcePath(t *testing.T, manifestName string) string
}

type internalManager struct {
	environmentType       string
	dockerComposeFilePath string
	dockerProjectName     string
}

func CreateAndStart(t *testing.T) Manager {
	t.Helper()

	mutex.Lock()
	defer mutex.Unlock()

	environmentName := os.Getenv("TEST_ENVIRONMENT")

	if environmentName == "" {
		environmentName = defaultEnv
	}

	if !slices.Contains(supportedEnvs, environmentName) {
		t.Fatalf("environment %s not registered", environmentName)
	}

	environment := GetManager(environmentName)
	environment.Start(t)

	return environment
}

func getExecutingDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

// Today we are assuming that the internal manager only deals with docker-compose,
// but in the future we can rename it do "dockerManager" and create another Manager to handle kubernetes environments

// This module assumes that no test will be run in parallel
// if we change this decision in the future, we will need to update the docker compose usage
// to use something like github.com/testcontainers/testcontainers-go
// (github.com/testcontainers/testcontainers-go/modules/compose in specific)

func GetManager(environmentType string) Manager {
	currentDir := getExecutingDir()
	dockerComposeFilepath := fmt.Sprintf("%s/%s/server-setup/docker-compose.yaml", currentDir, environmentType)

	atomic.AddInt64(&envCounter, 1)

	return &internalManager{
		environmentType:       environmentType,
		dockerComposeFilePath: dockerComposeFilepath,
		dockerProjectName:     fmt.Sprintf("tracetest-env-%d", envCounter),
	}
}

func (m *internalManager) Name() string {
	return m.environmentType
}

func (m *internalManager) Start(t *testing.T) {
	t.Helper()

	result, err := command.Exec(
		"docker", "compose",
		"--file", m.dockerComposeFilePath, // choose docker compose relative to the chosen environment
		"--project-name", m.dockerProjectName, // create a project name to isolate this scenario
		"up", "--detach")

	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)

	// TODO: think in a better way to assure readiness for Tracetest
	time.Sleep(1000 * time.Millisecond)
}

func (m *internalManager) Close(t *testing.T) {
	t.Helper()

	result, err := command.Exec(
		"docker", "compose",
		"--file", m.dockerComposeFilePath, // choose docker compose relative to the chosen environment
		"--project-name", m.dockerProjectName, // choose isolated project name
		"rm",
		"--force",   // bypass removal question
		"--volumes", // remove volumes attached to this project
		"--stop",    // force containers to stop
	)
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
}

func (m *internalManager) GetCLIConfigPath(t *testing.T) string {
	currentDir := getExecutingDir()
	return fmt.Sprintf("%s/%s/cli-config.yaml", currentDir, m.environmentType)
}

func (m *internalManager) GetManisfestResourcePath(t *testing.T, manifestName string) string {
	currentDir := getExecutingDir()
	return fmt.Sprintf("%s/%s/resources/%s.yaml", currentDir, m.environmentType, manifestName)
}
