# Tracetest CLI e2e tests

In this folder we have the End-to-end tests done on the CLI to guarantee that the CLI is working fine. 
The main idea is to test every CLI command against Tracetest server with different data stores and different Operating systems.


## Implementation status

| Linux              | Windows | MacOS  |
| ------------------ | ------- | ------ |
| :white_check_mark: | :soon:  | :soon: |

## Tracetest Data Store

| Jaeger             | Tempo  | OpenSearch | SignalFx | OTLP   | ElasticAPM | New Relic | Lightstep | Datadog | AWS X-Ray | Honeycomb |
| ------------------ | ------ | ---------- | -------- | ------ | ---------- | --------- | --------- | ------- | --------- | --------- |
| :white_check_mark: | :soon: | :soon:     | :soon:   | :soon: | :soon:     | :soon:    | :soon:    | :soon:  | :soon:    | :soon:    |

## CLI Commands to Test

### Misc and Flags

| CLI Command    | Tested             | Test scenarios |
| -------------- | ------------------ | -------------- |
| `version`      | :white_check_mark: | [VersionCommand](./testscenarios/version_test.go) |
| `help`         | :white_check_mark: | [HelpCommand](./testscenarios/help_test.go) |
| `--help`, `-h` | :white_check_mark: | [HelpCommand](./testscenarios/help_test.go) |
| `--config`     | :white_check_mark: | All scenarios |

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

### Tests

| CLI Command                                                 | Tested             | Test scenarios |
| ----------------------------------------------------------- | ------------------ | -------------- |
| `test list`                                                 | :yellow_circle:    | |
| `test run -d [test-definition]`                             | :yellow_circle:    | |
| `test run -d [test-definition] -e [environment-id]`         | :yellow_circle:    | |
| `test run -d [test-definition] -e [environment-definition]` | :yellow_circle:    | |
