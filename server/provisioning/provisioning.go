package provisioning

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

func WithResourceProvisioners(provs ...resourcemanager.Provisioner) option {
	return func(p *Provisioner) {
		for _, prov := range provs {
			p.AddResourceProvisioner(prov)
		}
	}
}

type option func(p *Provisioner)

func New(opts ...option) *Provisioner {
	p := &Provisioner{}
	for _, opt := range opts {
		opt(p)
	}

	return p
}

type Provisioner struct {
	provisioners []resourcemanager.Provisioner
}

func (p *Provisioner) AddResourceProvisioner(prov resourcemanager.Provisioner) {
	p.provisioners = append(p.provisioners, prov)
}

var (
	ErrEnvEmpty = errors.New("cannot read provisioning from env variable TRACETEST_PROVISIONING: variable is empty")
)

func (p Provisioner) FromEnv() error {
	envVar := os.Getenv("TRACETEST_PROVISIONING")
	if envVar == "" {
		return ErrEnvEmpty
	}

	data, err := base64.StdEncoding.DecodeString(envVar)
	if err != nil {
		return fmt.Errorf("cannot decode env variable TRACETEST_PROVISIONING: %w", err)
	}
	return p.do(data)
}

var ErrFileNotExists = errors.New("provisioning file does not exists")

func (p Provisioner) FromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists
		}
		return fmt.Errorf("cannot read provisioning file '%s'", path)
	}

	return p.do(data)
}

func (p Provisioner) do(data []byte) error {

	ctx := context.Background()

	d := yaml.NewDecoder(bytes.NewBuffer(data))
	for {
		var rawYaml map[string]any
		err := d.Decode(&rawYaml)

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("cannot unmarshal yaml: %w", err)
		}

		if rawYaml == nil {
			continue
		}

		success := false
		for _, p := range p.provisioners {
			err := p.Provision(ctx, rawYaml)
			if errors.Is(err, resourcemanager.ErrTypeNotSupported) {
				continue
			}
			if err != nil {
				return fmt.Errorf("cannot provision resource from yaml: %w", err)
			}
			success = true
		}

		if !success {
			return fmt.Errorf("invalid resource type from yaml")
		}
	}
	return nil
}
