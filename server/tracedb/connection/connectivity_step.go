package connection

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/kubeshop/tracetest/server/model"
)

type connectivityTestStep struct {
	endpoints []string
	protocol  model.Protocol
}

var _ TestStep = &connectivityTestStep{}

func (s *connectivityTestStep) TestConnection(_ context.Context) model.ConnectionTestStep {
	unreachableEndpoints := make([]string, 0)
	var connectionErr error
	for _, endpoint := range s.endpoints {
		err := CheckReachability(endpoint, s.protocol)
		if err != nil {
			unreachableEndpoints = append(unreachableEndpoints, fmt.Sprintf(`"%s"`, endpoint))
			connectionErr = multierror.Append(
				connectionErr,
				fmt.Errorf("cannot connect to endpoint '%s': %w", endpoint, err),
			)
		}
	}

	if len(s.endpoints) == 0 {
		return model.ConnectionTestStep{
			Message: "Tracetest tried to connect but no endpoints were provided",
			Error:   fmt.Errorf("no endpoints provided"),
		}
	}

	if connectionErr != nil {
		endpoints := strings.Join(unreachableEndpoints, ", ")
		return model.ConnectionTestStep{
			Message: fmt.Sprintf("Tracetest tried to connect to the following endpoints and failed: %s", endpoints),
			Status:  model.StatusFailed,
			Error:   connectionErr,
		}
	}

	endpoints := strings.Join(s.endpoints, ", ")
	return model.ConnectionTestStep{
		Message: fmt.Sprintf(`Tracetest connected to %s`, endpoints),
		Status:  model.StatusPassed,
	}
}

func ConnectivityStep(protocol model.Protocol, endpoints ...string) TestStep {
	return &connectivityTestStep{
		endpoints: endpoints,
		protocol:  protocol,
	}
}
