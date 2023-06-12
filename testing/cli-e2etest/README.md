# Tracetest CLI e2e Tests

In this folder we have the End-to-end tests done on the CLI to guarantee that the CLI is working fine. 
The main idea is to test every CLI command against the Tracetest server with different data stores and different operating systems.

## Implementation Status

| Linux              | Windows | MacOS  |
| ------------------ | ------- | ------ |
| :white_check_mark: | :soon:  | :soon: |

## Tracetest Data Store

| Jaeger             | Tempo  | OpenSearch | SignalFx | OTLP   | ElasticAPM | New Relic | Lightstep | Datadog | AWS X-Ray | Honeycomb |
| ------------------ | ------ | ---------- | -------- | ------ | ---------- | --------- | --------- | ------- | --------- | --------- |
| :white_check_mark: | :soon: | :soon:     | :soon:   | :soon: | :soon:     | :soon:    | :soon:    | :soon:  | :soon:    | :soon:    |

## CLI Commands to Test

### Misc and Flags

| CLI Command    | Test scenarios                                    |
| -------------- | ------------------------------------------------- |
| `version`      | [VersionCommand](./testscenarios/version_test.go) |
| `help`         | [HelpCommand](./testscenarios/help_test.go)       |
| `--help`, `-h` | [HelpCommand](./testscenarios/help_test.go)       |
| `--config`     | All scenarios                                     |

### Tests and Test Runner

| CLI Command                                                        | Test scenarios |
| ------------------------------------------------------------------ | -------------- |
| `test list`                                                        | |
| `test run -d [test-definition]`                                    | [TestRunTestWithGrpcTrigger](./testscenarios/test/run_test_with_grpc_trigger_test.go) |
| `test run -d [test-definition] -e [environment-id]`                | [TestRunTestWithHttpTriggerAndEnvironmentFile](./testscenarios/test/run_test_with_http_trigger_and_environment_file_test.go) |
| `test run -d [test-definition] -e [environment-definition]`        | [TestRunTestWithHttpTriggerAndEnvironmentFile](./testscenarios/test/run_test_with_http_trigger_and_environment_file_test.go) |
| `test run -d [transaction-definition]`                             | |
| `test run -d [transaction-definition] -e [environment-id]`         | |
| `test run -d [transaction-definition] -e [environment-definition]` | |

### Resources: Config

| CLI Command                                           | Test scenarios |
| ----------------------------------------------------- | -------------- |
| `apply config -f [config-file]`                       | [ApplyConfig](./testscenarios/config/apply_config_test.go) |
| `delete config --id current`                          | [DeleteConfig](./testscenarios/config/delete_config_test.go) |
| `get config --id current --output pretty`             | [GetConfig](./testscenarios/config/get_config_test.go), [ApplyConfig](./testscenarios/config/apply_config_test.go), [DeleteConfig](./testscenarios/config/delete_config_test.go) |
| `get config --id current --output json`               | [GetConfig](./testscenarios/config/get_config_test.go) |
| `get config --id current --output yaml`               | [GetConfig](./testscenarios/config/get_config_test.go) |
| `list config --output pretty`                         | [ListConfig](./testscenarios/config/list_config_test.go) |
| `list config --output json`                           | [ListConfig](./testscenarios/config/list_config_test.go) |
| `list config --output yaml`                           | [ListConfig](./testscenarios/config/list_config_test.go) |

### Resources: Data Store

| CLI Command                                              | Test scenarios |
| -------------------------------------------------------- | -------------- |
| `apply datastore -f [data-store-file]`                   | [ApplyDatastore](./testscenarios/datastore/apply_datastore_test.go) |
| `delete datastore --id current`                          | [DeleteDatastore](./testscenarios/datastore/delete_datastore_test.go) |
| `get datastore --id current --output pretty`             | [GetDatastore](./testscenarios/datastore/get_datastore_test.go), [ApplyDatastore](./testscenarios/datastore/apply_datastore_test.go), [DeleteDatastore](./testscenarios/datastore/delete_datastore_test.go) |
| `get datastore --id current --output json`               | [GetDatastore](./testscenarios/datastore/get_datastore_test.go) |
| `get datastore --id current --output yaml`               | [GetDatastore](./testscenarios/datastore/get_datastore_test.go) |
| `list datastore --output pretty`                         | [ListDatastore](./testscenarios/datastore/list_datastore_test.go) |
| `list datastore --output json`                           | [ListDatastore](./testscenarios/datastore/list_datastore_test.go) |
| `list datastore --output yaml`                           | [ListDatastore](./testscenarios/datastore/list_datastore_test.go) |

