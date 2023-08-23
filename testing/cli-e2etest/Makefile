PROJECT_ROOT=${PWD}
TRACETEST_CLI?="${PROJECT_ROOT}/../../dist/tracetest"
TEST_ENVIRONMENT?="jaeger"
TAG?="dev"
ENABLE_CLI_DEBUG?="false"
TEST_SCENARIO?=test

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

test: ## run tests for this application
	export TRACETEST_CLI=$(TRACETEST_CLI); \
	export TEST_ENVIRONMENT=$(TEST_ENVIRONMENT); \
	export TAG=$(TAG); \
	export ENABLE_CLI_DEBUG=$(ENABLE_CLI_DEBUG); \
	go clean -testcache; \
	go test -v -timeout 300s -p 1 ./...;

test/scenario: ## run tests for this application
	export TRACETEST_CLI=$(TRACETEST_CLI); \
	export TEST_ENVIRONMENT=$(TEST_ENVIRONMENT); \
	export TAG=$(TAG); \
	export ENABLE_CLI_DEBUG=$(ENABLE_CLI_DEBUG); \
	go clean -testcache; \
	go test -v -timeout 300s -p 1 "$(PROJECT_ROOT)/testscenarios/$(TEST_SCENARIO)";

test/debug: ## run tests for this application with debug mode enabled
	export ENABLE_CLI_DEBUG="true"; \
	make test;
