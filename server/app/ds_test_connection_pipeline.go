package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/testconnection"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func buildDataStoreTestPipeline(
	pool *pgxpool.Pool,
	dsTestListener *testconnection.Listener,
	tracer trace.Tracer,
	newTraceDBFn tracedb.FactoryFunc,
	appConfig *config.AppConfig,
	meter metric.Meter,
) *testconnection.DataStoreTestPipeline {
	requestWorker := testconnection.NewDsTestConnectionRequest(tracer, newTraceDBFn, appConfig.DataStorePipelineTestConnectionEnabled())
	notifyWorker := testconnection.NewDsTestConnectionNotify(dsTestListener, tracer)

	pgQueue := pipeline.NewPostgresQueueDriver[testconnection.Job](pool, pgChannelName)

	pipeline := pipeline.New(testconnection.NewConfigurer(meter),
		pipeline.Step[testconnection.Job]{Processor: requestWorker, Driver: pgQueue.Channel("datastore_test_connection_request")},
		pipeline.Step[testconnection.Job]{Processor: notifyWorker, Driver: pgQueue.Channel("datastore_test_connection_notify")},
	)

	return testconnection.NewDataStoreTestPipeline(pipeline, dsTestListener)
}
