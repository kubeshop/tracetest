package demoresource_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/config/demoresource"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
)

func TestPokeshopDemoResource(t *testing.T) {
	sampleDemo := demoresource.Demo{
		ID:      "1",
		Name:    "dev",
		Type:    demoresource.DemoTypePokeshop,
		Enabled: true,
		Pokeshop: &demoresource.PokeshopDemo{
			HTTPEndpoint: "http://dev-endpoint:1234",
			GRPCEndpoint: "dev-grpc:9091",
		},
	}

	secondSampleDemo := demoresource.Demo{
		ID:      "2",
		Name:    "staging",
		Type:    demoresource.DemoTypePokeshop,
		Enabled: true,
		Pokeshop: &demoresource.PokeshopDemo{
			HTTPEndpoint: "http://stg-endpoint:1234",
			GRPCEndpoint: "stg-grpc:9091",
		},
	}

	thirdSampleDemo := demoresource.Demo{
		ID:      "3",
		Name:    "production",
		Type:    demoresource.DemoTypePokeshop,
		Enabled: true,
		Pokeshop: &demoresource.PokeshopDemo{
			HTTPEndpoint: "http://prod-endpoint:1234",
			GRPCEndpoint: "prod-grpc:9091",
		},
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: demoresource.ResourceName,
		ResourceTypePlural:   demoresource.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			demoRepository := demoresource.NewRepository(db)

			manager := resourcemanager.New[demoresource.Demo](
				demoresource.ResourceName,
				demoresource.ResourceNamePlural,
				demoRepository,
				resourcemanager.WithIDGen(id.GenerateID),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			demoRepository := manager.Handler().(*demoresource.Repository)
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
					"grpcEndpoint": "dev-grpc:9091"
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
					"grpcEndpoint": "new-dev-grpc:9091"
				}
			}
		}`,
	})
}

func TestOpenTelemetryStoreDemoResource(t *testing.T) {
	sampleDemo := demoresource.Demo{
		ID:      "1",
		Name:    "dev",
		Type:    demoresource.DemoTypeOpentelemetryStore,
		Enabled: true,
		OpenTelemetryStore: &demoresource.OpenTelemetryStoreDemo{
			FrontendEndpoint:       "http://dev-frontend:9000",
			ProductCatalogEndpoint: "http://dev-product:8081",
			CartEndpoint:           "http://dev-cart:8082",
			CheckoutEndpoint:       "http://dev-checkout:8083",
		},
	}

	secondSampleDemo := demoresource.Demo{
		ID:      "2",
		Name:    "staging",
		Type:    demoresource.DemoTypePokeshop,
		Enabled: true,
		OpenTelemetryStore: &demoresource.OpenTelemetryStoreDemo{
			FrontendEndpoint:       "http://stg-frontend:9000",
			ProductCatalogEndpoint: "http://stg-product:8081",
			CartEndpoint:           "http://stg-cart:8082",
			CheckoutEndpoint:       "http://stg-checkout:8083",
		},
	}

	thirdSampleDemo := demoresource.Demo{
		ID:      "3",
		Name:    "production",
		Type:    demoresource.DemoTypePokeshop,
		Enabled: true,
		OpenTelemetryStore: &demoresource.OpenTelemetryStoreDemo{
			FrontendEndpoint:       "http://prod-frontend:9000",
			ProductCatalogEndpoint: "http://prod-product:8081",
			CartEndpoint:           "http://prod-cart:8082",
			CheckoutEndpoint:       "http://prod-checkout:8083",
		},
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: demoresource.ResourceName,
		ResourceTypePlural:   demoresource.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			demoRepository := demoresource.NewRepository(db)

			manager := resourcemanager.New[demoresource.Demo](
				demoresource.ResourceName,
				demoresource.ResourceNamePlural,
				demoRepository,
				resourcemanager.WithIDGen(id.GenerateID),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			demoRepository := manager.Handler().(*demoresource.Repository)
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
