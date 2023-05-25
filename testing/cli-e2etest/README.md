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

| CLI Command    | Tested             | Test scenarios                                    |
| -------------- | ------------------ | ------------------------------------------------- |
| `version`      | :white_check_mark: | [VersionCommand](./testscenarios/version_test.go) |
| `help`         | :white_check_mark: | [HelpCommand](./testscenarios/help_test.go)       |
| `--help`, `-h` | :white_check_mark: | [HelpCommand](./testscenarios/help_test.go)       |
| `--config`     | :white_check_mark: | All scenarios                                     |

### Tests and Test Runner

| CLI Command                                                        | Tested             | Test scenarios |
| ------------------------------------------------------------------ | ------------------ | -------------- |
| `test list`                                                        | :yellow_circle:    | |
| `test run -d [test-definition]`                                    | :yellow_circle:    | |
| `test run -d [test-definition] -e [environment-id]`                | :yellow_circle:    | |
| `test run -d [test-definition] -e [environment-definition]`        | :yellow_circle:    | |
| `test run -d [transaction-definition]`                             | :yellow_circle:    | |
| `test run -d [transaction-definition] -e [environment-id]`         | :yellow_circle:    | |
| `test run -d [transaction-definition] -e [environment-definition]` | :yellow_circle:    | |

### Resources: Config

| CLI Command                                           | Tested          | Test scenarios |
| ----------------------------------------------------- | ----------------| -------------- |
| `apply config -f [config-file]`                       | :yellow_circle: | |
| `delete config --id current`                          | :yellow_circle: | |
| `export config --id current --file [config-file]`     | :yellow_circle: | |
| `get config --id current --output pretty`             | :yellow_circle: | |
| `get config --id current --output json`               | :yellow_circle: | |
| `get config --id current --output yaml`               | :yellow_circle: | |
| `list config --output pretty`                         | :yellow_circle: | |
| `list config --output json`                           | :yellow_circle: | |
| `list config --output yaml`                           | :yellow_circle: | |
### Resources: Data Store

| CLI Command                                              | Tested             | Test scenarios |
| -------------------------------------------------------- | ------------------ | -------------- |
| `apply datastore -f [data-store-file]`                   | :white_check_mark: | [ApplyNewDatastore](./testscenarios/datastore/apply_new_datastore_test.go) |
| `delete datastore --id current`                          | :white_check_mark: | [DeleteDatastore](./testscenarios/datastore/delete_datastore_test.go) |
| `export datastore --id current --file [data-store-file]` | :yellow_circle:    | |
| `get datastore --id current --output pretty`             | :white_check_mark: | [ApplyNewDatastore](./testscenarios/datastore/apply_new_datastore_test.go), [DeleteDatastore](./testscenarios/datastore/delete_datastore_test.go) |
| `get datastore --id current --output json`               | :yellow_circle:    | |
| `get datastore --id current --output yaml`               | :yellow_circle:    | |
| `list datastore --output pretty`                         | :white_check_mark: | [ListDatastore](./testscenarios/datastore/list_datastore_test.go) |
| `list datastore --output json`                           | :white_check_mark: | [ListDatastore](./testscenarios/datastore/list_datastore_test.go) |
| `list datastore --output yaml`                           | :white_check_mark: | [ListDatastore](./testscenarios/datastore/list_datastore_test.go) |

### Resources: Demo

| CLI Command                                          | Tested             | Test scenarios |
| ---------------------------------------------------- | ------------------ | -------------- |
| `apply demo -f [new-demo-file]`                      | :yellow_circle:    | |
| `apply demo -f [existing-demo-file]`                 | :yellow_circle:    | |
| `delete demo --id [existing-id]`                     | :yellow_circle:    | |
| `delete demo --id [non-existing-id]`                 | :yellow_circle:    | |
| `export demo --id current --file [demo-file]`        | :yellow_circle:    | |
| `get demo --id [non-existing-id]`                    | :yellow_circle:    | |
| `get demo --id [existing-id] --output pretty`        | :yellow_circle:    | |
| `get demo --id [existing-id] --output json`          | :yellow_circle:    | |
| `get demo --id [existing-id] --output yaml`          | :yellow_circle:    | |
| `list demo --output pretty`                          | :yellow_circle:    | |
| `list demo --output json`                            | :yellow_circle:    | |
| `list demo --output yaml`                            | :yellow_circle:    | |
| `list demo --skip 1 --take 2`                        | :yellow_circle:    | |
| `list demo --sortBy name --sortDirection desc`       | :yellow_circle:    | |

