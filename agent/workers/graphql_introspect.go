package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/astprinter"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/introspection"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type GraphqlIntrospectWorker struct {
	client *client.Client
	logger *zap.Logger
	tracer trace.Tracer
	meter  metric.Meter
}

type GraphqlIntrospectOption func(*GraphqlIntrospectWorker)

func WithGraphqlIntrospectLogger(logger *zap.Logger) GraphqlIntrospectOption {
	return func(w *GraphqlIntrospectWorker) {
		w.logger = logger
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

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

type transportHeaders struct {
	headers http.Header
	wrapped http.RoundTripper
}

func (t *transportHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range t.headers {
		req.Header[k] = v
	}

	return t.wrapped.RoundTrip(req)
}

func (w *GraphqlIntrospectWorker) Introspect(ctx context.Context, request *proto.GraphqlIntrospectRequest) error {
	ctx, span := w.tracer.Start(ctx, "TestConnectionRequest Worker operation")
	defer span.End()

	headers := make(http.Header)
	for _, header := range request.Headers {
		headers.Add(header.Key, header.Value)
	}

	httpClient := http.Client{
		Transport: &transportHeaders{
			headers: headers,
			wrapped: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient(request.Url, &httpClient)

	graphqlRequest := &graphql.Request{
		Query: IntrospectionQuery,
	}

	var graphqlResponse graphql.Response
	err := graphqlClient.MakeRequest(ctx, graphqlRequest, &graphqlResponse)
	if err != nil {
		w.logger.Error("Could not make graphql request", zap.Error(err))
		span.RecordError(err)
		return err
	}

	json, err := json.Marshal(graphqlResponse.Data)
	if err != nil {
		w.logger.Error("Could not marshal graphql response", zap.Error(err))
		span.RecordError(err)
		return err
	}

	converter := introspection.JsonConverter{}
	buf := bytes.NewBuffer(json)

	doc, err := converter.GraphQLDocument(buf)
	if err != nil {
		w.logger.Error("Could not convert graphql document", zap.Error(err))
		span.RecordError(err)
		return err
	}

	outWriter := &bytes.Buffer{}
	err = astprinter.PrintIndent(doc, []byte("  "), outWriter)
	if err != nil {
		w.logger.Error("Could not print graphql document", zap.Error(err))
		span.RecordError(err)
		return err
	}

	response := &proto.GraphqlIntrospectResponse{
		RequestID:  request.RequestID,
		Successful: true,
		Schema:     outWriter.String(),
	}

	w.logger.Debug("Sending datastore connection test result", zap.Any("response", response))
	err = w.client.SendGraphqlIntrospectionResult(ctx, response)
	if err != nil {
		w.logger.Error("Could not send datastore connection test result", zap.Error(err))
		span.RecordError(err)
	} else {
		w.logger.Debug("Sent datastore connection test result")
	}

	return nil
}

var IntrospectionQuery = `
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
`
