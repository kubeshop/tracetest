package demo_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/demo"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestPokeshopDemoResource(t *testing.T) {
	sampleDemo := demo.Demo{
		ID:      "1",
		Name:    "dev",
		Type:    demo.DemoTypePokeshop,
		Enabled: true,
		Pokeshop: &demo.PokeshopDemo{
			HTTPEndpoint: "http://dev-endpoint:1234",
			GRPCEndpoint: "dev-grpc:9091",
			KafkaBroker:  "dev-kafka:9092",
		},
	}

	secondSampleDemo := demo.Demo{
		ID:      "2",
		Name:    "staging",
		Type:    demo.DemoTypePokeshop,
		Enabled: true,
		Pokeshop: &demo.PokeshopDemo{
			HTTPEndpoint: "http://stg-endpoint:1234",
			GRPCEndpoint: "stg-grpc:9091",
			KafkaBroker:  "stg-kafka:9092",
		},
	}

	thirdSampleDemo := demo.Demo{
		ID:      "3",
		Name:    "production",
		Type:    demo.DemoTypePokeshop,
		Enabled: true,
		Pokeshop: &demo.PokeshopDemo{
			HTTPEndpoint: "http://prod-endpoint:1234",
			GRPCEndpoint: "prod-grpc:9091",
			KafkaBroker:  "prod-kafka:9092",
		},
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: demo.ResourceName,
		ResourceTypePlural:   demo.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			demoRepository := demo.NewRepository(db)

			manager := resourcemanager.New[demo.Demo](
				demo.ResourceName,
				demo.ResourceNamePlural,
				demoRepository,
				resourcemanager.WithIDGen(id.GenerateID),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			demoRepository := manager.Handler().(*demo.Repository)
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				demoRepository.Create(context.TODO(), sampleDemo)
			case rmtests.OperationListSortSuccess:
				demoRepository.Create(context.TODO(), sampleDemo)
				demoRepository.Create(context.TODO(), secondSampleDemo)
				demoRepository.Create(context.TODO(), thirdSampleDemo)
			}
		},
		SampleJSON: `{
			"type": "Demo",
			"spec": {
				"id": "1",
				"name": "dev",
				"enabled": true,
				"type": "pokeshop",
				"pokeshop": {
					"httpEndpoint": "http://dev-endpoint:1234",
					"grpcEndpoint": "dev-grpc:9091",
					"kafkaBroker": "dev-kafka:9092"
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Demo",
			"spec": {
				"id": "1",
				"name": "new-dev",
				"enabled": true,
				"type": "pokeshop",
				"pokeshop": {
					"httpEndpoint": "http://new-dev-endpoint:1234",
					"grpcEndpoint": "new-dev-grpc:9091",
					"kafkaBroker": "new-dev-kafka:9092"
				}
			}
		}`,
	})
}

func TestOpenTelemetryStoreDemoResource(t *testing.T) {
	sampleDemo := demo.Demo{
		ID:      "1",
		Name:    "dev",
		Type:    demo.DemoTypeOpentelemetryStore,
		Enabled: true,
		OpenTelemetryStore: &demo.OpenTelemetryStoreDemo{
			FrontendEndpoint:       "http://dev-frontend:9000",
			ProductCatalogEndpoint: "http://dev-product:8081",
			CartEndpoint:           "http://dev-cart:8082",
			CheckoutEndpoint:       "http://dev-checkout:8083",
		},
	}

	secondSampleDemo := demo.Demo{
		ID:      "2",
		Name:    "staging",
		Type:    demo.DemoTypePokeshop,
		Enabled: true,
		OpenTelemetryStore: &demo.OpenTelemetryStoreDemo{
			FrontendEndpoint:       "http://stg-frontend:9000",
			ProductCatalogEndpoint: "http://stg-product:8081",
			CartEndpoint:           "http://stg-cart:8082",
			CheckoutEndpoint:       "http://stg-checkout:8083",
		},
	}

	thirdSampleDemo := demo.Demo{
		ID:      "3",
		Name:    "production",
		Type:    demo.DemoTypePokeshop,
		Enabled: true,
		OpenTelemetryStore: &demo.OpenTelemetryStoreDemo{
			FrontendEndpoint:       "http://prod-frontend:9000",
			ProductCatalogEndpoint: "http://prod-product:8081",
			CartEndpoint:           "http://prod-cart:8082",
			CheckoutEndpoint:       "http://prod-checkout:8083",
		},
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: demo.ResourceName,
		ResourceTypePlural:   demo.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			demoRepository := demo.NewRepository(db)

			manager := resourcemanager.New[demo.Demo](
				demo.ResourceName,
				demo.ResourceNamePlural,
				demoRepository,
				resourcemanager.WithIDGen(id.GenerateID),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			demoRepository := manager.Handler().(*demo.Repository)
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				demoRepository.Create(context.TODO(), sampleDemo)
			case rmtests.OperationListSortSuccess:
				demoRepository.Create(context.TODO(), sampleDemo)
				demoRepository.Create(context.TODO(), secondSampleDemo)
				demoRepository.Create(context.TODO(), thirdSampleDemo)
			}
		},
		SampleJSON: `{
			"type": "Demo",
			"spec": {
				"id": "1",
				"name": "dev",
				"enabled": true,
				"type": "otelstore",
				"opentelemetryStore": {
					"frontendEndpoint": "http://dev-frontend:9000",
					"productCatalogEndpoint": "http://dev-product:8081",
					"cartEndpoint": "http://dev-cart:8082",
					"checkoutEndpoint": "http://dev-checkout:8083"
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Demo",
			"spec": {
				"id": "1",
				"name": "new-dev",
				"enabled": true,
				"type": "otelstore",
				"opentelemetryStore": {
					"frontendEndpoint": "http://dev-frontend",
					"productCatalogEndpoint": "http://dev-product",
					"cartEndpoint": "http://dev-cart",
					"checkoutEndpoint": "http://dev-checkout"
				}
			}
		}`,
	})
}
