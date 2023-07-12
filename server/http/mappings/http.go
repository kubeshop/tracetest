package mappings

import (
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/test/trigger"
)

// out

func (m OpenAPI) HTTPHeaders(in []trigger.HTTPHeader) []openapi.HttpHeader {
	headers := make([]openapi.HttpHeader, len(in))
	for i, h := range in {
		headers[i] = openapi.HttpHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m OpenAPI) HTTPRequest(in *trigger.HTTPRequest) openapi.HttpRequest {
	if in == nil {
		return openapi.HttpRequest{}
	}

	return openapi.HttpRequest{
		Url:             in.URL,
		Method:          string(in.Method),
		Headers:         m.HTTPHeaders(in.Headers),
		Body:            in.Body,
		Auth:            m.Auth(in.Auth),
		SslVerification: in.SSLVerification,
	}
}

func (m OpenAPI) HTTPResponse(in *trigger.HTTPResponse) openapi.HttpResponse {
	if in == nil {
		return openapi.HttpResponse{}
	}
	return openapi.HttpResponse{
		Status:     in.Status,
		StatusCode: int32(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m OpenAPI) Auth(in *trigger.HTTPAuthenticator) openapi.HttpAuth {
	if in == nil {
		return openapi.HttpAuth{}
	}

	auth := openapi.HttpAuth{
		Type: in.Type,
	}
	switch in.Type {
	case "apiKey":
		auth.ApiKey = openapi.HttpAuthApiKey{
			Key:   in.APIKey.Key,
			Value: in.APIKey.Value,
			In:    string(in.APIKey.In),
		}
	case "basic":
		auth.Basic = openapi.HttpAuthBasic{
			Username: in.Basic.Username,
			Password: in.Basic.Password,
		}
	case "bearer":
		auth.Bearer = openapi.HttpAuthBearer{
			Token: in.Bearer.Bearer,
		}
	}

	return auth
}

// in

func (m Model) HTTPHeaders(in []openapi.HttpHeader) []trigger.HTTPHeader {
	headers := make([]trigger.HTTPHeader, len(in))
	for i, h := range in {
		headers[i] = trigger.HTTPHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m Model) HTTPRequest(in openapi.HttpRequest) *trigger.HTTPRequest {
	// ignore unset http requests
	if in.Url == "" {
		return nil
	}

	return &trigger.HTTPRequest{
		URL:             in.Url,
		Method:          trigger.HTTPMethod(in.Method),
		Headers:         m.HTTPHeaders(in.Headers),
		Body:            in.Body,
		Auth:            m.Auth(in.Auth),
		SSLVerification: in.SslVerification,
	}
}

func (m Model) HTTPResponse(in openapi.HttpResponse) *trigger.HTTPResponse {
	// ignore unset http responses
	if in.StatusCode == 0 {
		return nil
	}

	return &trigger.HTTPResponse{
		Status:     in.Status,
		StatusCode: int(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m Model) Auth(in openapi.HttpAuth) *trigger.HTTPAuthenticator {
	return &trigger.HTTPAuthenticator{
		Type: in.Type,
		APIKey: &trigger.APIKeyAuthenticator{
			Key:   in.ApiKey.Key,
			Value: in.ApiKey.Value,
			In:    trigger.APIKeyPosition(in.ApiKey.In),
		},
		Basic: &trigger.BasicAuthenticator{
			Username: in.Basic.Username,
			Password: in.Basic.Password,
		},
		Bearer: &trigger.BearerAuthenticator{
			Bearer: in.Bearer.Token,
		},
	}
}
