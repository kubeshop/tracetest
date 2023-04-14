package provisioning_test

import (
	"context"
	"database/sql"
	"encoding/base64"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config/configresource"
	"github.com/kubeshop/tracetest/server/config/demoresource"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestFromFile(t *testing.T) {
	t.Run("Inexistent", func(t *testing.T) {
		t.Parallel()
		provisioner := provisioning.New()

		err := provisioner.FromFile("notexists.yaml")
		assert.ErrorIs(t, err, provisioning.ErrFileNotExists)
	})

	db := testmock.MustGetRawTestingDatabase()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("FromFile", func(t *testing.T) {
				c := c
				t.Parallel()

				f := setup(db)

				err := f.provisioner.FromFile(c.file)
				assert.NoError(t, err)

				f.assert(t, c.expectations)
			})
		})
	}
}

func TestFromEnv(t *testing.T) {
	db := testmock.MustGetRawTestingDatabase()

	t.Run("Empty", func(t *testing.T) {
		provisioner := provisioning.New()

		err := provisioner.FromEnv()
		assert.ErrorIs(t, err, provisioning.ErrEnvEmpty)
	})

	t.Run("InvalidData", func(t *testing.T) {
		os.Setenv("TRACETEST_PROVISIONING", "this is not base64")

		provisioner := provisioning.New()

		err := provisioner.FromEnv()
		assert.ErrorContains(t, err, "cannot decode env variable TRACETEST_PROVISIONING")
	})

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Run("FromEnv", func(t *testing.T) {
				f := setup(db)

				fcontents, err := os.ReadFile(c.file)
				if err != nil {
					panic(err)
				}

				encoded := base64.StdEncoding.EncodeToString(fcontents)
				os.Setenv("TRACETEST_PROVISIONING", encoded)

				err = f.provisioner.FromEnv()
				assert.NoError(t, err)

				f.assert(t, c.expectations)

			})
		})
	}

}

type expectations struct {
	dataStore      *model.DataStore
	config         *configresource.Config
	pollingprofile *pollingprofile.PollingProfile
	demos          []demoresource.Demo
}

type provisioningFixture struct {
	provisioner     *provisioning.Provisioner
	configs         *configresource.Repository
	pollingProfiles *pollingprofile.Repository
	demos           *demoresource.Repository
	dataStores      model.DataStoreRepository
}

func (f provisioningFixture) assert(t *testing.T, expected expectations) {
	if expected.dataStore != nil {
		actual, err := f.dataStores.DefaultDataStore(context.TODO())
		require.NoError(t, err)

		// ignore ID for assertion
		expected.dataStore.ID = "0"
		actual.ID = "0"

		// ignore createdAt for assertion
		expected.dataStore.CreatedAt = actual.CreatedAt

		assert.Equal(t, *expected.dataStore, actual)
	}

	if expected.config != nil {
		actual := f.configs.Current(context.TODO())

		// ignore ID for assertion
		expected.config.ID = "0"
		actual.ID = "0"

		// ignore Name for assertion
		expected.config.Name = ""
		actual.Name = ""

		assert.Equal(t, *expected.config, actual)
	}

	if expected.pollingprofile != nil {
		actual := f.pollingProfiles.GetDefault(context.TODO())

		// ignore ID for assertion
		expected.pollingprofile.ID = "0"
		actual.ID = "0"

		assert.Equal(t, *expected.pollingprofile, actual)
	}

	if demosCount := len(expected.demos); demosCount > 0 {
		list, err := f.demos.List(context.TODO(), demosCount, 0, "", "", "")
		require.NoError(t, err)

		require.Equal(t, demosCount, len(list))

		for i, actual := range list {
			// ignore ID for assertion
			actual.ID = "0"
			expected.demos[i].ID = "0"

			assert.Equal(t, expected.demos[i], actual)
		}
	}
}

func setup(db *sql.DB) provisioningFixture {
	testDB, err := testdb.Postgres(testdb.WithDB(db))
	if err != nil {
		panic(err)
	}
	f := provisioningFixture{
		configs:         configresource.NewRepository(db),
		pollingProfiles: pollingprofile.NewRepository(db),
		demos:           demoresource.NewRepository(db),
		dataStores:      testDB,
	}

	configManager := resourcemanager.New[configresource.Config](
		configresource.ResourceName,
		configresource.ResourceNamePlural,
		f.configs,
		resourcemanager.WithOperations(configresource.Operations...),
	)

	pollingProfilesManager := resourcemanager.New[pollingprofile.PollingProfile](
		pollingprofile.ResourceName,
		pollingprofile.ResourceNamePlural,
		f.pollingProfiles,
		resourcemanager.WithOperations(pollingprofile.Operations...),
	)

	demoManager := resourcemanager.New[demoresource.Demo](
		demoresource.ResourceName,
		demoresource.ResourceNamePlural,
		f.demos,
		resourcemanager.WithOperations(demoresource.Operations...),
	)

	dataStoreManager := resourcemanager.New[testdb.DataStoreResource](
		testdb.DataStoreResourceName,
		testdb.DataStoreResourceNamePlural,
		testdb.NewDataStoreResourceProvisioner(f.dataStores),
		resourcemanager.WithOperations(resourcemanager.OperationNoop),
	)

	f.provisioner = provisioning.New(provisioning.WithResourceProvisioners(
		configManager,
		pollingProfilesManager,
		demoManager,
		dataStoreManager,
	))

	return f
}

