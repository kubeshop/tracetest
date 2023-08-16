# Tracetest CLI e2e Tests

In this folder we have the End-to-end tests done on the CLI to guarantee that the CLI is working fine. 
The main idea is to test every CLI command against the Tracetest server with different data stores and different operating systems.

## Implementation Status

| Linux              | Windows | MacOS  |
| ------------------ | ------- | ------ |
| :white_check_mark: | :soon:  | :soon: |

## Tracetest Data Store

| Jaeger             | Tempo  | OpenSearch | SignalFx | OTLP   | ElasticAPM | New Relic | Lightstep | Datadog | AWS X-Ray | Honeycomb | Dynatrace |
| ------------------ | ------ | ---------- | -------- | ------ | ---------- | --------- | --------- | ------- | --------- | --------- | --------- |
| :white_check_mark: | :soon: | :soon:     | :soon:   | :soon: | :soon:     | :soon:    | :soon:    | :soon:  | :soon:    | :soon:    | :soon:    |

## CLI Commands to Test

### Misc and Flags

| CLI Command    | Test scenarios                                    |
| -------------- | ------------------------------------------------- |
| `version`      | [VersionCommand](./testscenarios/version_test.go) |
| `help`         | [HelpCommand](./testscenarios/help_test.go)       |
| `--help`, `-h` | [HelpCommand](./testscenarios/help_test.go)       |
| `--config`     | All scenarios                                     |

### Run Tests and TestSuites

| CLI Command                                                        | Test scenarios |
| ------------------------------------------------------------------ | -------------- |
| `run test -f [test-definition]`                                    | [RunTestWithGrpcTrigger](./testscenarios/test/run_test_with_grpc_trigger_test.go) |
| `run test -f [test-definition] --vars [variableset-id]`                | [RunTestWithHttpTriggerAndVariableSetFile](./testscenarios/test/run_test_with_http_trigger_and_variableset_file_test.go) |
| `run test -f [test-definition] --vars [variableset-definition]`        | [RunTestWithHttpTriggerAndVariableSetFile](./testscenarios/test/run_test_with_http_trigger_and_variableset_file_test.go) |
| `run testsuite -f [testsuite-definition]`                             | [RunTestSuite](./testscenarios/testsuite//run_testsuite_test.go) |
| `run testsuite -f [testsuite-definition] --vars [variableset-id]`         | |
| `run testsuite -f [testsuite-definition] --vars [variableset-definition]` | |

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

### Resources: VariableSet

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply variableset -f [new-variableset-file]`               | [ApplyVariableSet](./testscenarios/variableset/apply_variableset_test.go) |
| `apply variableset -f [existing-variableset-file]`          | [ApplyVariableSet](./testscenarios/variableset/apply_variableset_test.go) |
| `delete variableset --id [existing-id]`                     | [DeleteVariableSet](./testscenarios/variableset/delete_variableset_test.go) |
| `delete variableset --id [non-existing-id]`                 | [DeleteVariableSet](./testscenarios/variableset/delete_variableset_test.go) |
| `get variableset --id [non-existing-id]`                    | [GetVariableSet](./testscenarios/variableset/get_variableset_test.go), [DeleteVariableSet](./testscenarios/variableset/delete_variableset_test.go) |
| `get variableset --id [existing-id] --output pretty`        | [GetVariableSet](./testscenarios/variableset/get_variableset_test.go) |
| `get variableset --id [existing-id] --output json`          | [GetVariableSet](./testscenarios/variableset/get_variableset_test.go) |
| `get variableset --id [existing-id] --output yaml`          | [GetVariableSet](./testscenarios/variableset/get_variableset_test.go) |
| `list variableset --output pretty`                          | [ListVariableSet](./testscenarios/variableset/list_variableset_test.go) |
| `list variableset --output json`                            | [ListVariableSet](./testscenarios/variableset/list_variableset_test.go) |
| `list variableset --output yaml`                            | [ListVariableSet](./testscenarios/variableset/list_variableset_test.go) |
| `list variableset --skip 1 --take 2`                        | [ListVariableSet](./testscenarios/variableset/list_variableset_test.go) |
| `list variableset --sortBy name --sortDirection asc`        | [ListVariableSet](./testscenarios/variableset/list_variableset_test.go) |

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

### Resources: TestRunner

| CLI Command                                                           | Test scenarios |
| --------------------------------------------------------------------- | -------------- |
| `apply testrunner -f [testrunner-file]`                               | [ApplyTestRunner](./testscenarios/testrunner/apply_testrunner_test.go) |
| `delete testrunner --id current`                                      | [DeleteTestRunner](./testscenarios/testrunner/delete_testrunner_test.go) |
| `get testrunner --id current --output pretty`                         | [GetTestRunner](./testscenarios/testrunner/get_testrunner_test.go) |
| `get testrunner --id current --output json`                           | [GetTestRunner](./testscenarios/testrunner/get_testrunner_test.go) |
| `get testrunner --id current --output yaml`                           | [GetTestRunner](./testscenarios/testrunner/get_testrunner_test.go) |
| `list testrunner --output pretty`                                     | [ListTestRunner](./testscenarios/testrunner/list_testrunner_test.go) |
| `list testrunner --output json`                                       | [ListTestRunner](./testscenarios/testrunner/list_testrunner_test.go) |
| `list testrunner --output yaml`                                       | [ListTestRunner](./testscenarios/testrunner/list_testrunner_test.go) |

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

### Resources: TestSuites

| CLI Command                                                 | Test scenarios |
| ----------------------------------------------------------- | -------------- |
| `apply testsuite -f [new-testsuite-file]`               | [ApplyTestSuite](./testscenarios/testsuite/apply_testsuite_test.go) |
| `apply testsuite -f [existing-testsuite-file]`          | [ApplyTestSuite](./testscenarios/testsuite/apply_testsuite_test.go) |
| `delete testsuite --id [existing-id]`                     | [DeleteTestSuite](./testscenarios/testsuite/delete_testsuite_test.go) |
| `delete testsuite --id [non-existing-id]`                 | [DeleteTestSuite](./testscenarios/testsuite/delete_testsuite_test.go) |
| `get testsuite --id [non-existing-id]`                    | [GetTestSuite](./testscenarios/testsuite/get_testsuite_test.go), [DeleteTestSuite](./testscenarios/testsuite/delete_testsuite_test.go) |
| `get testsuite --id [existing-id] --output pretty`        | [GetTestSuite](./testscenarios/testsuite/get_testsuite_test.go) |
| `get testsuite --id [existing-id] --output json`          | [GetTestSuite](./testscenarios/testsuite/get_testsuite_test.go) |
| `get testsuite --id [existing-id] --output yaml`          | [GetTestSuite](./testscenarios/testsuite/get_testsuite_test.go) |
| `list testsuite --output pretty`                          | [ListTestSuite](./testscenarios/testsuite/list_testsuites_test.go) |
| `list testsuite --output json`                            | [ListTestSuite](./testscenarios/testsuite/list_testsuites_test.go) |
| `list testsuite --output yaml`                            | [ListTestSuite](./testscenarios/testsuite/list_testsuites_test.go) |
| `list testsuite --skip 1 --take 2`                        | [ListTestSuite](./testscenarios/testsuite/list_testsuites_test.go) |
| `list testsuite --sortBy name --sortDirection asc`        | [ListTestSuite](./testscenarios/testsuite/list_testsuites_test.go) |

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
