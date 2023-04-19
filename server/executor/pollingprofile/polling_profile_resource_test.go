package pollingprofile_test

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	rmtests "github.com/kubeshop/tracetest/server/resourcemanager/testutil"
	"github.com/kubeshop/tracetest/server/testmock"
)

func TestPollingProfileResource(t *testing.T) {
	db := testmock.MustGetRawTestingDatabase()
	// sampleProfile := pollingprofile.PollingProfile{
	// 	ID:       "1",
	// 	Name:     "test",
	// 	Default:  true,
	// 	Strategy: pollingprofile.Periodic,
	// 	Periodic: &pollingprofile.PeriodicPollingConfig{
	// 		RetryDelay: "10s",
	// 		Timeout:    "30m",
	// 	},
	// }

	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: pollingprofile.ResourceName,
		ResourceTypePlural:   pollingprofile.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router) resourcemanager.Manager {
			db := testmock.MustCreateRandomMigratedDatabase(db)
			pollingProfileRepo := pollingprofile.NewRepository(db)

			manager := resourcemanager.New[pollingprofile.PollingProfile](
				pollingprofile.ResourceName,
				pollingprofile.ResourceNamePlural,
				pollingProfileRepo,
				resourcemanager.WithOperations(pollingprofile.Operations...),
				resourcemanager.WithIDGen(id.GenerateID),
			)
			manager.RegisterRoutes(router)

			return manager
		},
		// Prepare: func(t *testing.T, op rmtests.Operation, manager resourcemanager.Manager) {
		// 	pollingProfileRepo := manager.Handler().(*pollingprofile.Repository)
		// 	switch op {
		// 	case rmtests.OperationGetSuccess,
		// 		pollingProfileRepo.Update(context.TODO(), sampleProfile)
		// 	case rmtests.OperationListPaginatedSuccess:
		// 		pollingProfileRepo.Create(context.TODO(), sampleProfile)
		// 		pollingProfileRepo.Create(context.TODO(), secondSampleProfile)
		// 		pollingProfileRepo.Create(context.TODO(), thirdSampleProfile)
		// 	}
		// },
		SampleJSON: `{
			"type": "PollingProfile",
			"spec": {
				"id": "current",
				"name": "default",
				"default": true,
				"strategy": "periodic",
				"periodic": {
					"timeout": "1m",
					"retryDelay": "5s"
				}
			}
		}`,
		SampleJSONUpdated: `{
			"type": "PollingProfile",
			"spec": {
				"id": "current",
				"name": "long test",
				"default": true,
				"strategy": "periodic",
				"periodic": {
					"timeout": "1h",
					"retryDelay": "25s"
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
