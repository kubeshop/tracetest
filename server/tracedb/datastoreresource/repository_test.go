package datastoreresource_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	datastore "github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"github.com/stretchr/testify/require"
)

var (
	excludedOperations = rmtests.ExcludeOperations(
		rmtests.OperationUpdateNotFound,
		rmtests.OperationGetNotFound,
		rmtests.OperationDeleteNotFound,
	)
	operationsWithoutPostAssert = rmtests.IgnorePostAssertForOperations(rmtests.OperationDeleteSuccess)
	jsonComparer                = rmtests.JSONComparer(compareJSONDataStores)
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

func registerManagerFn(router *mux.Router, db *sql.DB) resourcemanager.Manager {
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
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "awsxray",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"awsxray": {
					"region": "some-region",
					"accessKeyId": "some-access-key",
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
				"awsxray": {
					"region": "some-region-updated",
					"accessKeyId": "some-access-key-updated",
					"secretAccessKey": "some-access-key-updated",
					"sessionToken": "some-session-token-updated",
					"useDefaultAuth": true
				}
			}
		}`,
	}

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
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
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "elasticapm",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"elasticapm": {
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
				"elasticapm": {
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

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
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

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
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

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
}

func TestDataStoreResource_OpenSearch(t *testing.T) {
	sample := datastore.DataStore{
		ID:        "current",
		Name:      "default",
		Type:      datastore.DataStoreTypeOpenSearch,
		Default:   true,
		CreatedAt: "2023-03-09T17:53:10.256383Z",
		Values: datastore.DataStoreValues{
			OpenSearch: &datastore.ElasticSearchConfig{
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
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "opensearch",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"opensearch": {
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
				"opensearch": {
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

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
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
		SampleJSON: `{
			"type": "DataStore",
			"spec": {
				"id": "current",
				"name": "default",
				"type": "signalfx",
				"default": true,
				"createdAt": "2023-03-09T17:53:10.256383Z",
				"signalfx": {
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
				"signalfx": {
					"realm": "some-realm-updated",
					"token": "some-token-updated"
				}
			}
		}`,
	}

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
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

	rmtests.TestResourceType(t, testSpec, excludedOperations, jsonComparer, operationsWithoutPostAssert)
}
