package environment

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

var (
	mutex         = sync.Mutex{}
	defaultEnv    = "jaeger"
	supportedEnvs = []string{"jaeger", "tempo"}
)

type Manager interface {
	Start(t *testing.T)
	Close(t *testing.T)
	GetCLIConfig(t *testing.T) string
}

type internalManager struct {
	environmentType string
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

	environment := &internalManager{
		environmentType: environmentName,
	}

	environment.Start(t)

	return environment
}

func getExecutingDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

// TODO: this module assumes that no test will be run in parallel
// if we change this decision in the future, we will need to update the docker compose usage
// to use something like github.com/testcontainers/testcontainers-go
// (github.com/testcontainers/testcontainers-go/modules/compose in specific)

func (m *internalManager) Start(t *testing.T) {
	t.Helper()

	currentDir := getExecutingDir()
	dockerComposeFilepath := fmt.Sprintf("%s/%s/server-setup/docker-compose.yaml", currentDir, m.environmentType)

	result, err := command.Exec("docker", "compose", "-f", dockerComposeFilepath, "up", "-d")
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)

	// TODO: think in a better way to assure readiness for Tracetest
	time.Sleep(1 * time.Second)
}

func (m *internalManager) Close(t *testing.T) {
	t.Helper()

	currentDir := getExecutingDir()
	dockerComposeFilepath := fmt.Sprintf("%s/%s/server-setup/docker-compose.yaml", currentDir, m.environmentType)

	result, err := command.Exec("docker", "compose", "-f", dockerComposeFilepath, "rm", "--force", "--volumes", "--stop")
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
}

func (m *internalManager) GetCLIConfig(t *testing.T) string {
	currentDir := getExecutingDir()
	return fmt.Sprintf("%s/%s/cli-config.yaml", currentDir, m.environmentType)
}

func (m *internalManager) GetManisfestResource(t *testing.T, manifestName string) string {
	currentDir := getExecutingDir()
	return fmt.Sprintf("%s/%s/resources/%s.yaml", currentDir, m.environmentType, manifestName)
}
