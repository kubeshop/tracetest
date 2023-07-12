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
| `test run -d [test-definition]`                                    | [RunTestWithGrpcTrigger](./testscenarios/test/run_test_with_grpc_trigger_test.go) |
| `test run -d [test-definition] -e [environment-id]`                | [RunTestWithHttpTriggerAndEnvironmentFile](./testscenarios/test/run_test_with_http_trigger_and_environment_file_test.go) |
| `test run -d [test-definition] -e [environment-definition]`        | [RunTestWithHttpTriggerAndEnvironmentFile](./testscenarios/test/run_test_with_http_trigger_and_environment_file_test.go) |
| `test run -d [transaction-definition]`                             | [RunTransaction](./testscenarios/transaction//run_transaction_test.go) |
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
| `apply demo -f [new-demo-file]`                      | [ApplyDemo](./testscenarios/demo/apply_demo_test.go) |
| `apply demo -f [existing-demo-file]`                 | [ApplyDemo](./testscenarios/demo/apply_demo_test.go) |
| `delete demo --id [existing-id]`                     | [DeleteDemo](./testscenarios/demo/delete_demo_test.go) |
| `delete demo --id [non-existing-id]`                 | [DeleteDemo](./testscenarios/demo/delete_demo_test.go) |
| `get demo --id [non-existing-id]`                    | [GetDemo](./testscenarios/demo/get_demo_test.go), [DeleteDemo](./testscenarios/demo/delete_demo_test.go) |
| `get demo --id [existing-id] --output pretty`        | [GetDemo](./testscenarios/demo/get_demo_test.go) |
| `get demo --id [existing-id] --output json`          | [GetDemo](./testscenarios/demo/get_demo_test.go) |
| `get demo --id [existing-id] --output yaml`          | [GetDemo](./testscenarios/demo/get_demo_test.go) |
| `list demo --output pretty`                          | [ListDemo](./testscenarios/demo/list_demos_test.go) |
| `list demo --output json`                            | [ListDemo](./testscenarios/demo/list_demos_test.go) |
| `list demo --output yaml`                            | [ListDemo](./testscenarios/demo/list_demos_test.go) |
| `list demo --skip 1 --take 1`                        | [ListDemo](./testscenarios/demo/list_demos_test.go) |
| `list demo --sortBy name --sortDirection asc`        | [ListDemo](./testscenarios/demo/list_demos_test.go) |

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
| `get environment --id [existing-id] --output yaml`          | [GetEnvironment](./testscenarios/environment/get_environment_test.go) |
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

### Resources: Analyzer

| CLI Command                                           | Test scenarios |
| ----------------------------------------------------- | -------------- |
| `apply analyzer -f [analyzer-file]`                     | [ApplyAnalyzer](./testscenarios/analyzer/apply_analyzer_test.go) |
| `delete analyzer --id current`                          | [DeleteAnalyzer](./testscenarios/analyzer/delete_analyzer_test.go) |
| `get analyzer --id current --output pretty`             | [GetAnalyzer](./testscenarios/analyzer/get_analyzer_test.go) |
| `get analyzer --id current --output json`               | [GetAnalyzer](./testscenarios/analyzer/get_analyzer_test.go) |
| `get analyzer --id current --output yaml`               | [GetAnalyzer](./testscenarios/analyzer/get_analyzer_test.go) |
| `list analyzer --output pretty`                         | [ListAnalyzer](./testscenarios/analyzer/list_analyzer_test.go) |
| `list analyzer --output json`                           | [ListAnalyzer](./testscenarios/analyzer/list_analyzer_test.go) |
| `list analyzer --output yaml`                           | [ListAnalyzer](./testscenarios/analyzer/list_analyzer_test.go) |

### Resources: Transactions

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply transaction -f [new-transaction-file]`               | [ApplyTransaction](./testscenarios/transaction/apply_transaction_test.go) |
| `apply transaction -f [existing-transaction-file]`          | [ApplyTransaction](./testscenarios/transaction/apply_transaction_test.go) |
| `delete transaction --id [existing-id]`                     | [DeleteTransaction](./testscenarios/transaction/delete_transaction_test.go) |
| `delete transaction --id [non-existing-id]`                 | [DeleteTransaction](./testscenarios/transaction/delete_transaction_test.go) |
| `get transaction --id [non-existing-id]`                    | [GetTransaction](./testscenarios/transaction/get_transaction_test.go), [DeleteTransaction](./testscenarios/transaction/delete_transaction_test.go) |
| `get transaction --id [existing-id] --output pretty`        | [GetTransaction](./testscenarios/transaction/get_transaction_test.go) |
| `get transaction --id [existing-id] --output json`          | [GetTransaction](./testscenarios/transaction/get_transaction_test.go) |
| `get transaction --id [existing-id] --output yaml`          | [GetTransaction](./testscenarios/transaction/get_transaction_test.go) |
| `list transaction --output pretty`                          | [ListTransaction](./testscenarios/transaction/list_transactions_test.go) |
| `list transaction --output json`                            | [ListTransaction](./testscenarios/transaction/list_transactions_test.go) |
| `list transaction --output yaml`                            | [ListTransaction](./testscenarios/transaction/list_transactions_test.go) |
| `list transaction --skip 1 --take 2`                        | [ListTransaction](./testscenarios/transaction/list_transactions_test.go) |
| `list transaction --sortBy name --sortDirection asc`        | [ListTransaction](./testscenarios/transaction/list_transactions_test.go) |

### Resources: Tests

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply test -f [new-test-file]`                             | [ApplyTest](./testscenarios/test/apply_test_test.go) |
| `apply test -f [existing-test-file]`                        | [ApplyTest](./testscenarios/test/apply_test_test.go) |
| `delete test --id [existing-id]`                            | [DeleteTest](./testscenarios/test/delete_test_test.go) |
| `delete test --id [non-existing-id]`                        | [DeleteTest](./testscenarios/test/delete_test_test.go) |
| `get test --id [non-existing-id]`                           | [GetTest](./testscenarios/test/get_test_test.go), [DeleteTest](./testscenarios/test/delete_test_test.go) |
| `get test --id [existing-id] --output pretty`               | [GetTest](./testscenarios/test/get_test_test.go) |
| `get test --id [existing-id] --output json`                 | [GetTest](./testscenarios/test/get_test_test.go) |
| `get test --id [existing-id] --output yaml`                 | [GetTest](./testscenarios/test/get_test_test.go) |
| `list test --output pretty`                                 | [ListTest](./testscenarios/test/list_test_test.go) |
| `list test --output json`                                   | [ListTest](./testscenarios/test/list_test_test.go) |
| `list test --output yaml`                                   | [ListTest](./testscenarios/test/list_test_test.go) |
| `list test --skip 1 --take 2`                               | [ListTest](./testscenarios/test/list_test_test.go) |
| `list test --sortBy name --sortDirection asc`               | [ListTest](./testscenarios/test/list_test_test.go) |
