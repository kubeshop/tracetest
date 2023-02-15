package provisioning_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	t.Run("Inexistent", func(t *testing.T) {
		t.Parallel()

		provisioner := provisioning.New(&testdb.MockRepository{})

		err := provisioner.FromFile("notexists.yaml")
		assert.ErrorContains(t, err, "cannot read provisioning file 'notexists.yaml'")
	})
}
