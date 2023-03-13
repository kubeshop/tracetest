package provisioning_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	db := testmock.MustGetRawTestingDatabase()

	t.Run("Inexistent", func(t *testing.T) {
		t.Parallel()

		configDB := configresource.NewRepository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)

		provisioner := provisioning.New(&testdb.MockRepository{}, configDB)

		err := provisioner.FromFile("notexists.yaml")
		assert.ErrorIs(t, err, provisioning.ErrFileNotExists)
	})
}

func TestFromEnv(t *testing.T) {
	db := testmock.MustGetRawTestingDatabase()

	t.Run("Empty", func(t *testing.T) {
		configDB := configresource.NewRepository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)

		provisioner := provisioning.New(&testdb.MockRepository{}, configDB)

		err := provisioner.FromEnv()
		assert.ErrorIs(t, err, provisioning.ErrEnvEmpty)
	})

	t.Run("InvalidData", func(t *testing.T) {
		configDB := configresource.NewRepository(
			testmock.MustCreateRandomMigratedDatabase(db),
		)

		os.Setenv("TRACETEST_PROVISIONING", "this is not base64")

		provisioner := provisioning.New(&testdb.MockRepository{}, configDB)

		err := provisioner.FromEnv()
		assert.ErrorContains(t, err, "cannot decode env variable TRACETEST_PROVISIONING")
	})

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			configDB := configresource.NewRepository(
				testmock.MustCreateRandomMigratedDatabase(db),
			)

			t.Run("FromEnv", func(t *testing.T) {
				c := c

				expectedDataStore := model.DataStore{
					IsDefault: true,
					Name:      string(c.dsType),
					Type:      c.dsType,
					Values:    c.values,
				}

				mockRepo := &testdb.MockRepository{}
				provisioner := provisioning.New(mockRepo, configDB)

				mockRepo.
					On("CreateDataStore", expectedDataStore).
					Return(expectedDataStore, nil)

				fcontents, err := os.ReadFile(c.file)
				if err != nil {
					panic(err)
				}

				encoded := base64.StdEncoding.EncodeToString(fcontents)
				os.Setenv("TRACETEST_PROVISIONING", encoded)

				err = provisioner.FromEnv()
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			})
		})
	}

}
