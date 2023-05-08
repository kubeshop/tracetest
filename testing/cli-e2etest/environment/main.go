package environment

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

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
}

type internalManager struct {
	environmentType string
}

func CreateAndStart(t *testing.T) Manager {
	t.Helper()

	mutex.Lock()
	defer mutex.Unlock()

	environmentName := os.Getenv("TEST_ENV")

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

func getDockerComposePath(environmentType string) string {
	currentDir := getExecutingDir()

	return fmt.Sprintf("%s/%s/server-setup/docker-compose.yaml", currentDir, environmentType)
}

// TODO: this module assumes that no test will be run in parallel
// if we change this decision in the future, we will need to update the docker compose usage
// to use something like github.com/testcontainers/testcontainers-go
// (github.com/testcontainers/testcontainers-go/modules/compose in specific)

func (m *internalManager) Start(t *testing.T) {
	t.Helper()

	dockerComposeFilepath := getDockerComposePath(m.environmentType)

	_, exitCode, err := command.Exec("docker", "compose", "-f", dockerComposeFilepath, "up", "-d")
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}

func (m *internalManager) Close(t *testing.T) {
	t.Helper()

	dockerComposeFilepath := getDockerComposePath(m.environmentType)

	_, exitCode, err := command.Exec("docker", "compose", "-f", dockerComposeFilepath, "stop")
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)

	_, exitCode, err = command.Exec("docker", "compose", "-f", dockerComposeFilepath, "rm")
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)
}
