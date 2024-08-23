package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/stretchr/testify/require"
)

func TestGraphqlIntrospectionWorker(t *testing.T) {
	ctx := ContextWithTracingEnabled()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	worker := workers.NewGraphqlIntrospectWorker(client)

	client.OnGraphqlIntrospectionRequest(func(ctx context.Context, pr *proto.GraphqlIntrospectRequest) error {
		return worker.Introspect(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	request := &proto.GraphqlIntrospectRequest{
		RequestID: "test",
		Url:       "https://swapi-graphql.netlify.app/.netlify/functions/index",
	}

	controlPlane.SendGraphqlIntrospectionRequest(ctx, request)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastGraphqlIntrospectionResponse()
	require.NotNil(t, resp, "agent did not send graphql introspection response back to server")
}
