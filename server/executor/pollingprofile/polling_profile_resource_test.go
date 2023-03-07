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
	sampleConfig := pollingprofile.PollingProfile{
		ID:       "1",
		Name:     "test",
		Strategy: pollingprofile.Periodic,
		Periodic: &pollingprofile.PeriodicPollingConfig{
			RetryDelay: "10s",
			Timeout:    "30m",
		},
	}

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceType: "PollingProfile",
		RegisterManagerFn: func(router *mux.Router) any {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			pollingProfileRepo := pollingprofile.Repository(db)

			manager := resourcemanager.New[pollingprofile.PollingProfile]("PollingProfile", pollingProfileRepo, id.GenerateID)
			manager.RegisterRoutes(router)

			return pollingProfileRepo
		},
		Prepare: func(t *testing.T, op rmtests.Operation, bridge any) {
			pollingProfileRepo := bridge.(resourcemanager.ResourceHandler[pollingprofile.PollingProfile])
			switch op {
			case rmtests.OperationGetSuccess,
				rmtests.OperationUpdateSuccess,
				rmtests.OperationDeleteSuccess,
				rmtests.OperationListSuccess:
				pollingProfileRepo.Create(context.TODO(), sampleConfig)
			}
		},
		SampleJSON: `{
			"type": "PollingProfile",
			"spec": {
				"id": "1",
				"name": "test",
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
				"strategy": "periodic",
				"periodic": {
					"retryDelay": "25s",
					"timeout": "1h"
				}
			}
		}`,
	})
}
