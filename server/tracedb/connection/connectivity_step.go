package connection

import (
	"context"
	"fmt"
	"strings"
)

type connectivityTestStep struct {
	endpoints []string
}

var _ TestStep = &connectivityTestStep{}

func (s *connectivityTestStep) TestConnection(ctx context.Context) ConnectionTestStepResult {
	unreachableEndpoints := make([]string, 0)
	var connectionErr error
	for _, endpoint := range s.endpoints {
		reachable, err := IsReachable(endpoint)
		if !reachable {
			unreachableEndpoints = append(unreachableEndpoints, fmt.Sprintf(`"%s"`, endpoint))
			connectionErr = err
		}
	}

	if len(unreachableEndpoints) > 0 {
		endpoints := strings.Join(unreachableEndpoints, ", ")
		return ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf("Tracetest tried to connect to the following endpoints and failed: %s", endpoints),
			Error:                connectionErr,
		}
	}

	endpoints := strings.Join(s.endpoints, ", ")
	return ConnectionTestStepResult{
		OperationDescription: fmt.Sprintf(`Tracetest connected to "%s"`, endpoints),
	}
}

func ConnectivityStep(endpoints ...string) TestStep {
	return &connectivityTestStep{
		endpoints: endpoints,
	}
}
