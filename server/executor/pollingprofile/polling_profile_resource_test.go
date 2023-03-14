package pollingprofile_test

import (
	"context"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testmock"
)

func TestPollingProfileResource(t *testing.T) {
	db := testmock.MustGetRawTestingDatabase()
	sampleProfile := pollingprofile.PollingProfile{
		ID:       "1",
		Name:     "test",
		Default:  true,
		Strategy: pollingprofile.Periodic,
		Periodic: &pollingprofile.PeriodicPollingConfig{
			RetryDelay: "10s",
			Timeout:    "30m",
		},
	}

	secondSampleProfile := pollingprofile.PollingProfile{
		ID:       "2",
		Name:     "fast test",
		Default:  false,
		Strategy: pollingprofile.Periodic,
		Periodic: &pollingprofile.PeriodicPollingConfig{
			RetryDelay: "1s",
			Timeout:    "1m",
		},
	}

	thirdSampleProfile := pollingprofile.PollingProfile{
		ID:       "3",
		Name:     "long running test",
		Default:  false,
		Strategy: pollingprofile.Periodic,
		Periodic: &pollingprofile.PeriodicPollingConfig{
			RetryDelay: "2m",
			Timeout:    "45m",
		},
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "PollingProfile",
		RegisterManagerFn: func(router *mux.Router) resourcemanager.Manager {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			pollingProfileRepo := pollingprofile.NewRepository(db)

			manager := resourcemanager.New[pollingprofile.PollingProfile]("PollingProfile", pollingProfileRepo, resourcemanager.WithIDGen(id.GenerateID))
			manager.RegisterRoutes(router)

			return manager
		},
		Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
			pollingProfileRepo := manager.Handler().(resourcemanager.Create[pollingprofile.PollingProfile])
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				pollingProfileRepo.Create(context.TODO(), sampleProfile)
			case rmtests.OperationListPaginatedSuccess:
				pollingProfileRepo.Create(context.TODO(), sampleProfile)
				pollingProfileRepo.Create(context.TODO(), secondSampleProfile)
				pollingProfileRepo.Create(context.TODO(), thirdSampleProfile)
			}
		},
		SampleJSON: `{
			"type": "PollingProfile",
			"spec": {
				"id": "1",
				"name": "test",
				"default": true,
				"strategy": "periodic",
				"periodic": {
					"retryDelay": "10s",
					"timeout": "30m"
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "PollingProfile",
			"spec": {
				"id": "1",
				"name": "long test",
				"default": true,
				"strategy": "periodic",
				"periodic": {
					"retryDelay": "25s",
					"timeout": "1h"
				}
			}
		}`,
	},
		// TODO: remove this when we support multiple profiles
		rmtests.ExcludeOperations(
			rmtests.OperationGetNotFound,
			rmtests.OperationUpdateNotFound,
		))
}
