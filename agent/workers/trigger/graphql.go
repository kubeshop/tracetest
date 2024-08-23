package trigger

import (
	"context"
	"fmt"
)

func GRAPHQL(httpTriggerer Triggerer) Triggerer {
	return &graphqlTriggerer{httpTriggerer}
}

type graphqlTriggerer struct {
	httpTriggerer Triggerer
}

func (te *graphqlTriggerer) Trigger(ctx context.Context, triggerConfig Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: TriggerResult{
			Type: te.Type(),
		},
	}

	if triggerConfig.Type != TriggerTypeGraphql {
		return response, fmt.Errorf(`trigger type "%s" not supported by HTTP triggerer`, triggerConfig.Type)
	}

	triggerConfig = mapGraphqlToHttp(triggerConfig)

	response, err := te.httpTriggerer.Trigger(ctx, triggerConfig, opts)
	if err != nil {
		return response, fmt.Errorf("error triggering Graphql trigger: %w", err)
	}

	return mapHttpToGraphql(response), nil
}

func (t *graphqlTriggerer) Type() TriggerType {
	return TriggerTypeGraphql
}

const TriggerTypeGraphql TriggerType = "graphql"

func mapGraphqlToHttp(triggerConfig Trigger) Trigger {
	return Trigger{
		Type: TriggerTypeHTTP,
		HTTP: &HTTPRequest{
			URL:             triggerConfig.Graphql.URL,
			Method:          HTTPMethodPOST,
			Body:            triggerConfig.Graphql.Body,
			Headers:         triggerConfig.Graphql.Headers,
			SSLVerification: triggerConfig.Graphql.SSLVerification,
		},
	}
}

func mapHttpToGraphql(response Response) Response {
	return Response{
		TraceID:        response.TraceID,
		SpanID:         response.SpanID,
		SpanAttributes: response.SpanAttributes,
		Result: TriggerResult{
			Type: TriggerTypeGraphql,
			Graphql: &GraphqlResponse{
				Status:     response.Result.HTTP.Status,
				StatusCode: response.Result.HTTP.StatusCode,
				Headers:    response.Result.HTTP.Headers,
				Body:       response.Result.HTTP.Body,
			},
		},
	}
}

type GraphqlRequest struct {
	URL             string             `expr_enabled:"true" json:"url,omitempty"`
	Body            string             `expr_enabled:"true" json:"body,omitempty"`
	Headers         []HTTPHeader       `json:"headers,omitempty"`
	Auth            *HTTPAuthenticator `json:"auth,omitempty"`
	SSLVerification bool               `json:"sslVerification,omitempty"`
}

type GraphqlResponse struct {
	Status     string
	StatusCode int
	Headers    []HTTPHeader
	Body       string
}
