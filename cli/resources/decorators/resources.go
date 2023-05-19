package resources_decorators

import (
	global_decorators "github.com/kubeshop/tracetest/cli/global/decorators"
	global_types "github.com/kubeshop/tracetest/cli/global/types"
	resources_actions "github.com/kubeshop/tracetest/cli/resources/actions"
	resources_formatters "github.com/kubeshop/tracetest/cli/resources/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type resources struct {
	global_decorators.ResultHandler
	global_decorators.Analytics
	ResourceRegistry resources_actions.ResourceRegistry
	ResourceActions  resources_actions.ResourceActions
}

type Resources interface {
	global_decorators.Logger
	global_decorators.Analytics
	global_decorators.Version
	global_decorators.Config
	global_decorators.ResultHandler
	SetResourceRegistry(resources_actions.ResourceRegistry)
	GetResourceRegistry() resources_actions.ResourceRegistry
	SetResourceActions(resources_actions.ResourceActions)
	GetResourceActions() resources_actions.ResourceActions
}

var _ Resources = &resources{}

func WithResources(command global_types.Command) global_types.Command {
	resultHandler, err := command.(global_decorators.ResultHandler)
	if !err {
		panic("command must implement Config interface")
	}

	analytics, err := command.(global_decorators.Analytics)
	if !err {
		panic("command must implement Config interface")
	}

	resources := resources{
		ResultHandler: resultHandler,
		Analytics:     analytics,
	}

	cmd := resources.Get()
	cmd.PreRun = resources.preRun(cmd.PreRun)
	resources.Set(cmd)

	return global_decorators.Decorate[Resources](
		resources,
		global_decorators.WithAnalytics("apply", "resources"),
		global_decorators.WithLogger,
		global_decorators.WithConfig,
		global_decorators.WithVersion,
		global_decorators.WithResultHandler,
	)
}

func (d *resources) preRun(next global_decorators.CobraFn) global_decorators.CobraFn {
	return func(cmd *cobra.Command, args []string) {
		d.setupResourceRegistry(cmd, args)
		d.setupResourceActions(cmd, args)

		next(cmd, args)
	}
}

func (d *resources) SetResourceRegistry(registry resources_actions.ResourceRegistry) {
	d.ResourceRegistry = registry
}

func (d *resources) GetResourceRegistry() resources_actions.ResourceRegistry {
	return d.ResourceRegistry
}

func (d *resources) SetResourceActions(actions resources_actions.ResourceActions) {
	d.ResourceActions = actions
}

func (d *resources) GetResourceActions() resources_actions.ResourceActions {
	return d.ResourceActions
}

func (d *resources) setupResourceRegistry(cmd *cobra.Command, args []string) {
	config := d.GetConfig()
	logger := d.GetLogger()
	resourceRegistry := resources_actions.NewResourceRegistry()
	baseOptions := []resources_actions.ResourceArgsOption{resources_actions.WithLogger(logger), resources_actions.WithConfig(*config)}

	configOptions := append(
		baseOptions,
		resources_actions.WithClient(utils.GetResourceAPIClient("configs", *config)),
		resources_actions.WithFormatter(resources_formatters.NewConfigFormatter()),
	)
	configActions := resources_actions.NewConfigActions(configOptions...)
	resourceRegistry.Register(configActions)

	pollingOptions := append(
		baseOptions,
		resources_actions.WithClient(utils.GetResourceAPIClient("pollingprofiles", *config)),
		resources_actions.WithFormatter(resources_formatters.NewPollingFormatter()),
	)
	pollingActions := resources_actions.NewPollingActions(pollingOptions...)
	resourceRegistry.Register(pollingActions)

	demoOptions := append(
		baseOptions,
		resources_actions.WithClient(utils.GetResourceAPIClient("demos", *config)),
		resources_actions.WithFormatter(resources_formatters.NewDemoFormatter()),
	)
	demoActions := resources_actions.NewDemoActions(demoOptions...)
	resourceRegistry.Register(demoActions)

	dataStoreOptions := append(
		baseOptions,
		resources_actions.WithClient(utils.GetResourceAPIClient("datastores", *config)),
		resources_actions.WithFormatter(resources_formatters.NewDatastoreFormatter()),
	)
	dataStoreActions := resources_actions.NewDataStoreActions(dataStoreOptions...)
	resourceRegistry.Register(dataStoreActions)

	environmentOptions := append(
		baseOptions,
		resources_actions.WithClient(utils.GetResourceAPIClient("environments", *config)),
		resources_actions.WithFormatter(resources_formatters.NewEnvironmentsFormatter()),
	)
	environmentActions := resources_actions.NewEnvironmentsActions(environmentOptions...)
	resourceRegistry.Register(environmentActions)

	d.SetResourceRegistry(resourceRegistry)
}

func (d *resources) setupResourceActions(cmd *cobra.Command, args []string) {
	resourceType := args[0]
	resourceActions, err := d.GetResourceRegistry().Get(resourceType)
	if err != nil {
		d.Error(err)
		return
	}

	d.SetResourceActions(resourceActions)
}
