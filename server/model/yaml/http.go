package yaml

import (
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

type HTTPHeaders []HTTPHeader

func (hs HTTPHeaders) Model() []model.HTTPHeader {
	if len(hs) == 0 {
		return nil
	}
	mh := make([]model.HTTPHeader, 0, len(hs))
	for _, h := range hs {
		mh = append(mh, model.HTTPHeader{
			Key:   h.Key,
			Value: h.Value,
		})
	}

	return mh
}

type HTTPRequest struct {
	URL            string             `yaml:"url"`
	Method         string             `yaml:"method"`
	Headers        HTTPHeaders        `yaml:"headers,omitempty"`
	Authentication HTTPAuthentication `yaml:"authentication,omitempty"`
	Body           string             `yaml:"body,omitempty"`
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

	if err := r.Authentication.Validate(); err != nil {
		return fmt.Errorf("http authentication must be valid: %w", err)
	}

	return nil
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
	Type   string         `yaml:"type,omitempty"`
	Basic  HTTPBasicAuth  `yaml:"basic,omitempty"`
	ApiKey HTTPAPIKeyAuth `yaml:"apiKey,omitempty"`
	Bearer HTTPBearerAuth `yaml:"bearer,omitempty"`
}

func (a HTTPAuthentication) Model() *model.HTTPAuthenticator {
	var props map[string]string
	switch a.Type {
	case "":
		// auth not set
		return nil
	case "apiKey":
		props = map[string]string{
			"key":   a.ApiKey.Key,
			"value": a.ApiKey.Value,
			"in":    a.ApiKey.In,
		}
	case "basic":
		props = map[string]string{
			"username": a.Basic.User,
			"password": a.Basic.Password,
		}
	case "bearer":
		props = map[string]string{
			"token": a.Bearer.Token,
		}
	}

	return &model.HTTPAuthenticator{
		Type:  a.Type,
		Props: props,
	}
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
		if err := a.ApiKey.Validate(); err != nil {
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
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (ba HTTPBasicAuth) Validate() error {
	if ba.User == "" {
		return fmt.Errorf("user cannot be empty")
	}

	if ba.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	return nil
}

type HTTPBearerAuth struct {
	Token string `yaml:"token"`
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
