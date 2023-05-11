package resources_setup

import (
	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	resources_actions "github.com/kubeshop/tracetest/cli/resources/actions"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type Setup struct {
	*global_setup.Setup
	ResourceRegistry resources_actions.ResourceRegistry
	ResourceActions  resources_actions.ResourceActions
	options          []SetupOption
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

func WithResourceRegistry() SetupOption {
	return func(setup *Setup, args []string) error {
		resourceRegistry := resources_actions.NewResourceRegistry()
		baseOptions := []resources_actions.ResourceArgsOption{resources_actions.WithLogger(setup.Logger), resources_actions.WithConfig(*setup.Config)}

		configOptions := append(
			baseOptions,
			resources_actions.WithClient(utils.GetResourceAPIClient("configs", *setup.Config)),
			resources_actions.WithFormatter(resources_formatters.NewConfigFormatter()),
		)
		configActions := resources_actions.NewConfigActions(configOptions...)
		resourceRegistry.Register(configActions)

		pollingOptions := append(
			baseOptions,
			resources_actions.WithClient(utils.GetResourceAPIClient("pollingprofiles", *setup.Config)),
			resources_actions.WithFormatter(resources_formatters.NewPollingFormatter()),
		)
		pollingActions := resources_actions.NewPollingActions(pollingOptions...)
		resourceRegistry.Register(pollingActions)

		demoOptions := append(
			baseOptions,
			resources_actions.WithClient(utils.GetResourceAPIClient("demos", *setup.Config)),
			resources_actions.WithFormatter(resources_formatters.NewDemoFormatter()),
		)
		demoActions := resources_actions.NewDemoActions(demoOptions...)
		resourceRegistry.Register(demoActions)

		dataStoreOptions := append(
			baseOptions,
			resources_actions.WithClient(utils.GetResourceAPIClient("datastores", *setup.Config)),
			resources_actions.WithFormatter(resources_formatters.NewDatastoreFormatter()),
		)
		dataStoreActions := resources_actions.NewDataStoreActions(dataStoreOptions...)
		resourceRegistry.Register(dataStoreActions)

		environmentOptions := append(
			baseOptions,
			resources_actions.WithClient(utils.GetResourceAPIClient("environments", *setup.Config)),
			resources_actions.WithFormatter(resources_formatters.NewEnvironmentsFormatter()),
		)
		environmentActions := resources_actions.NewEnvironmentsActions(environmentOptions...)
		resourceRegistry.Register(environmentActions)

		setup.ResourceRegistry = resourceRegistry
		return nil
	}
}

func WithResourceActions() SetupOption {
	return func(setup *Setup, args []string) error {
		resourceType := args[0]
		resourceActions, err := setup.ResourceRegistry.Get(resourceType)
		if err != nil {
			return err
		}

		setup.ResourceActions = resourceActions
		return nil
	}
}
