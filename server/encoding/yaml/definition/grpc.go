package definition

import (
	"fmt"
)

type GRPC struct {
	ProtobufFile string             `json:"protobufFile" yaml:"protobufFile"`
	Address      string             `json:"address" yaml:"address"`
	Method       string             `json:"method" yaml:"method"`
	Metadata     []GRPCHeader       `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Auth         HTTPAuthentication `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	Request      string             `json:"request,omitempty" yaml:"request,omitempty"`
}

func (r GRPC) Validate() error {
	if r.ProtobufFile == "" {
		return fmt.Errorf("protobufFile cannot be empty")
	}

	if r.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if r.Method == "" {
		return fmt.Errorf("method cannot be empty")
	}

	for _, header := range r.Metadata {
		if err := header.Validate(); err != nil {
			return fmt.Errorf("http header must be valid: %w", err)
		}
	}

	if err := r.Auth.Validate(); err != nil {
		return fmt.Errorf("http authentication must be valid: %w", err)
	}

	return nil
}

type GRPCHeader struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

func (h GRPCHeader) Validate() error {
	if h.Key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	return nil
}
