package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/processor"
	"github.com/kubeshop/tracetest/cli/processor/trigger_preprocessor"
	"github.com/kubeshop/tracetest/cli/runner"
)

var resourceParams = &resourceParameters{}

var (
	testSuiteRunner = runner.TestSuiteRunner(
		testSuiteClient,
		openapiClient,
		formatters.TestSuiteRun(func() string { return cliConfig.UI() }, true),
	)

	runnerRegistry = runner.NewRegistry(cliLogger).
			Register(runner.TestRunner(
			testClient,
			openapiClient,
			formatters.TestRun(func() string { return cliConfig.UI() }, true),
		)).
		Register(testSuiteRunner).
		RegisterProxy("transaction", testSuiteRunner.Name())
)

var (
	httpClient = &resourcemanager.HTTPClient{}

	variableSetPreprocessor = processor.VariableSet(cliLogger)
	variableSetClient       = GetVariableSetClient(httpClient, variableSetPreprocessor)

	pollingProfileClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"pollingprofile", "pollingprofiles",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "STRATEGY", Path: "spec.strategy"},
			},
		}),
		resourcemanager.WithResourceType("PollingProfile"),
	)

	envTokensClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"environmenttoken", "environmenttokens",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "ROLE", Path: "spec.role"},
				{Header: "REVOKED", Path: "spec.isRevoked"},
				{Header: "ISSUED AT", Path: "spec.issuedAt"},
				{Header: "EXPIRES AT", Path: "spec.expiresAt"},
			},
		}),
		resourcemanager.WithResourceType("EnvironmentToken"),
	)

	orgInvitesClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"invite", "invites",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "TO", Path: "spec.to"},
				{Header: "Role", Path: "spec.role"},
				{Header: "TYPE", Path: "spec.type"},
				{Header: "STATUS", Path: "spec.status"},
			},
		}),
		resourcemanager.WithResourceType("Invite"),
	)

	environmentPostproessor = processor.Environment(cliLogger, func(ctx context.Context, input fileutil.File) (fileutil.File, error) {
		client, err := envResources.Get(strings.ToLower(input.Type()))
		if err != nil {
			return input, fmt.Errorf("cannot get client for resource type '%s': %w", input.Type(), err)
		}

		_, err = client.Apply(ctx, input, resourcemanager.Formats.Get(resourcemanager.FormatYAML))
		if err != nil {
			return input, fmt.Errorf("cannot apply resource: %w", err)
		}

		return input, nil
	}, func(ctx context.Context, envID string) error {
		cliConfig.EnvironmentID = envID
		setupResources()
		return nil
	})

	environmentClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"environment", "environments",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
			},
		}),
		resourcemanager.WithResourceType("Environment"),
		resourcemanager.WithApplyPostProcessor(environmentPostproessor.Postprocess),
	)

	triggerPreprocessorRegistry = trigger_preprocessor.NewRegistry(cliLogger).
					Register(trigger_preprocessor.GRPC(cliLogger)).
					Register(trigger_preprocessor.PLAYWRIGHTENGINE(cliLogger)).
					Register(trigger_preprocessor.GRAPHQL(cliLogger)).
					Register(trigger_preprocessor.HTTP(cliLogger))

	testPreprocessor = processor.Test(cliLogger, triggerPreprocessorRegistry, func(ctx context.Context, input fileutil.File) (fileutil.File, error) {
		updated, err := pollingProfileClient.Apply(ctx, input, resourcemanager.Formats.Get(resourcemanager.FormatYAML))
		if err != nil {
			return input, fmt.Errorf("cannot apply polling profile: %w", err)
		}

		return fileutil.New(input.AbsPath(), []byte(updated)), nil
	})

	testApplyFn = func(ctx context.Context, input fileutil.File) (fileutil.File, error) {
		updated, err := testClient.Apply(ctx, input, resourcemanager.Formats.Get(resourcemanager.FormatYAML))
		if err != nil {
			return input, fmt.Errorf("cannot apply test: %w", err)
		}

		return fileutil.New(input.AbsPath(), []byte(updated)), nil
	}

	testSuiteApplyFn = func(ctx context.Context, input fileutil.File) (fileutil.File, error) {
		updated, err := testSuiteClient.Apply(ctx, input, resourcemanager.Formats.Get(resourcemanager.FormatYAML))
		if err != nil {
			return input, fmt.Errorf("cannot apply suite: %w", err)
		}

		return fileutil.New(input.AbsPath(), []byte(updated)), nil
	}

	monitorPreprocessor = processor.Monitor(cliLogger, testSuiteApplyFn, testApplyFn)

	monitorClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"monitor", "monitors",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "VERSION", Path: "spec.version"},
				{Header: "RUNS", Path: "spec.summary.runs"},
				{Header: "LAST RUN TIME", Path: "spec.summary.lastRunTime"},
				{Header: "LAST RUN STATE", Path: "spec.summary.lastRunState"},
				{Header: "URL", Path: "spec.url"},
			},
			ItemModifier: func(item *gabs.Container) error {
				// set spec.summary.steps to the number of steps in the test suite
				id, ok := item.Path("spec.id").Data().(string)
				if !ok {
					return fmt.Errorf("monitor id '%s' is not a string", id)
				}

				url := cliConfig.URL() + "/monitor/" + id
				item.SetP(url, "spec.url")

				if err := formatItemDate(item, "spec.summary.lastRunTime"); err != nil {
					return err
				}

				return nil
			},
		}),
		resourcemanager.WithApplyPreProcessor(monitorPreprocessor.Preprocess),
	)

	testClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"test", "tests",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "VERSION", Path: "spec.version"},
				{Header: "TRIGGER TYPE", Path: "spec.trigger.type"},
				{Header: "RUNS", Path: "spec.summary.runs"},
				{Header: "LAST RUN TIME", Path: "spec.summary.lastRun.time"},
				{Header: "LAST RUN SUCCESSES", Path: "spec.summary.lastRun.passes"},
				{Header: "LAST RUN FAILURES", Path: "spec.summary.lastRun.fails"},
				{Header: "URL", Path: "spec.url"},
			},
			ItemModifier: func(item *gabs.Container) error {
				// set spec.summary.steps to the number of steps in the test suite
				id, ok := item.Path("spec.id").Data().(string)
				if !ok {
					return fmt.Errorf("test id '%s' is not a string", id)
				}

				url := cliConfig.URL() + "/test/" + id
				item.SetP(url, "spec.url")

				if err := formatItemDate(item, "spec.summary.lastRun.time"); err != nil {
					return err
				}

				return nil
			},
		}),
		resourcemanager.WithApplyPreProcessor(testPreprocessor.Preprocess),
	)

	testSuitePreprocessor = processor.TestSuite(cliLogger, testApplyFn)

	testSuiteClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"testsuite", "testsuites",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "VERSION", Path: "spec.version"},
				{Header: "STEPS", Path: "spec.summary.steps"},
				{Header: "RUNS", Path: "spec.summary.runs"},
				{Header: "LAST RUN TIME", Path: "spec.summary.lastRun.time"},
				{Header: "LAST RUN SUCCESSES", Path: "spec.summary.lastRun.passes"},
				{Header: "LAST RUN FAILURES", Path: "spec.summary.lastRun.fails"},
			},
			ItemModifier: func(item *gabs.Container) error {
				// set spec.summary.steps to the number of steps in the test suite
				item.SetP(len(item.Path("spec.steps").Children()), "spec.summary.steps")

				if err := formatItemDate(item, "spec.summary.lastRun.time"); err != nil {
					return err
				}

				return nil
			},
		}),
		resourcemanager.WithResourceType("TestSuite"),
		resourcemanager.WithApplyPreProcessor(testSuitePreprocessor.Preprocess),
		resourcemanager.WithDeprecatedAlias("Transaction"),
	)

	deprecatedTransactionsClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"transaction", "transactions",
		resourcemanager.WithProxyResource("testsuite"),
	)

	envResources = resourcemanager.NewRegistry().
			Register(
			resourcemanager.NewClient(
				httpClient, cliLogger,
				"config", "configs",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "ANALYTICS ENABLED", Path: "spec.analyticsEnabled"},
					},
				}),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient, cliLogger,
				"analyzer", "analyzers",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "ENABLED", Path: "spec.enabled"},
						{Header: "MINIMUM SCORE", Path: "spec.minimumScore"},
						{Header: "PLUGINS", Path: "spec.total.plugins"},
					},
					ItemModifier: func(item *gabs.Container) error {
						item.SetP(len(item.Path("spec.plugins").Children()), "spec.total.plugins")

						return nil
					},
				}),
			),
		).
		Register(monitorClient).
		Register(envTokensClient).
		Register(orgInvitesClient).
		Register(pollingProfileClient).
		Register(
			resourcemanager.NewClient(
				httpClient, cliLogger,
				"demo", "demos",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "TYPE", Path: "spec.type"},
						{Header: "ENABLED", Path: "spec.enabled"},
					},
				}),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient, cliLogger,
				"datastore", "datastores",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "DEFAULT", Path: "spec.default"},
					},
					ItemModifier: func(item *gabs.Container) error {
						isDefault := item.Path("spec.default").Data().(bool)
						if !isDefault {
							item.SetP("", "spec.default")
						} else {
							item.SetP("*", "spec.default")
						}
						return nil
					},
				}),
				resourcemanager.WithResourceType("DataStore"),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient, cliLogger,
				"testrunner", "testrunners",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "REQUIRED GATES", Path: "spec.gates"},
					},
					ItemModifier: func(item *gabs.Container) error {
						gates := []string{}
						for _, gate := range item.Path("spec.requiredGates").Children() {
							gates = append(gates, "- "+gate.Data().(string))
						}
						item.SetP(strings.Join(gates, "\n"), "spec.gates")
						return nil
					},
				}),
				resourcemanager.WithResourceType("TestRunner"),
			),
		).
		Register(variableSetClient).
		Register(testSuiteClient).
		Register(testClient)

	resources = envResources.
			Register(organizationsClient).
			Register(environmentClient).
			Register(environmentMeClient).

		// deprecated resources
		Register(deprecatedTransactionsClient)

	organizationsClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"organization", "organizations",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "id"},
				{Header: "NAME", Path: "name"},
			},
		}),
		resourcemanager.WithListPath("elements"),
	)

	environmentMeClient = resourcemanager.NewClient(
		httpClient, cliLogger,
		"env", "environments",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "id"},
				{Header: "NAME", Path: "name"},
			},
		}),
		resourcemanager.WithPrefixGetter(func() string { return fmt.Sprintf("/organizations/%s/", cliConfig.OrganizationID) }),
		resourcemanager.WithListPath("elements"),
	)
)

