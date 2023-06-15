package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

// DataStore
type DataStore struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type DataStoreResource struct {
	Type string    `json:"type"`
	Spec DataStore `json:"spec"`
}

// Environment

type EnvironmentKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Environment struct {
	ID     string                `json:"id"`
	Name   string                `json:"name"`
	Values []EnvironmentKeyValue `json:"values"`
}

type EnvironmentResource struct {
	Type string      `json:"type"`
	Spec Environment `json:"spec"`
}

// Config

type Config struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	AnalyticsEnabled bool `json:"analyticsEnabled"`
}

type ConfigResource struct {
	Type string `json:"type"`
	Spec Config `json:"spec"`
}

// PollingProfile

type PollingProfilePeriodicStrategy struct {
	Timeout              string `json:"timeout"`
	RetryDelay           string `json:"retryDelay"`
	SelectorMatchRetries string `json:"selectorMatchRetries"`
}

type PollingProfile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Default bool   `json:"default"`

	Strategy string                         `json:"strategy"`
	Periodic PollingProfilePeriodicStrategy `json:"periodic"`
}

type PollingProfileResource struct {
	Type string         `json:"type"`
	Spec PollingProfile `json:"spec"`
}

// Demo

type DemoPokeshop struct {
	HttpEndpoint string `json:"httpEndpoint"`
	GrpcEndpoint string `json:"grpcEndpoint"`
}

type DemoOTelStore struct {
	FrontendEndpoint       string `json:"frontendEndpoint"`
	ProductCatalogEndpoint string `json:"productCatalogEndpoint"`
	CartEndpoint           string `json:"cartEndpoint"`
	CheckoutEndpoint       string `json:"checkoutEndpoint"`
}

type Demo struct {
	Id        string        `json:"id"`
	Name      string        `json:"name"`
	Enabled   bool          `json:"enabled"`
	Type      string        `json:"type"`
	OTelStore DemoOTelStore `json:"opentelemetryStore"`
	Pokeshop  DemoPokeshop  `json:"pokeshop"`
}

type DemoResource struct {
	Type string `json:"type"`
	Spec Demo   `json:"spec"`
}

type AnalyzerResource struct {
	Type string   `json:"type"`
	Spec Analyzer `json:"spec"`
}

type Analyzer struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Enabled      bool             `json:"enabled"`
	MinimumScore int              `json:"minimumScore"`
	Plugins      []AnalyzerPlugin `json:"plugins"`
}

type AnalyzerPlugin struct {
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
	Required bool   `json:"required"`
}
