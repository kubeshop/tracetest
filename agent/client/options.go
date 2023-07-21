package client

type Option func(*Client)

func WithAPIKey(apiKey string) Option {
	return func(c *Client) {
		c.config.APIKey = apiKey
	}
}

func WithAgentName(name string) Option {
	return func(c *Client) {
		c.config.AgentName = name
	}
}
