package variable

import "os"

type VariableProvider interface {
	GetVariable(string) (string, error)
}

type EnvironmentVariableProvider struct{}

func (p EnvironmentVariableProvider) GetVariable(name string) (string, error) {
	return os.Getenv(name), nil
}
