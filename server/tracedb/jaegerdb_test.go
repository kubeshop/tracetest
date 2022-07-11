package tracedb_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestJaegerGetTraceByID(t *testing.T) {
	t.Skip("TODO: docker-compose jaeger")
	db, err := tracedb.New(config.Config{
		TracingBackend: config.TracingBackend{
			DataStore: config.TracingBackendDataStoreConfig{
				Type: tracedb.JAEGER_BACKEND,
				Jaeger: configgrpc.GRPCClientSettings{
					Endpoint:   "localhost:16685",
					TLSSetting: configtls.TLSClientSetting{Insecure: true},
				},
			},
		},
	})
	assert.NoError(t, err)

	defer db.Close()
	trace, err := db.GetTraceByID(context.Background(), "0194fdc2fa2ffcc041d3ff12045b73c8")
	assert.NoError(t, err)

	buf := bytes.Buffer{}
	m := jsonpb.Marshaler{}
	err = m.Marshal(&buf, trace)
	assert.NoError(t, err)
	fmt.Printf("\n%s\n", buf.String())
}
