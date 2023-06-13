TRACETEST_CLI?="../../dist/tracetest"
TEST_ENVIRONMENT?="jaeger"
TAG?="dev"
ENABLE_CLI_DEBUG?="false"

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

test: # run tests for this application
	go clean -testcache
	go test -timeout 300s -p 1 ./...

test/debug: # run tests for this application with debug mode enabled
	ENABLE_CLI_DEBUG="true"
	make test
