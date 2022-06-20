package definition

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/openapi"
)

func validateRequest(r openapi.HttpRequest) error {
	if r.Url == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if r.Method == "" {
		return fmt.Errorf("method cannot be empty")
	}

	for _, header := range r.Headers {
		if err := validateHeader(header); err != nil {
			return fmt.Errorf("http header must be valid: %w", err)
		}
	}

	if err := validateAuthentication(r.Auth); err != nil {
		return fmt.Errorf("http authentication must be valid: %w", err)
	}

	return nil
}

func validateHeader(h openapi.HttpHeader) error {
	if h.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	return nil
}

func validateAuthentication(a openapi.HttpAuth) error {
	switch a.Type {
	case "basic":
		if err := validateBasicAuth(a.Basic); err != nil {
			return fmt.Errorf("basic authentication must be valid: %w", err)
		}
		return nil
	case "bearer":
		if err := validateBearerAuth(a.Bearer); err != nil {
			return fmt.Errorf("bearer authentication must be valid: %w", err)
		}
		return nil
	case "apiKey":
		if err := validateApiKeyAuth(a.ApiKey); err != nil {
			return fmt.Errorf("apiKey authentication must be valid: %w", err)
		}
		return nil
	case "":
		// No authentication
		return nil
	default:
		// Any other type that is not supported
		return fmt.Errorf("type \"%s\" is not supported", a.Type)
	}
}

func validateBasicAuth(a openapi.HttpAuthBasic) error {
	if a.Username == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if a.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

func validateBearerAuth(a openapi.HttpAuthBearer) error {
	if a.Token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	return nil
}

func validateApiKeyAuth(a openapi.HttpAuthApiKey) error {
	if a.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if a.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	supportedInOptions := map[string]bool{
		"query":  true,
		"header": true,
	}

	if _, ok := supportedInOptions[a.In]; !ok {
		return fmt.Errorf("in option \"%s\" is not supported. Only \"query\" and \"header\" are supported", a.In)
	}

	return nil
}
