package definition

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

func validateRequest(r model.HTTPRequest) error {
	if r.URL == "" {
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

func validateHeader(h model.HTTPHeader) error {
	if h.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	return nil
}

func validateAuthentication(a *model.HTTPAuthenticator) error {
	switch a.Type {
	case "basic":
		if err := validateBasicAuth(a); err != nil {
			return fmt.Errorf("basic authentication must be valid: %w", err)
		}
		return nil
	case "bearer":
		if err := validateBearerAuth(a); err != nil {
			return fmt.Errorf("bearer authentication must be valid: %w", err)
		}
		return nil
	case "apiKey":
		if err := validateApiKeyAuth(a); err != nil {
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

func validateBasicAuth(a *model.HTTPAuthenticator) error {
	if a.Props["username"] == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if a.Props["password"] == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

func validateBearerAuth(a *model.HTTPAuthenticator) error {
	if a.Props["token"] == "" {
		return fmt.Errorf("token cannot be empty")
	}

	return nil
}

func validateApiKeyAuth(a *model.HTTPAuthenticator) error {
	if a.Props["key"] == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if a.Props["value"] == "" {
		return fmt.Errorf("value cannot be empty")
	}

	supportedInOptions := map[string]bool{
		"query":  true,
		"header": true,
	}

	in := a.Props["in"]

	if _, ok := supportedInOptions[in]; !ok {
		return fmt.Errorf("in option \"%s\" is not supported. Only \"query\" and \"header\" are supported", in)
	}

	return nil
}
