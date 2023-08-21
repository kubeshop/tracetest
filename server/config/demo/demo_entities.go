package demo

import (
	"github.com/kubeshop/tracetest/server/pkg/id"
)

type DemoType string

const (
	DemoTypePokeshop           DemoType = "pokeshop"
	DemoTypeOpentelemetryStore DemoType = "otelstore"
)

const (
	ResourceName       = "Demo"
	ResourceNamePlural = "Demos"
)

type Demo struct {
	ID                 id.ID                   `json:"id"`
	Name               string                  `json:"name"`
	Type               DemoType                `json:"type"`
	Enabled            bool                    `json:"enabled"`
	Pokeshop           *PokeshopDemo           `json:"pokeshop,omitempty"`
	OpenTelemetryStore *OpenTelemetryStoreDemo `json:"opentelemetryStore,omitempty"`
}

func (d Demo) HasID() bool {
	return d.ID != ""
}

func (d Demo) GetID() id.ID {
	return d.ID
}

func (d Demo) Validate() error {
	return nil
}

type PokeshopDemo struct {
	HTTPEndpoint  string `json:"httpEndpoint,omitempty"`
	GRPCEndpoint  string `json:"grpcEndpoint,omitempty"`
	KafkaEndpoint string `json:"kafkaEndpoint,omitempty"`
}

type OpenTelemetryStoreDemo struct {
	FrontendEndpoint       string `json:"frontendEndpoint,omitempty"`
	ProductCatalogEndpoint string `json:"productCatalogEndpoint,omitempty"`
	CartEndpoint           string `json:"cartEndpoint,omitempty"`
	CheckoutEndpoint       string `json:"checkoutEndpoint,omitempty"`
}