### Resources: Demo

| CLI Command                                          | Test scenarios |
| ---------------------------------------------------- | -------------- |
| `apply demo -f [new-demo-file]`                      | |
| `apply demo -f [existing-demo-file]`                 | |
| `delete demo --id [existing-id]`                     | |
| `delete demo --id [non-existing-id]`                 | |
| `get demo --id [non-existing-id]`                    | |
| `get demo --id [existing-id] --output pretty`        | |
| `get demo --id [existing-id] --output json`          | |
| `get demo --id [existing-id] --output yaml`          | |
| `list demo --output pretty`                          | |
| `list demo --output json`                            | |
| `list demo --output yaml`                            | |
| `list demo --skip 1 --take 2`                        | |
| `list demo --sortBy name --sortDirection asc`        | |

### Resources: Environment

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply environment -f [new-environment-file]`               | [ApplyEnvironment](./testscenarios/environment/apply_environment_test.go) |
| `apply environment -f [existing-environment-file]`          | [ApplyEnvironment](./testscenarios/environment/apply_environment_test.go) |
| `delete environment --id [existing-id]`                     | [DeleteEnvironment](./testscenarios/environment/delete_environment_test.go) |
| `delete environment --id [non-existing-id]`                 | [DeleteEnvironment](./testscenarios/environment/delete_environment_test.go) |
| `get environment --id [non-existing-id]`                    | [GetEnvironment](./testscenarios/environment/get_environment_test.go), [DeleteEnvironment](./testscenarios/environment/delete_environment_test.go) |
| `get environment --id [existing-id] --output pretty`        | [GetEnvironment](./testscenarios/environment/get_environment_test.go) |
| `get environment --id [existing-id] --output json`          | [GetEnvironment](./testscenarios/environment/get_environment_test.go) |
| `get environment --id [existing-id] --output yaml`          | [GetEnvironment](./testscenarios/environment/get_environment_test.go), [ApplyEnvironment](./testscenarios/environment/apply_environment_test.go) |
| `list environment --output pretty`                          | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |
| `list environment --output json`                            | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |
| `list environment --output yaml`                            | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |
| `list environment --skip 1 --take 2`                        | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |
| `list environment --sortBy name --sortDirection asc`        | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |

### Resources: PollingProfile

| CLI Command                                                           | Test scenarios |
| --------------------------------------------------------------------- | -------------- |
| `apply pollingprofile -f [pollingprofile-file]`                       | [ApplyPollingProfile](./testscenarios/pollingprofile/apply_pollingprofile_test.go) |
| `delete pollingprofile --id current`                                  | [DeletePollingProfile](./testscenarios/pollingprofile/delete_pollingprofile_test.go) |
| `get pollingprofile --id current --output pretty`                     | [GetPollingProfile](./testscenarios/pollingprofile/get_pollingprofile_test.go), [ApplyPollingProfile](./testscenarios/pollingprofile/apply_pollingprofile_test.go), [DeletePollingProfile](./testscenarios/pollingprofile/delete_pollingprofile_test.go) |
| `get pollingprofile --id current --output json`                       | [GetPollingProfile](./testscenarios/pollingprofile/get_pollingprofile_test.go) |
| `get pollingprofile --id current --output yaml`                       | [GetPollingProfile](./testscenarios/pollingprofile/get_pollingprofile_test.go) |
| `list pollingprofile --output pretty`                                 | [ListPollingProfile](./testscenarios/pollingprofile/list_pollingprofile_test.go) |
| `list pollingprofile --output json`                                   | [ListPollingProfile](./testscenarios/pollingprofile/list_pollingprofile_test.go) |
| `list pollingprofile --output yaml`                                   | [ListPollingProfile](./testscenarios/pollingprofile/list_pollingprofile_test.go) |

### Resources: Transactions

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply transaction -f [new-transaction-file]`               | |
| `apply transaction -f [existing-transaction-file]`          | |
| `delete transaction --id [existing-id]`                     | |
| `delete transaction --id [non-existing-id]`                 | |
| `get transaction --id [non-existing-id]`                    | |
| `get transaction --id [existing-id] --output pretty`        | |
| `get transaction --id [existing-id] --output json`          | |
| `get transaction --id [existing-id] --output yaml`          | |
| `list transaction --output pretty`                          | |
| `list transaction --output json`                            | |
| `list transaction --output yaml`                            | |
| `list transaction --skip 1 --take 2`                        | |
| `list transaction --sortBy name --sortDirection asc`        | |
