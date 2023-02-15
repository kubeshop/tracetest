package provisioning

import (
	"fmt"
	"io/ioutil"

	"github.com/kubeshop/tracetest/server/model"
)

func New(db model.Repository) provisioner {
	return provisioner{db}
}

type provisioner struct {
	b model.Repository
}

func (p provisioner) FromFile(path string) error {
	_, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read provisioning file '%s'", path)
	}

	return nil
}
