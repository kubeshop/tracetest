package config

type Mode string

const (
	Mode_Desktop Mode = "desktop"
	Mode_Verbose Mode = "verbose"
)

type Flags struct {
	ServerURL         string
	OrganizationID    string
	EnvironmentID     string
	CI                bool
	AgentApiKey       string
	Token             string
	Mode              Mode
	LogLevel          string
	CollectorEndpoint string
}

func (f Flags) AutomatedEnvironmentCanBeInferred() bool {
	return f.CI || f.AgentApiKey != "" || f.Token != ""
}
