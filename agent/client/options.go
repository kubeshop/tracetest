package client

import "time"

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

func WithPingPeriod(period time.Duration) Option {
	return func(c *Client) {
		c.config.PingPeriod = period
	}
}
