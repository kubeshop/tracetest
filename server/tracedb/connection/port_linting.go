package connection

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
)

type portLinter struct {
	dataStoreName string
	endpoints     []string
	expectedPorts []string
}

var _ TestStep = &portLinter{}

func PortLinter(dataStoreName string, expectedPorts []string, endpoints ...string) TestStep {
	return &portLinter{
		dataStoreName: dataStoreName,
		endpoints:     endpoints,
		expectedPorts: expectedPorts,
	}
}

func (s *portLinter) TestConnection(ctx context.Context) model.ConnectionTestStep {
	for _, endpoint := range s.endpoints {
		port := parsePort(endpoint)

		if !sliceContains(s.expectedPorts, port) {
			suggestedPorts := formatAvailablePortsMessage(s.expectedPorts)
			return model.ConnectionTestStep{
				Message: fmt.Sprintf(`For %s, port "%s" is not the default port for accessing traces programmatically. Typically, %s uses port %s. If you continue experiencing issues, you may want to verify the correct port to specify.`, s.dataStoreName, port, s.dataStoreName, suggestedPorts),
				Status:  model.StatusWarning,
			}
		}
	}

	return model.ConnectionTestStep{
		Message: `You are using a commonly used port`,
		Status:  model.StatusPassed,
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

// Generates the ports separated by commas and "or"
// ["123"] => 123
// ["123", "345"] => 123 or 345
// ["123", "345", "567"] => 123, 345, or 567
func formatAvailablePortsMessage(ports []string) string {
	if len(ports) == 1 {
		return ports[0]
	}

	allPortsExceptLast := ports[0 : len(ports)-1]
	portsSeparatedByComma := strings.Join(allPortsExceptLast, ", ")
	lastPort := ports[len(ports)-1]

	if len(ports) == 2 {
		return fmt.Sprintf("%s or %s", portsSeparatedByComma, lastPort)
	}

	return fmt.Sprintf("%s, or %s", portsSeparatedByComma, lastPort)
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

	port := regexGroups[1]

	if port == "1" {
		if strings.Contains(url, "http") {
			return "80"
		} else {
			return "443"
		}
	}

	return port
}
