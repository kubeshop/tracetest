TRACETEST_COMMAND?="../../dist/tracetest"
TEST_ENVIRONMENT?="jaeger"
TAG?="dev"

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

test: # run tests for this application
	go clean -testcache
	go test -v -timeout 150s -p 1 ./...
