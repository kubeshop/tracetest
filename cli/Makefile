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

test:
	@go test -coverprofile=coverage.out ./...
