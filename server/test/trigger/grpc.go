package trigger

import "google.golang.org/grpc/metadata"

const TriggerTypeGRPC TriggerType = "grpc"

type GRPCHeader struct {
	Key   string `expr_enabled:"true" json:"key"`
	Value string `expr_enabled:"true" json:"value"`
}

type GRPCRequest struct {
	ProtobufFile string             `json:"protobufFile,omitempty" expr_enabled:"true"`
	Address      string             `json:"address,omitempty" expr_enabled:"true"`
	Service      string             `json:"service,omitempty" expr_enabled:"true"`
	Method       string             `json:"method,omitempty" expr_enabled:"true"`
	Request      string             `json:"request,omitempty" expr_enabled:"true"`
	Metadata     []GRPCHeader       `json:"metadata,omitempty"`
	Auth         *HTTPAuthenticator `json:"auth,omitempty"`
}

func (a GRPCRequest) Headers() []string {
	h := []string{}

	for _, md := range a.Metadata {
		// ignore invalid values
		if md.Key == "" {
			continue
		}

		h = append(h, md.Key+": "+md.Value)
	}

	return h
}

func (a GRPCRequest) MD() *metadata.MD {
	md := metadata.MD{}

	for _, header := range a.Metadata {
		// ignore invalid values
		if header.Key == "" {
			continue
		}

		md[header.Key] = []string{header.Value}
	}

	return &md
}

func (a GRPCRequest) Authenticate() {
	if a.Auth == nil {
		return
	}

	a.Auth.AuthenticateGRPC()
}

type GRPCResponse struct {
	Status     string
	StatusCode int
	Metadata   []GRPCHeader
	Body       string
}
