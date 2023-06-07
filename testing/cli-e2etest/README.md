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
| `apply config -f [config-file]`                       | |
| `delete config --id current`                          | |
| `get config --id current --output pretty`             | |
| `get config --id current --output json`               | |
| `get config --id current --output yaml`               | |
| `list config --output pretty`                         | |
| `list config --output json`                           | |
| `list config --output yaml`                           | |
### Resources: Data Store

| CLI Command                                              | Test scenarios |
| -------------------------------------------------------- | -------------- |
| `apply datastore -f [data-store-file]`                   | [ApplyNewDatastore](./testscenarios/datastore/apply_new_datastore_test.go) |
| `delete datastore --id current`                          | [DeleteDatastore](./testscenarios/datastore/delete_datastore_test.go) |
| `get datastore --id current --output pretty`             | [ApplyNewDatastore](./testscenarios/datastore/apply_new_datastore_test.go), [DeleteDatastore](./testscenarios/datastore/delete_datastore_test.go) |
| `get datastore --id current --output json`               | |
| `get datastore --id current --output yaml`               | |
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
| `apply environment -f [new-environment-file]`               | [ApplyNewEnvironment](./testscenarios/environment/apply_new_environment_test.go) |
| `apply environment -f [existing-environment-file]`          | [ApplyNewEnvironment](./testscenarios/environment/apply_new_environment_test.go) |
| `delete environment --id [existing-id]`                     | [DeleteEnvironment](./testscenarios/environment/delete_environment_test.go) |
| `delete environment --id [non-existing-id]`                 | |
| `get environment --id [non-existing-id]`                    | [DeleteEnvironment](./testscenarios/environment/delete_environment_test.go) |
| `get environment --id [existing-id] --output pretty`        | |
| `get environment --id [existing-id] --output json`          | |
| `get environment --id [existing-id] --output yaml`          | [ApplyNewEnvironment](./testscenarios/environment/apply_new_environment_test.go) |
| `list environment --output pretty`                          | |
| `list environment --output json`                            | |
| `list environment --output yaml`                            | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |
| `list environment --skip 1 --take 2`                        | |
| `list environment --sortBy name --sortDirection asc`        | [ListEnvironment](./testscenarios/environment/list_environments_test.go) |

### Resources: PollingProfile

| CLI Command                                                           | Test scenarios |
| --------------------------------------------------------------------- | -------------- |
| `apply pollingprofile -f [pollingprofile-file]`                       | |
| `delete pollingprofile --id current`                                  | |
| `export pollingprofile --id current --file [pollingprofile-file]`     | |
| `get pollingprofile --id current --output pretty`                     | |
| `get pollingprofile --id current --output json`                       | |
| `get pollingprofile --id current --output yaml`                       | |
| `list pollingprofile --output pretty`                                 | |
| `list pollingprofile --output json`                                   | |
| `list pollingprofile --output yaml`                                   | |

### Resources: Transactions

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply transaction -f [new-transaction-file]`               | |
| `apply transaction -f [existing-transaction-file]`          | |
| `delete transaction --id [existing-id]`                     | |
| `delete transaction --id [non-existing-id]`                 | |
| `export transaction --id current --file [transaction-file]` | |
| `get transaction --id [non-existing-id]`                    | |
| `get transaction --id [existing-id] --output pretty`        | |
| `get transaction --id [existing-id] --output json`          | |
| `get transaction --id [existing-id] --output yaml`          | |
| `list transaction --output pretty`                          | |
| `list transaction --output json`                            | |
| `list transaction --output yaml`                            | |
| `list transaction --skip 1 --take 2`                        | |
| `list transaction --sortBy name --sortDirection asc`        | |
