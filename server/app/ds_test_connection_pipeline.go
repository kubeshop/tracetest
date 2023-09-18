package app

import (
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/testconnection"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

func buildDataStoreTestPipeline(
	dsTestListener *testconnection.Listener,
	tracer trace.Tracer,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
) *testconnection.DataStoreTestPipeline {
	requestWorker := testconnection.NewDsTestConnectionRequest(dsTestListener, tracer, newTraceDBFn, dsRepo)

	pipeline := pipeline.New(&testconnection.Configurer[testconnection.Job]{},
		pipeline.Step[testconnection.Job]{Processor: requestWorker, Driver: pipeline.NewInMemoryQueueDriver[testconnection.Job]("datastore_test_connection")},
	)

	return testconnection.NewDataStoreTestPipeline(pipeline, dsTestListener)
}
