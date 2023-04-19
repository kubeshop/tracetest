package datastoreresource_test

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testmock"
	datastore "github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"github.com/stretchr/testify/require"
)

func compareJSONDataStores(t require.TestingT, operation rmtests.Operation, firstValue, secondValue string) {
	if operation == rmtests.OperationUpdateSuccess {
		require.JSONEq(t, firstValue, secondValue)
		return
	}

	expected := rmtests.RemoveFieldFromJSONResource("createdAt", firstValue)
	obtained := rmtests.RemoveFieldFromJSONResource("createdAt", secondValue)

	require.JSONEq(t, expected, obtained)
}

func registerManagerFn(router *mux.Router) resourcemanager.Manager {
	db := testmock.CreateMigratedDatabase()
	dataStoreRepository := datastore.NewRepository(db)

	manager := resourcemanager.New[datastore.DataStore](
		datastore.ResourceName,
		datastore.ResourceNamePlural,
		dataStoreRepository,
		resourcemanager.WithOperations(datastore.Operations...),
		resourcemanager.WithIDGen(id.GenerateID),
	)
	manager.RegisterRoutes(router)

	return manager
}

func cleanup(t *testing.T, manager resourcemanager.Manager) {
	repository := manager.Handler().(*datastore.Repository)
	err := repository.Close()
	require.NoError(t, err)
}

func getScenarioPreparation(sample datastore.DataStore) func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
	return func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
		repository := manager.Handler().(*datastore.Repository)

		if op == rmtests.OperationGetSuccess {
			// on get scenario we need to have one data store registered
			repository.Update(context.TODO(), sample)
		}
	}
}

func TestDataStoreResource_AWSXRay(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeAwsXRay,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			AwsXRay: &datastore.AWSXRayConfig{
				Region:          "some-region",
				AccessKeyID:     "some-access-key",
				SecretAccessKey: "some-secret-access-key",
				SessionToken:    "some-session-token",
				UseDefaultAuth:  true,
			},
		},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "awsxray",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"awsXRay": {
					"region": "some-region",
					"accessKeyID": "some-access-key",
					"secretAccessKey": "some-secret-access-key",
					"sessionToken": "some-session-token",
					"useDefaultAuth": true
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "awsxray",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"awsXRay": {
					"region": "some-region-updated",
					"accessKeyID": "some-access-key-updated",
					"secretAccessKey": "some-access-key-updated",
					"sessionToken": "some-session-token-updated",
					"useDefaultAuth": true
				}
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestDataStoreResource_ElasticAPM(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeElasticAPM,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			ElasticApm: &datastore.ElasticSearchConfig{
				Addresses:          []string{"1.2.3.4"},
				Username:           "some-user",
				Password:           "some-password",
				Index:              "an-index",
				Certificate:        "certificate.cert",
				InsecureSkipVerify: true,
			},
		},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "elasticapm",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"elasticAPM": {
					"addresses": ["1.2.3.4"],
					"username": "some-user",
					"password": "some-password",
					"index": "an-index",
					"certificate": "certificate.cert",
					"insecureSkipVerify": true
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "elasticapm",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"elasticAPM": {
					"addresses": ["4.3.2.1"],
					"username": "some-user-updated",
					"password": "some-password-updated",
					"index": "an-index-updated",
					"certificate": "certificate.cert-updated",
					"insecureSkipVerify": true
				}
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestDataStoreResource_Jaeger(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeJaeger,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			Jaeger: &datastore.GRPCClientSettings{
				Endpoint: "some-endpoint",
				TLS: &datastore.TLS{
					Insecure: true,
				},
			},
		},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "jaeger",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"jaeger": {
					"endpoint": "some-endpoint",
					"tls": {
						"insecure": true
					}
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "jaeger",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"jaeger": {
					"endpoint": "some-endpoint",
					"tls": {
						"insecure": true
					}
				}
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestDataStoreResource_OTLP(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeOTLP,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values:    datastore.DataStoreValues{},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "otlp",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z"
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "otlp",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z"
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestDataStoreResource_OpenSearch(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeOpenSearch,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			ElasticApm: &datastore.ElasticSearchConfig{
				Addresses:          []string{"1.2.3.4"},
				Username:           "some-user",
				Password:           "some-password",
				Index:              "an-index",
				Certificate:        "certificate.cert",
				InsecureSkipVerify: true,
			},
		},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "opensearch",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"elasticAPM": {
					"addresses": ["1.2.3.4"],
					"username": "some-user",
					"password": "some-password",
					"index": "an-index",
					"certificate": "certificate.cert",
					"insecureSkipVerify": true
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "opensearch",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"elasticAPM": {
					"addresses": ["4.3.2.1"],
					"username": "some-user-updated",
					"password": "some-password-updated",
					"index": "an-index-updated",
					"certificate": "certificate.cert-updated",
					"insecureSkipVerify": true
				}
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestDataStoreResource_SignalFX(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeSignalFX,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			SignalFx: &datastore.SignalFXConfig{
				Realm: "some-realm",
				Token: "some-token",
			},
		},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "signalfx",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"signalFX": {
					"realm": "some-realm",
					"token": "some-token"
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "signalfx",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"signalFX": {
					"realm": "some-realm-updated",
					"token": "some-token-updated"
				}
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}

func TestDataStoreResource_Tempo(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeTempo,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			Tempo: &datastore.MultiChannelClientConfig{
				Type: datastore.MultiChannelClientTypeGRPC,
				Grpc: &datastore.GRPCClientSettings{
					Endpoint: "some-endpoint",
					TLS: &datastore.TLS{
						Insecure: true,
					},
				},
			},
		},
	}

	testSpec := rmtests.ResourceTypeTest{
		ResourceTypeSingular: datastore.ResourceName,
		ResourceTypePlural:   datastore.ResourceNamePlural,
		RegisterManagerFn:    registerManagerFn,
		Prepare:              getScenarioPreparation(sample),
		Cleanup:              cleanup,
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "tempo",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"tempo": {
					"type": "grpc",
					"grpc": {
						"endpoint": "some-endpoint",
						"tls": {
							"insecure": true
						}
					}
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "another data store",
				"type": "tempo",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"tempo": {
					"type": "http",
					"http": {
						"url": "some-url",
						"headers": {
							"authorization": "something"
						}
					}
				}
			}
		}`,
	}

	excludedOperations := rmtests.ExcludeOperations(rmtests.OperationUpdateNotFound, rmtests.OperationGetNotFound)
	jsonComparer := rmtests.JSONComparer(compareJSONDataStores)

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer)
}