### Resources: Environment

| CLI Command                                                 | Tested             | Test scenarios |
| ----------------------------------------------------------- | ------------------ | -------------- |
| `apply environment -f [new-environment-file]`               | :yellow_circle:    | |
| `apply environment -f [existing-environment-file]`          | :yellow_circle:    | |
| `delete environment --id [existing-id]`                     | :yellow_circle:    | |
| `delete environment --id [non-existing-id]`                 | :yellow_circle:    | |
| `export environment --id current --file [environment-file]` | :yellow_circle:    | |
| `get environment --id [non-existing-id]`                    | :yellow_circle:    | |
| `get environment --id [existing-id] --output pretty`        | :yellow_circle:    | |
| `get environment --id [existing-id] --output json`          | :yellow_circle:    | |
| `get environment --id [existing-id] --output yaml`          | :yellow_circle:    | |
| `list environment --output pretty`                          | :yellow_circle:    | |
| `list environment --output json`                            | :yellow_circle:    | |
| `list environment --output yaml`                            | :yellow_circle:    | |
| `list environment --skip 1 --take 2`                        | :yellow_circle:    | |
| `list environment --sortBy name --sortDirection desc`       | :yellow_circle:    | |

### Resources: PollingProfile

| CLI Command                                                           | Tested          | Test scenarios |
| --------------------------------------------------------------------- | --------------- | -------------- |
| `apply pollingprofile -f [pollingprofile-file]`                       | :yellow_circle: | |
| `delete pollingprofile --id current`                                  | :yellow_circle: | |
| `export pollingprofile --id current --file [pollingprofile-file]`     | :yellow_circle: | |
| `get pollingprofile --id current --output pretty`                     | :yellow_circle: | |
| `get pollingprofile --id current --output json`                       | :yellow_circle: | |
| `get pollingprofile --id current --output yaml`                       | :yellow_circle: | |
| `list pollingprofile --output pretty`                                 | :yellow_circle: | |
| `list pollingprofile --output json`                                   | :yellow_circle: | |
| `list pollingprofile --output yaml`                                   | :yellow_circle: | |

### Resources: Transactions

| CLI Command                                                 | Tested             | Test scenarios |
| ----------------------------------------------------------- | ------------------ | -------------- |
| `apply transaction -f [new-transaction-file]`               | :yellow_circle:    | |
| `apply transaction -f [existing-transaction-file]`          | :yellow_circle:    | |
| `delete transaction --id [existing-id]`                     | :yellow_circle:    | |
| `delete transaction --id [non-existing-id]`                 | :yellow_circle:    | |
| `export transaction --id current --file [transaction-file]` | :yellow_circle:    | |
| `get transaction --id [non-existing-id]`                    | :yellow_circle:    | |
| `get transaction --id [existing-id] --output pretty`        | :yellow_circle:    | |
| `get transaction --id [existing-id] --output json`          | :yellow_circle:    | |
| `get transaction --id [existing-id] --output yaml`          | :yellow_circle:    | |
| `list transaction --output pretty`                          | :yellow_circle:    | |
| `list transaction --output json`                            | :yellow_circle:    | |
| `list transaction --output yaml`                            | :yellow_circle:    | |
| `list transaction --skip 1 --take 2`                        | :yellow_circle:    | |
| `list transaction --sortBy name --sortDirection desc`       | :yellow_circle:    | |
