package misc_setup

import (
	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type Setup struct {
	*global_setup.Setup
	Client  *openapi.APIClient
	options []SetupOption
}

type SetupOption = func(setup *Setup, args []string) error

func NewSetup(global *global_setup.Setup, options ...SetupOption) *Setup {
	return &Setup{
		Setup:   global,
		options: options,
	}
}

func (s *Setup) PreRun(cmd *cobra.Command, args []string) (string, error) {
	_, err := s.Setup.PreRun(cmd, args)
	if err != nil {
		return "", err
	}

	for _, option := range s.options {
		err := option(s, args)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func WithApiClient() SetupOption {
	return func(setup *Setup, _ []string) error {
		setup.Client = utils.GetAPIClient(*setup.Config)

		return nil
	}
}
