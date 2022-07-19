package model

const TriggerTypeGRPC TriggerType = "grpc"

type GRPCHeader struct {
	Key, Value string
}

type GRPCRequest struct {
	ProtobufFile string
	Address      string
	Service      string
	Method       string
	Metadata     []GRPCHeader
	Auth         *HTTPAuthenticator
	Request      string
}

func (a GRPCRequest) Headers() []string {
	h := []string{}

	for _, md := range a.Metadata {
		h = append(h, md.Key+": "+md.Value)
	}

	return h
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
