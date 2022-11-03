package mappings

import (
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

// out

func (m OpenAPI) HTTPHeaders(in []model.HTTPHeader) []openapi.HttpHeader {
	headers := make([]openapi.HttpHeader, len(in))
	for i, h := range in {
		headers[i] = openapi.HttpHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m OpenAPI) HTTPRequest(in *model.HTTPRequest) openapi.HttpRequest {
	if in == nil {
		return openapi.HttpRequest{}
	}

	return openapi.HttpRequest{
		Url:     in.URL,
		Method:  string(in.Method),
		Headers: m.HTTPHeaders(in.Headers),
		Body:    in.Body,
		Auth:    m.Auth(in.Auth),
	}
}

func (m OpenAPI) HTTPResponse(in *model.HTTPResponse) openapi.HttpResponse {
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

func (m OpenAPI) Auth(in *model.HTTPAuthenticator) openapi.HttpAuth {
	if in == nil {
		return openapi.HttpAuth{}
	}

	auth := openapi.HttpAuth{
		Type: in.Type,
	}
	switch in.Type {
	case "apiKey":
		auth.ApiKey = openapi.HttpAuthApiKey{
			Key:   in.Props["key"],
			Value: in.Props["value"],
			In:    in.Props["in"],
		}
	case "basic":
		auth.Basic = openapi.HttpAuthBasic{
			Username: in.Props["username"],
			Password: in.Props["password"],
		}
	case "bearer":
		auth.Bearer = openapi.HttpAuthBearer{
			Token: in.Props["bearer"],
		}
	}

	return auth
}

// in

func (m Model) HTTPHeaders(in []openapi.HttpHeader) []model.HTTPHeader {
	headers := make([]model.HTTPHeader, len(in))
	for i, h := range in {
		headers[i] = model.HTTPHeader{Key: h.Key, Value: h.Value}
	}

	return headers
}

func (m Model) HTTPRequest(in openapi.HttpRequest) *model.HTTPRequest {
	// ignore unset http requests
	if in.Url == "" {
		return nil
	}

	return &model.HTTPRequest{
		URL:     in.Url,
		Method:  model.HTTPMethod(in.Method),
		Headers: m.HTTPHeaders(in.Headers),
		Body:    in.Body,
		Auth:    m.Auth(in.Auth),
	}
}

func (m Model) HTTPResponse(in openapi.HttpResponse) *model.HTTPResponse {
	// ignore unset http responses
	if in.StatusCode == 0 {
		return nil
	}

	return &model.HTTPResponse{
		Status:     in.Status,
		StatusCode: int(in.StatusCode),
		Headers:    m.HTTPHeaders(in.Headers),
		Body:       in.Body,
	}
}

func (m Model) Auth(in openapi.HttpAuth) *model.HTTPAuthenticator {
	var props map[string]string

	switch in.Type {
	case "apiKey":
		props = map[string]string{
			"key":   in.ApiKey.Key,
			"value": in.ApiKey.Value,
			"in":    in.ApiKey.In,
		}
	case "basic":
		props = map[string]string{
			"username": in.Basic.Username,
			"password": in.Basic.Password,
		}
	case "bearer":
		props = map[string]string{
			"token": in.Bearer.Token,
		}
	}

	return &model.HTTPAuthenticator{
		Type:  in.Type,
		Props: props,
	}
}
