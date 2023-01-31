VERSION ?= "dev"
TRACETEST_ENV ?= "dev"
ANALYTICS_BE_KEY ?= ""

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

build: ## build this project using goreleaser
ifeq (, $(shell which goreleaser))
	go install github.com/goreleaser/goreleaser@v1.14.1
endif
	VERSION=$(VERSION) \
	TRACETEST_ENV=$(TRACETEST_ENV) \
	ANALYTICS_BE_KEY=$(ANALYTICS_BE_KEY) \
	goreleaser build --single-target --rm-dist --snapshot -f ../.goreleaser.yaml
	cp `cat ./dist/artifacts.json | jq -rc '.[0].path'` ./dist/tracetest

test: ## execute unit tests
	@go test -coverprofile=coverage.out ./...
