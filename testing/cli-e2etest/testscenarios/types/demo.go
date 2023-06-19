package types

// Note: these types are very similar to the types on the server folder
// however they are defined here to avoid bias with the current implementation

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
