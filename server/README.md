# Server architecture and coding guidelines

## Basic code organization rules

### 1. Resources

- `resource_name` must be singular. **snake_case** for filenames, **CamelCaseUpperFirst** for struct names.
- Resource package names must be in **singular**.
- Resource files must be named `{resource_name}_entities.go` and  `{resource_name}_repository.go`
  > This format makes it easy to fuzzy search files
- `{resource_name}_entities.go` must contain only one "main" entity and zero or more sub-entities (i.e. [demoresource](https://github.com/kubeshop/tracetest/blob/main/server/config/demoresource/demo_resource.go))
- Resource packages can have one or  more "main" entities, each living on its own file. For example, the `transactions` package could have `transaction_entities.go` and `transaction_runs_entities.go`
- Resource packages can have "subpackages". All the same rules apply. The package must **NOT** be prefixed with the parent package (i.e. `test/trigger` is good, `test/testtrigger` is bad), with the exception of generic names (i.e. `test/util` should be `test/testutil`)

Example tests resource package:

```
test/
	trigger/
		trigger_entities.go	
		http_entities.go

	test_repository.go
	test_entities.go
	test_run_repository.go
	test_run_entities.go
```