var cases = []struct {
	name         string
	file         string
	expectations expectations
}{
	{
		name: "AllSettings",
		file: "./testdata/all_settings.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "Jaeger",
				IsDefault: true,
				Type:      model.DataStoreTypeJaeger,
				Values: model.DataStoreValues{
					Jaeger: &configgrpc.GRPCClientSettings{
						Endpoint:   "jaeger-query:16685",
						TLSSetting: configtls.TLSClientSetting{Insecure: true},
					},
				},
			},
			config: &configresource.Config{
				AnalyticsEnabled: true,
			},
			pollingprofile: &pollingprofile.PollingProfile{
				Name:     "Custom Profile",
				Default:  true,
				Strategy: pollingprofile.Periodic,
				Periodic: &pollingprofile.PeriodicPollingConfig{
					Timeout:    "2h",
					RetryDelay: "30m",
				},
			},
			demos: []demoresource.Demo{
				{
					Name:    "pokeshop",
					Type:    demoresource.DemoTypePokeshop,
					Enabled: true,
					Pokeshop: &demoresource.PokeshopDemo{
						HTTPEndpoint: "http://localhost/api",
						GRPCEndpoint: "localhost:8080",
					},
				},
				{
					Name:    "otel",
					Type:    demoresource.DemoTypeOpentelemetryStore,
					Enabled: true,
					OpenTelemetryStore: &demoresource.OpenTelemetryStoreDemo{
						FrontendEndpoint: "http://frontend:8080/",
					},
				},
			},
		},
	},
	{
		name: "JaegerGRPC",
		file: "./testdata/jaeger_grpc.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "Jaeger",
				IsDefault: true,
				Type:      model.DataStoreTypeJaeger,
				Values: model.DataStoreValues{
					Jaeger: &configgrpc.GRPCClientSettings{
						Endpoint:   "jaeger-query:16685",
						TLSSetting: configtls.TLSClientSetting{Insecure: true},
					},
				},
			},
		},
	},
	{
		name: "TempoGRPC",
		file: "./testdata/tempo_grpc.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "Tempo (gRPC)",
				IsDefault: true,
				Type:      model.DataStoreTypeTempo,
				Values: model.DataStoreValues{
					Tempo: &model.BaseClientConfig{
						Grpc: configgrpc.GRPCClientSettings{
							Endpoint:   "tempo:9095",
							TLSSetting: configtls.TLSClientSetting{Insecure: true},
						},
					},
				},
			},
		},
	},
	{
		name: "TempoHTTP",
		file: "./testdata/tempo_http.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "Tempo (HTTP)",
				IsDefault: true,
				Type:      model.DataStoreTypeTempo,
				Values: model.DataStoreValues{
					Tempo: &model.BaseClientConfig{
						Http: model.HttpClientConfig{
							Url:        "tempo:80",
							TLSSetting: configtls.TLSClientSetting{Insecure: true},
						},
					},
				},
			},
		},
	},
	{
		name: "OpenSearch",
		file: "./testdata/opensearch.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "OpenSearch",
				IsDefault: true,
				Type:      model.DataStoreTypeOpenSearch,
				Values: model.DataStoreValues{
					OpenSearch: &model.ElasticSearchDataStoreConfig{
						Addresses: []string{"http://opensearch:9200"},
						Index:     "traces",
					},
				},
			},
		},
	},
	{
		name: "SignalFX",
		file: "./testdata/signalfx.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "SignalFX",
				IsDefault: true,
				Type:      model.DataStoreTypeSignalFX,
				Values: model.DataStoreValues{
					SignalFx: &model.SignalFXDataStoreConfig{
						Token: "thetoken",
						Realm: "us1",
					},
				},
			},
		},
	},
	{
		name: "ElasitcAPM",
		file: "./testdata/elastic_apm.yaml",
		expectations: expectations{
			dataStore: &model.DataStore{
				Name:      "elastic APM",
				IsDefault: true,
				Type:      model.DataStoreTypeElasticAPM,
				Values: model.DataStoreValues{
					ElasticApm: &model.ElasticSearchDataStoreConfig{
						Addresses:          []string{"https://es01:9200"},
						Username:           "elastic",
						Password:           "changeme",
						Index:              "traces-apm-default",
						InsecureSkipVerify: true,
					},
				},
			},
		},
	},
}
