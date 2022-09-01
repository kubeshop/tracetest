VERSION ?= "dev"
TRACETEST_ENV ?= "dev"
ANALYTICS_BE_KEY ?= ""


PROJECT_ROOT=${shell dirname $${PWD}}

OPENAPI_GENERATOR_VER=v5.4.0
OPENAPI_GENERATOR_IMAGE=openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VER)
OPENAPI_GENERATOR_CLI=docker run --rm -u ${shell id -u}  -v "$(PROJECT_ROOT):/local" -w "/local" ${OPENAPI_GENERATOR_IMAGE}
OPENAPI_TARGET_DIR=./openapi/

generate-client:
	rm -rf $(OPENAPI_TARGET_DIR)
	mkdir -p ./tmp
	mkdir -p $(OPENAPI_TARGET_DIR)

	$(OPENAPI_GENERATOR_CLI)  generate \
		-i api/openapi.yaml \
		-g go \
		-o ./cli/tmp \
		--generate-alias-as-model
	cp ./tmp/*.go $(OPENAPI_TARGET_DIR)
	rm -rf ./tmp

	go fmt ./...; cd ..

build:
ifeq (, $(shell which goreleaser))
	go install github.com/goreleaser/goreleaser@latest
endif
	VERSION=$(VERSION) \
	TRACETEST_ENV=$(TRACETEST_ENV) \
	ANALYTICS_BE_KEY=$(ANALYTICS_BE_KEY) \
	goreleaser build --single-target --rm-dist --snapshot -f ../.goreleaser.yaml
	cp `cat ./dist/artifacts.json | jq -rc '.[0].path'` ./dist/tracetest

simple-build:
	go build -o tracetest main.go

test: mockserver
	@go test -coverprofile=coverage.out ./...
	@make --no-print-directory kill-mock

mockserver:
ifeq (, $(shell which prism))
	npm install -g @stoplight/prism-cli@4.8.0
endif
	@make --no-print-directory kill-mock
	@echo "Starting prism mock server"
	@prism mock ../api/openapi.yaml &
	@sleep 2

kill-mock:
	@echo "Stopping prism mock server"
	@ps -ef | grep /bin/prism | grep -v grep | awk '{print $$2}' | xargs -r kill -9