func resourceList() string {
	return strings.Join(resources.List(), "|")
}

func runnableResourceList() string {
	return strings.Join(runnerRegistry.List(), "|")
}

func setupResources() {
	extraHeaders := http.Header{}
	extraHeaders.Set("x-client-id", analytics.ClientID())
	extraHeaders.Set("x-source", "cli")
	extraHeaders.Set("x-organization-id", cliConfig.OrganizationID)
	extraHeaders.Set("x-environment-id", cliConfig.EnvironmentID)
	extraHeaders.Set("Authorization", fmt.Sprintf("Bearer %s", cliConfig.Jwt))

	// if cliConfig has SkipVerify set to true, use that value.
	// otherwise use the value from the flag
	if cliConfig.SkipVerify {
		skipVerify = true
	}

	// To avoid a ciruclar reference initialization when setting up the registry and its resources,
	// we create the resources with a pointer to an unconfigured HTTPClient.
	// When each command is run, this function is run in the PreRun stage, before any of the actual `Run` code is executed
	// We take this chance to configure the HTTPClient with the correct URL and headers.
	// To make this configuration propagate to all the resources, we need to replace the pointer to the HTTPClient.
	// For more details, see https://github.com/kubeshop/tracetest/pull/2832#discussion_r1245616804

	hc := resourcemanager.NewHTTPClient(fmt.Sprintf("%s%s", cliConfig.URL(), cliConfig.Path()), extraHeaders, skipVerify)
	*httpClient = *hc
}

