package variableset_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/variableset"
	"github.com/stretchr/testify/require"
)

func TestVariableSetRepository(t *testing.T) {
	sampleVariableSet := variableset.VariableSet{
		ID:          "dev",
		Name:        "dev",
		Description: "dev variables",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Values: []variableset.VariableSetValue{
			{Key: "URL", Value: "http://dev-app.com"},
		},
	}

	secondVariableSet := variableset.VariableSet{
		ID:          "new-dev",
		Name:        "new dev",
		Description: "new dev variables",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Values: []variableset.VariableSetValue{
			{Key: "URL", Value: "http://dev-app.com"},
			{Key: "AUTH", Value: "user:pass"},
		},
	}

	thirdVariableSet := variableset.VariableSet{
		ID:          "stg",
		Name:        "staging",
		Description: "staging variables",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Values: []variableset.VariableSetValue{
			{Key: "URL", Value: "http://stg-app.com"},
			{Key: "AUTH", Value: "user:pass"},
		},
	}

	resourceTypeTest := rmtests.ResourceTypeTest{
		ResourceTypeSingular: variableset.ResourceName,
		ResourceTypePlural:   variableset.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router, db *sql.DB) resourcemanager.Manager {
			variableSetRepository := variableset.NewRepository(db)

			manager := resourcemanager.New[variableset.VariableSet](
				variableset.ResourceName,
				variableset.ResourceNamePlural,
				variableSetRepository,
				resourcemanager.WithIDGen(id.GenerateID),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			variableSetRepository := manager.Handler().(*variableset.Repository)
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				variableSetRepository.Create(context.TODO(), sampleVariableSet)
			case rmtests.OperationListSortSuccess:
				variableSetRepository.Create(context.TODO(), sampleVariableSet)
				variableSetRepository.Create(context.TODO(), secondVariableSet)
				variableSetRepository.Create(context.TODO(), thirdVariableSet)
			}
		},
		SampleJSON: `{
			"type": "VariableSet",
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
			"type": "VariableSet",
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

	rmtests.TestResourceType(t, resourceTypeTest, rmtests.JSONComparer(VariableSetJSONComparer))
}

func VariableSetJSONComparer(t require.TestingT, operation rmtests.Operation, firstValue, secondValue string) {
	expected := rmtests.RemoveFieldFromJSONResource("createdAt", firstValue)
	obtained := rmtests.RemoveFieldFromJSONResource("createdAt", secondValue)

	require.JSONEq(t, expected, obtained)
}
