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
	rmtests.TestResourceType(t, rmtests.ResourceTypeTest{
		ResourceTypeSingular: pollingprofile.ResourceName,
		ResourceTypePlural:   pollingprofile.ResourceNamePlural,
		RegisterManagerFn: func(router *mux.Router) resourcemanager.Manager {
			db := testmock.CreateMigratedDatabase()
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
		rmtests.ExcludeOperations(
			rmtests.OperationGetNotFound,
			rmtests.OperationUpdateNotFound,
		))
}
