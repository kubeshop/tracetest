package trigger

import "google.golang.org/grpc/metadata"

const TriggerTypeGRPC TriggerType = "grpc"

type GRPCHeader struct {
	Key   string `expr_enabled:"true"`
	Value string `expr_enabled:"true"`
}

type GRPCRequest struct {
	ProtobufFile string `expr_enabled:"true"`
	Address      string `expr_enabled:"true"`
	Service      string `expr_enabled:"true"`
	Method       string `expr_enabled:"true"`
	Request      string `expr_enabled:"true"`
	Metadata     []GRPCHeader
	Auth         *HTTPAuthenticator
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
