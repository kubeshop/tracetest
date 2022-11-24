VERSION ?= "dev"
TRACETEST_ENV ?= "dev"
ANALYTICS_BE_KEY ?= ""

build:
ifeq (, $(shell which goreleaser))
	go install github.com/goreleaser/goreleaser@latest
endif
	VERSION=$(VERSION) \
	TRACETEST_ENV=$(TRACETEST_ENV) \
	ANALYTICS_BE_KEY=$(ANALYTICS_BE_KEY) \
	goreleaser build --single-target --rm-dist --snapshot -f ../.goreleaser.yaml
	cp `cat ./dist/artifacts.json | jq -rc '.[0].path'` ./dist/tracetest

test:
	@go test -coverprofile=coverage.out ./...
