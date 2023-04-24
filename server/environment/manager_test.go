package environment_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/stretchr/testify/require"
)

func TestEnvironment(t *testing.T) {
	sampleEnvironment := environment.Environment{
		ID:          "dev",
		Name:        "dev",
		Description: "dev variables",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Values: []environment.EnvironmentValue{
			{Key: "URL", Value: "http://dev-app.com"},
		},
	}

	secondEnvironment := environment.Environment{
		ID:          "new-dev",
		Name:        "new dev",
		Description: "new dev variables",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Values: []environment.EnvironmentValue{
			{Key: "URL", Value: "http://dev-app.com"},
			{Key: "AUTH", Value: "user:pass"},
		},
	}

	thirdEnvironment := environment.Environment{
		ID:          "stg",
		Name:        "staging",
		Description: "staging variables",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Values: []environment.EnvironmentValue{
			{Key: "URL", Value: "http://stg-app.com"},
			{Key: "AUTH", Value: "user:pass"},
		},
	}

	resourceTypeTest := rmtests.ResourceTypeTest{
		ResourceTypeSingular: environment.ResourceName,
		ResourceTypePlural:   environment.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			environmentRepository := environment.NewRepository(db)

			manager := resourcemanager.New[environment.Environment](
				environment.ResourceName,
				environment.ResourceNamePlural,
				environmentRepository,
				resourcemanager.WithIDGen(id.GenerateID),
				resourcemanager.WithOperations(environment.Operations...),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			environmentRepository := manager.Handler().(environment.Repository)
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				environmentRepository.Create(context.TODO(), sampleEnvironment)
			case rmtests.OperationListPaginatedSuccess:
				environmentRepository.Create(context.TODO(), sampleEnvironment)
				environmentRepository.Create(context.TODO(), secondEnvironment)
				environmentRepository.Create(context.TODO(), thirdEnvironment)
			}
		},
		SampleJSON: `{
			"type": "Environment",
			"spec": {
				"id": "dev",
				"name": "dev",
				"description": "dev variables",
				"values": [
					{ "key": "URL", "value": "http://dev-app.com" }
				]
			}
		}`,
		SampleJSONUpdated: `{
			"type": "Environment",
			"spec": {
				"id": "dev",
				"name": "new-dev",
				"description": "new dev variables",
				"values": [
					{ "key": "URL", "value": "http://dev-app.com" },
					{ "key": "AUTH", "value": "user:password" }
				]
			}
		}`,
	}

	rmtests.TestResourceType(t, resourceTypeTest, rmtests.JSONComparer(environmentJSONComparer))
}

func environmentJSONComparer(t require.TestingT, operation rmtests.Operation, firstValue, secondValue string) {
	expected := rmtests.RemoveFieldFromJSONResource("createdAt", firstValue)
	obtained := rmtests.RemoveFieldFromJSONResource("createdAt", secondValue)

	require.JSONEq(t, expected, obtained)
}
