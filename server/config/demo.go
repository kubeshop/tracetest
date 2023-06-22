package config

type demo struct {
	Enabled   []string `yaml:",omitempty" mapstructure:"enabled"`
	Endpoints struct {
		PokeshopHttp       string `yaml:",omitempty" mapstructure:"pokeshopHttp"`
		PokeshopGrpc       string `yaml:",omitempty" mapstructure:"pokeshopGrpc"`
		OtelFrontend       string `yaml:",omitempty" mapstructure:"otelFrontend"`
		OtelProductCatalog string `yaml:",omitempty" mapstructure:"otelProductCatalog"`
		OtelCart           string `yaml:",omitempty" mapstructure:"otelCart"`
		OtelCheckout       string `yaml:",omitempty" mapstructure:"otelCheckout"`
	} `yaml:",omitempty" mapstructure:"endpoints"`
}

func (c *AppConfig) DemoEnabled() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.Demo.Enabled
}

func (c *AppConfig) DemoEndpoints() map[string]string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return map[string]string{
		"PokeshopHttp":       c.config.Demo.Endpoints.PokeshopHttp,
		"PokeshopGrpc":       c.config.Demo.Endpoints.PokeshopGrpc,
		"OtelFrontend":       c.config.Demo.Endpoints.OtelFrontend,
		"OtelProductCatalog": c.config.Demo.Endpoints.OtelProductCatalog,
		"OtelCart":           c.config.Demo.Endpoints.OtelCart,
		"OtelCheckout":       c.config.Demo.Endpoints.OtelCheckout,
	}
}