func formatItemDate(item *gabs.Container, path string) error {
	rawDate := item.Path(path).Data()
	if rawDate == nil {
		return nil
	}
	dateStr := rawDate.(string)
	// if field is empty, do nothing
	if dateStr == "" {
		return nil
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse datetime field '%s' (value '%s'): %s", path, dateStr, err)
	}

	if date.IsZero() {
		// sometime the date comes like 0000-00-00T00:00:00Z... show nothing in that case
		item.SetP("", path)
		return nil
	}

	item.SetP(date.Format(time.DateTime), path)
	return nil
}

func GetVariableSetClient(httpClient *resourcemanager.HTTPClient, preprocessor processor.Preprocessor) resourcemanager.Client {
	variableSetClient := resourcemanager.NewClient(
		httpClient, cmdutil.GetLogger(),
		"variableset", "variablesets",
		resourcemanager.WithTableConfig(resourcemanager.TableConfig{
			Cells: []resourcemanager.TableCellConfig{
				{Header: "ID", Path: "spec.id"},
				{Header: "NAME", Path: "spec.name"},
				{Header: "DESCRIPTION", Path: "spec.description"},
			},
		}),
		resourcemanager.WithResourceType("VariableSet"),
		resourcemanager.WithApplyPreProcessor(preprocessor.Preprocess),
		resourcemanager.WithDeprecatedAlias("Environment"),
	)

	return variableSetClient
}
