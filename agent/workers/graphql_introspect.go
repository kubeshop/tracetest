package workers

import (
	"context"
	"strings"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"github.com/kubeshop/tracetest/agent/workers/trigger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type GraphqlIntrospectWorker struct {
	client         *client.Client
	logger         *zap.Logger
	tracer         trace.Tracer
	meter          metric.Meter
	graphqlTrigger trigger.Triggerer
}

type GraphqlIntrospectOption func(*GraphqlIntrospectWorker)

func WithGraphqlIntrospectLogger(logger *zap.Logger) GraphqlIntrospectOption {
	return func(w *GraphqlIntrospectWorker) {
		w.logger = logger
	}
}

func WithGraphqlIntrospectTrigger(trigger trigger.Triggerer) GraphqlIntrospectOption {
	return func(w *GraphqlIntrospectWorker) {
		w.graphqlTrigger = trigger
	}
}

func WithGraphqlIntrospectTracer(tracer trace.Tracer) GraphqlIntrospectOption {
	return func(w *GraphqlIntrospectWorker) {
		w.tracer = tracer
	}
}

func WithGraphqlIntrospectMeter(meter metric.Meter) GraphqlIntrospectOption {
	return func(w *GraphqlIntrospectWorker) {
		w.meter = meter
	}
}

func NewGraphqlIntrospectWorker(client *client.Client, opts ...GraphqlIntrospectOption) *GraphqlIntrospectWorker {
	worker := &GraphqlIntrospectWorker{
		client: client,
		tracer: telemetry.GetNoopTracer(),
		logger: zap.NewNop(),
		meter:  telemetry.GetNoopMeter(),
	}

	http := trigger.HTTP()
	worker.graphqlTrigger = trigger.GRAPHQL(http)

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

func (w *GraphqlIntrospectWorker) Introspect(ctx context.Context, r *proto.GraphqlIntrospectRequest) error {
	ctx, span := w.tracer.Start(ctx, "TestConnectionRequest Worker operation")
	defer span.End()

	request := mapProtoToGraphqlRequest(r)

	response, err := w.graphqlTrigger.Trigger(ctx, request, nil)
	if err != nil {
		w.logger.Error("Could not send graphql introspection request", zap.Error(err))
		span.RecordError(err)
		return err
	}

	w.logger.Debug("Sending graphql introspection result", zap.Any("response", response))
	err = w.client.SendGraphqlIntrospectionResult(ctx, mapGraphqlToProtoResponse(r, response.Result.Graphql))
	if err != nil {
		w.logger.Error("Could not send graphql introspection result", zap.Error(err))
		span.RecordError(err)
	} else {
		w.logger.Debug("Sent graphql introspection result")
	}

	return nil
}

func mapProtoToGraphqlRequest(r *proto.GraphqlIntrospectRequest) trigger.Trigger {
	headers := make([]trigger.HTTPHeader, 0)
	for _, header := range r.Graphql.Headers {
		headers = append(headers, trigger.HTTPHeader{
			Key:   header.Key,
			Value: header.Value,
		})
	}

	headers = append(headers, trigger.HTTPHeader{
		Value: "application/json",
		Key:   "Content-Type",
	})

	request := &trigger.GraphqlRequest{
		URL:             r.Graphql.Url,
		SSLVerification: r.Graphql.SSLVerification,
		Headers:         headers,
	}

	request.Body = trigger.GraphqlBody{
		Query: IntrospectionQuery,
	}

	return trigger.Trigger{
		Type:    trigger.TriggerTypeGraphql,
		Graphql: request,
	}
}

func mapGraphqlToProtoResponse(r *proto.GraphqlIntrospectRequest, resp *trigger.GraphqlResponse) *proto.GraphqlIntrospectResponse {
	headers := make([]*proto.HttpHeader, 0)
	for _, header := range resp.Headers {
		headers = append(headers, &proto.HttpHeader{
			Key:   header.Key,
			Value: header.Value,
		})
	}

	return &proto.GraphqlIntrospectResponse{
		RequestID: r.RequestID,
		Response: &proto.HttpResponse{
			StatusCode: int32(resp.StatusCode),
			Status:     resp.Status,
			Headers:    headers,
			Body:       []byte(resp.Body),
		},
	}
}

var IntrospectionQuery = strings.ReplaceAll(`
  query IntrospectionQuery {
    __schema {
      queryType { name }
      mutationType { name }
      subscriptionType { name }
      types {
        ...FullType
      }
      directives {
        name
        description
		locations
        args {
          ...InputValue
        }
      }
    }
  }

  fragment FullType on __Type {
    kind
    name
    description
    fields(includeDeprecated: true) {
      name
      description
      args {
        ...InputValue
      }
      type {
        ...TypeRef
      }
      isDeprecated
      deprecationReason
    }
    inputFields {
      ...InputValue
    }
    interfaces {
      ...TypeRef
    }
    enumValues(includeDeprecated: true) {
      name
      description
      isDeprecated
      deprecationReason
    }
    possibleTypes {
      ...TypeRef
    }
  }

  fragment InputValue on __InputValue {
    name
    description
    type { ...TypeRef }
    defaultValue
  }

  fragment TypeRef on __Type {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                }
              }
            }
          }
        }
      }
    }
  }
`, "\n", "")
