package connection

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type portLinter struct {
	endpoints     []string
	expectedPorts []string
}

var _ TestStep = &portLinter{}

func PortLinter(expectedPorts []string, endpoints ...string) TestStep {
	return &portLinter{
		endpoints:     endpoints,
		expectedPorts: expectedPorts,
	}
}

func (s *portLinter) TestConnection(ctx context.Context) ConnectionTestStepResult {
	for _, endpoint := range s.endpoints {
		port := parsePort(endpoint)

		if !sliceContains(s.expectedPorts, port) {
			suggestedPorts := strings.Join(s.expectedPorts, ", ")
			return ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`port "%s" is not used frequently for connecting to this data store. Consider using %s instead if you experience problems connecting to it`, port, suggestedPorts),
				Status:               StatusWarning,
			}
		}
	}

	return ConnectionTestStepResult{
		OperationDescription: `You are using a commonly used port`,
		Status:               StatusPassed,
	}
}

func sliceContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}

	return false
}

var extractPortRegex = regexp.MustCompile("([0-9]+).*")

func parsePort(url string) string {
	index := strings.LastIndex(url, ":")
	if index < 0 {
		return ""
	}

	substring := url[index+1:]
	regexGroups := extractPortRegex.FindStringSubmatch(substring)
	if len(regexGroups) < 2 {
		return ""
	}

	return regexGroups[1]
}
