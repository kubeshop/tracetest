package workers_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers"
	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGraphqlTrigger struct{}

func (m mockGraphqlTrigger) Trigger(ctx context.Context, triggerConfig trigger.Trigger, opts *trigger.Options) (trigger.Response, error) {
	return trigger.Response{
		Result: trigger.TriggerResult{
			Type: trigger.TriggerTypeGraphql,
			Graphql: &trigger.GraphqlResponse{
				StatusCode: 200,
			},
		},
	}, nil
}

func (m mockGraphqlTrigger) Type() trigger.TriggerType {
	return trigger.TriggerTypeGraphql
}

func TestGraphqlIntrospectionWorker(t *testing.T) {
	ctx := ContextWithTracingEnabled()
	controlPlane := mocks.NewGrpcServer()

	client, err := client.Connect(ctx, controlPlane.Addr(), client.WithInsecure())
	require.NoError(t, err)

	worker := workers.NewGraphqlIntrospectWorker(client, workers.WithGraphqlIntrospectTrigger(mockGraphqlTrigger{}))

	client.OnGraphqlIntrospectionRequest(func(ctx context.Context, pr *proto.GraphqlIntrospectRequest) error {
		return worker.Introspect(ctx, pr)
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	request := &proto.GraphqlIntrospectRequest{
		RequestID: "test",
		Graphql: &proto.GraphqlRequest{
			Url: "https://swapi-graphql.netlify.app/.netlify/functions/index",
			Headers: []*proto.HttpHeader{
				{Key: "Content-Type", Value: "application/json"},
			},
		},
	}

	controlPlane.SendGraphqlIntrospectionRequest(ctx, request)
	time.Sleep(1 * time.Second)

	resp := controlPlane.GetLastGraphqlIntrospectionResponse()
	require.NotNil(t, resp, "agent did not send graphql introspection response back to server")

	assert.Equal(t, resp.Data.Response.StatusCode, int32(200))

}
