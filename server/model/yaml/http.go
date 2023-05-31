package yaml

import (
	"fmt"
)

type HTTPRequest struct {
	URL             string              `yaml:"url"`
	Method          string              `yaml:"method"`
	Headers         []HTTPHeader        `yaml:"headers,omitempty" dc:"headers"`
	Authentication  *HTTPAuthentication `yaml:"authentication,omitempty" dc:"auth"`
	Body            string              `yaml:"body,omitempty"`
	SSLVerification bool                `yaml:"sslVerification,omitempty"`
}

func (r HTTPRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	if r.Method == "" {
		return fmt.Errorf("method cannot be empty")
	}

	for _, header := range r.Headers {
		if err := header.Validate(); err != nil {
			return fmt.Errorf("http header must be valid: %w", err)
		}
	}

	if err := r.validateAuth(); err != nil {
		return fmt.Errorf("http authentication must be valid: %w", err)
	}

	return nil
}

func (r HTTPRequest) validateAuth() error {
	if r.Authentication == nil {
		return nil
	}

	return r.Authentication.Validate()
}

type HTTPHeader struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func (h HTTPHeader) Validate() error {
	if h.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	return nil
}

type HTTPAuthentication struct {
	Type   string          `yaml:"type,omitempty"`
	Basic  *HTTPBasicAuth  `yaml:"basic,omitempty"`
	APIKey *HTTPAPIKeyAuth `yaml:"apiKey,omitempty"`
	Bearer *HTTPBearerAuth `yaml:"bearer,omitempty"`
}

func (a HTTPAuthentication) Validate() error {
	switch a.Type {
	case "basic":
		if err := a.Basic.Validate(); err != nil {
			return fmt.Errorf("basic authentication must be valid: %w", err)
		}
		return nil
	case "bearer":
		if err := a.Bearer.Validate(); err != nil {
			return fmt.Errorf("bearer authentication must be valid: %w", err)
		}
		return nil
	case "apiKey":
		if err := a.APIKey.Validate(); err != nil {
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

type HTTPBasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (ba HTTPBasicAuth) Validate() error {
	if ba.Username == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if ba.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

type HTTPBearerAuth struct {
	Token string `yaml:"token" dc:"bearer"`
}

func (ba HTTPBearerAuth) Validate() error {
	if ba.Token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	return nil
}

type HTTPAPIKeyAuth struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
	In    string `yaml:"in"`
}

func (aka HTTPAPIKeyAuth) Validate() error {
	if aka.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if aka.Value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	supportedInOptions := map[string]bool{
		"query":  true,
		"header": true,
	}

	if _, ok := supportedInOptions[aka.In]; !ok {
		return fmt.Errorf("in option \"%s\" is not supported. Only \"query\" and \"header\" are supported", aka.In)
	}

	return nil
}
