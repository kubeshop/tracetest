package provisioning_test

import (
	"context"
	"database/sql"
	"encoding/base64"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/config/demo"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/provisioning"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromFile(t *testing.T) {
	t.Run("Inexistent", func(t *testing.T) {
		t.Parallel()
		provisioner := provisioning.New()

		err := provisioner.FromFile("notexists.yaml")
		assert.ErrorIs(t, err, provisioning.ErrFileNotExists)
	})

	db := testmock.CreateMigratedDatabase()
	defer db.Close()

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
	db := testmock.CreateMigratedDatabase()
	defer db.Close()

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
	dataStore      *datastore.DataStore
	config         *config.Config
	pollingprofile *pollingprofile.PollingProfile
	demos          []demo.Demo
}

type provisioningFixture struct {
	provisioner     *provisioning.Provisioner
	configs         *config.Repository
	pollingProfiles *pollingprofile.Repository
	demos           *demo.Repository
	dataStores      *datastore.Repository
}

func (f provisioningFixture) assert(t *testing.T, expected expectations) {
	if expected.dataStore != nil {
		actual, err := f.dataStores.Current(context.TODO())
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
	f := provisioningFixture{
		configs:         config.NewRepository(db),
		pollingProfiles: pollingprofile.NewRepository(db),
		demos:           demo.NewRepository(db),
		dataStores:      datastore.NewRepository(db),
	}

	configManager := resourcemanager.New[config.Config](
		config.ResourceName,
		config.ResourceNamePlural,
		f.configs,
		resourcemanager.WithOperations(config.Operations...),
	)

	pollingProfilesManager := resourcemanager.New[pollingprofile.PollingProfile](
		pollingprofile.ResourceName,
		pollingprofile.ResourceNamePlural,
		f.pollingProfiles,
		resourcemanager.WithOperations(pollingprofile.Operations...),
	)

	demoManager := resourcemanager.New[demo.Demo](
		demo.ResourceName,
		demo.ResourceNamePlural,
		f.demos,
		resourcemanager.WithOperations(demo.Operations...),
	)

	dataStoreManager := resourcemanager.New[datastore.DataStore](
		datastore.ResourceName,
		datastore.ResourceNamePlural,
		f.dataStores,
		resourcemanager.WithOperations(datastore.Operations...),
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
			dataStore: &datastore.DataStore{
				Name:    "Jaeger",
				Default: true,
				Type:    datastore.DataStoreTypeJaeger,
				Values: datastore.DataStoreValues{
					Jaeger: &datastore.GRPCClientSettings{
						Endpoint: "jaeger-query:16685",
						TLS:      &datastore.TLS{Insecure: true},
					},
				},
			},
			config: &config.Config{
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
			demos: []demo.Demo{
				{
					Name:    "otel",
					Type:    demo.DemoTypeOpentelemetryStore,
					Enabled: true,
					OpenTelemetryStore: &demo.OpenTelemetryStoreDemo{
						FrontendEndpoint: "http://frontend:8080/",
					},
				},
				{
					Name:    "pokeshop",
					Type:    demo.DemoTypePokeshop,
					Enabled: true,
					Pokeshop: &demo.PokeshopDemo{
						HTTPEndpoint: "http://localhost/api",
						GRPCEndpoint: "localhost:8080",
					},
				},
			},
		},
	},
	{
		name: "JaegerGRPC",
		file: "./testdata/jaeger_grpc.yaml",
		expectations: expectations{
			dataStore: &datastore.DataStore{
				Name:    "Jaeger",
				Default: true,
				Type:    datastore.DataStoreTypeJaeger,
				Values: datastore.DataStoreValues{
					Jaeger: &datastore.GRPCClientSettings{
						Endpoint: "jaeger-query:16685",
						TLS:      &datastore.TLS{Insecure: true},
					},
				},
			},
		},
	},
	{
		name: "TempoGRPC",
		file: "./testdata/tempo_grpc.yaml",
		expectations: expectations{
			dataStore: &datastore.DataStore{
				Name:    "Tempo (gRPC)",
				Default: true,
				Type:    datastore.DataStoreTypeTempo,
				Values: datastore.DataStoreValues{
					Tempo: &datastore.MultiChannelClientConfig{
						Grpc: &datastore.GRPCClientSettings{
							Endpoint: "tempo:9095",
							TLS:      &datastore.TLS{Insecure: true},
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
			dataStore: &datastore.DataStore{
				Name:    "Tempo (HTTP)",
				Default: true,
				Type:    datastore.DataStoreTypeTempo,
				Values: datastore.DataStoreValues{
					Tempo: &datastore.MultiChannelClientConfig{
						Http: &datastore.HttpClientConfig{
							Url: "tempo:80",
							TLS: &datastore.TLS{Insecure: true},
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
			dataStore: &datastore.DataStore{
				Name:    "OpenSearch",
				Default: true,
				Type:    datastore.DataStoreTypeOpenSearch,
				Values: datastore.DataStoreValues{
					OpenSearch: &datastore.ElasticSearchConfig{
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
			dataStore: &datastore.DataStore{
				Name:    "SignalFX",
				Default: true,
				Type:    datastore.DataStoreTypeSignalFX,
				Values: datastore.DataStoreValues{
					SignalFx: &datastore.SignalFXConfig{
						Token: "thetoken",
						Realm: "us1",
					},
				},
			},
		},
	},
	{
		name: "ElasticAPM",
		file: "./testdata/elastic_apm.yaml",
		expectations: expectations{
			dataStore: &datastore.DataStore{
				Name:    "elastic APM",
				Default: true,
				Type:    datastore.DataStoreTypeElasticAPM,
				Values: datastore.DataStoreValues{
					ElasticApm: &datastore.ElasticSearchConfig{
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
