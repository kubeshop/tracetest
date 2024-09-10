// Tracetest module run https://tracetest.io/ assertions in your pipeline

// This module allow to add Tracetest capabilities to your CI pipelines. It
// currently supports either starting the Tracetest agent, or running commands via the CLI.
//
// If you want to see a real case on how this module can be used in a larger stack,
// you can check Tracetest's official pokeshop example here: https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-pokeshop

package main

import (
	"context"
	"dagger/tracetest/internal/dagger"
)

type Tracetest struct {

	// +private
	APIKey *dagger.Secret

	// +private
	Environment string

	// +private
	Organization string

	// +private
	Version string

	// +private
	Server string

	// +private
	Ctr *dagger.Container
}

func New(

	// +optional
	apiKey *dagger.Secret,

	// +optional
	environment string,

	// +optional
	organization string,

	// +optional
	// +default="https://app.tracetest.io"
	server string,

	// +optional
	// +default="latest"
	version string,

) *Tracetest {
	return &Tracetest{
		APIKey:       apiKey,
		Environment:  environment,
		Organization: organization,
		Version:      version,
		Server:       server,
		Ctr:          dag.Container().From("kubeshop/tracetest:" + version),
	}
}

// Retruns a Tracetest container configured to be used
// as a Dagger service
func (m *Tracetest) Agent() *dagger.Container {
	return m.Ctr.With(func(c *dagger.Container) *dagger.Container {
		if m.APIKey != nil {
			c = c.WithSecretVariable("TRACETEST_API_KEY", m.APIKey)
		}
		if m.Environment != "" {
			c = c.WithEnvVariable("TRACETEST_ENVIRONMENT_ID", m.Environment)
		}
		return c
	}).
		WithExposedPort(4317).WithExposedPort(4318)
}

func (m *Tracetest) WithService(name string, svc *dagger.Service) *Tracetest {
	m.Ctr = m.Ctr.WithServiceBinding(name, svc)
	return m
}

// Runs Tracetest CLI commands.
func (m *Tracetest) CLI(ctx context.Context) *dagger.Container {

	c := m.Ctr.WithEntrypoint([]string{"tracetest"})

	if m.APIKey != nil {
		c = c.WithSecretVariable("TRACETEST_API_KEY", m.APIKey)
	}

	if m.Environment != "" {
		c = c.WithEnvVariable("TRACETEST_ENVIRONMENT_ID", m.Environment)
	}

	args := []string{"configure", "-s", m.Server}

	if m.Organization != "" {
		args = append(args, "--organization", m.Organization)
	}

	if m.Environment != "" {
		args = append(args, "--environment", m.Environment)
	}

	if m.APIKey != nil {
		plainKey, _ := m.APIKey.Plaintext(ctx)
		args = append(args, "--token", plainKey)
	}

	c = c.WithExec(args, dagger.ContainerWithExecOpts{UseEntrypoint: true})
	return c

}
