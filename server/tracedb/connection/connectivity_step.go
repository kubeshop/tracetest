package connection

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type connectivityTestStep struct {
	endpoints []string
	protocol  Protocol
}

var _ TestStep = &connectivityTestStep{}

func (s *connectivityTestStep) TestConnection(_ context.Context) ConnectionTestStepResult {
	unreachableEndpoints := make([]string, 0)
	var connectionErr error
	for _, endpoint := range s.endpoints {
		err := IsReachable(endpoint, s.protocol)
		if err != nil {
			unreachableEndpoints = append(unreachableEndpoints, fmt.Sprintf(`"%s"`, endpoint))
			connectionErr = multierror.Append(
				connectionErr,
				fmt.Errorf("cannot connect to endpoint '%s': %w", endpoint, err),
			)
		}
	}

	if connectionErr != nil {
		endpoints := strings.Join(unreachableEndpoints, ", ")
		return ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf("Tracetest tried to connect to the following endpoints and failed: %s", endpoints),
			Error:                connectionErr,
		}
	}

	endpoints := strings.Join(s.endpoints, ", ")
	return ConnectionTestStepResult{
		OperationDescription: fmt.Sprintf(`Tracetest connected to %s`, endpoints),
	}
}

func ConnectivityStep(protocol Protocol, endpoints ...string) TestStep {
	return &connectivityTestStep{
		endpoints: endpoints,
	}
}
