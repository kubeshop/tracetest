package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/preprocessor"
)

var resourceParams = &resourceParameters{}
var (
	httpClient = &resourcemanager.HTTPClient{}

	testPreprocessor = preprocessor.Test(cliLogger)
	testClient       = resourcemanager.NewClient(
		httpClient,
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
				// set spec.summary.steps to the number of steps in the transaction
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

	transactionPreprocessor = preprocessor.Transaction(cliLogger, func(ctx context.Context, input fileutil.File) (fileutil.File, error) {
		formatYAML, err := resourcemanager.Formats.Get(resourcemanager.FormatYAML)
		if err != nil {
			return input, fmt.Errorf("cannot get yaml format: %w", err)
		}
		updated, err := testClient.Apply(ctx, input, formatYAML)
		if err != nil {
			return input, fmt.Errorf("cannot apply test: %w", err)
		}

		return fileutil.New(input.AbsPath(), []byte(updated)), nil
	})

	resources = resourcemanager.NewRegistry().
			Register(
			resourcemanager.NewClient(
				httpClient,
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
				httpClient,
				"analyzer", "analyzers",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "ENABLED", Path: "spec.enabled"},
						{Header: "MINIMUM SCORE", Path: "spec.minimumScore"},
					},
				}),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient,
				"pollingprofile", "pollingprofiles",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "STRATEGY", Path: "spec.strategy"},
					},
				}),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient,
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
				httpClient,
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
				resourcemanager.WithDeleteSuccessMessage("DataStore removed. Defaulting back to no-tracing mode"),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient,
				"environment", "environments",
				resourcemanager.WithTableConfig(resourcemanager.TableConfig{
					Cells: []resourcemanager.TableCellConfig{
						{Header: "ID", Path: "spec.id"},
						{Header: "NAME", Path: "spec.name"},
						{Header: "DESCRIPTION", Path: "spec.description"},
					},
				}),
			),
		).
		Register(
			resourcemanager.NewClient(
				httpClient,
				"transaction", "transactions",
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
						// set spec.summary.steps to the number of steps in the transaction
						item.SetP(len(item.Path("spec.steps").Children()), "spec.summary.steps")

						if err := formatItemDate(item, "spec.summary.lastRun.time"); err != nil {
							return err
						}

						return nil
					},
				}),
				resourcemanager.WithApplyPreProcessor(transactionPreprocessor.Preprocess),
			),
		).
		Register(testClient)
)

func resourceList() string {
	return strings.Join(resources.List(), "|")
}

func setupResources() {
	extraHeaders := http.Header{}
	extraHeaders.Set("x-client-id", analytics.ClientID())
	extraHeaders.Set("x-source", "cli")

	// To avoid a ciruclar reference initialization when setting up the registry and its resources,
	// we create the resources with a pointer to an unconfigured HTTPClient.
	// When each command is run, this function is run in the PreRun stage, before any of the actual `Run` code is executed
	// We take this chance to configure the HTTPClient with the correct URL and headers.
	// To make this configuration propagate to all the resources, we need to replace the pointer to the HTTPClient.
	// For more details, see https://github.com/kubeshop/tracetest/pull/2832#discussion_r1245616804
	hc := resourcemanager.NewHTTPClient(cliConfig.URL(), extraHeaders)
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
