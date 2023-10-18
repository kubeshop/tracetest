VERSION?="dev"
TAG?=$(VERSION)
GORELEASER_VERSION=1.21.2-pro

PROJECT_ROOT=${PWD}

CURRENT_GORELEASER_VERSION := $(shell goreleaser --version | head -n 9 | tail -n 1 |  tr -s ' ' | cut -d' ' -f2-)
goreleaser-version:
ifneq "$(CURRENT_GORELEASER_VERSION)" "$(GORELEASER_VERSION)"
	@printf "\033[0;31m Bad goreleaser version $(CURRENT_GORELEASER_VERSION), please install $(GORELEASER_VERSION)\033[0m\n\n"
	@printf "\033[0;31m Tracetest requires goreleaser pro installed (licence not necessary for local builds)\033[0m\n\n"
	@printf "\033[0;33m See https://goreleaser.com/install/ \033[0m\n\n"
	@exit 1
endif


CLI_SRC_FILES := $(shell find cli -type f)
dist/tracetest: goreleaser-version generate-cli $(CLI_SRC_FILES)
	goreleaser build --single-target --clean --snapshot --id cli
	find ./dist -name 'tracetest' -exec cp {} ./dist \;

SERVER_SRC_FILES := $(shell find server -type f)
dist/tracetest-server: goreleaser-version generate-server $(SERVER_SRC_FILES)
	goreleaser build --single-target --clean --snapshot --id server
	find ./dist -name 'tracetest-server' -exec cp {} ./dist \;

web/node_modules: web/package.json web/package-lock.json
	cd web; npm install

WEB_SRC_FILES := $(shell find web -type f -not -path "*node_modules*" -not -path "*build*" -not -path "*cypress/videos*" -not -path "*cypress/screenshots*")
web/build: web/node_modules $(WEB_SRC_FILES)
	cd web; npm run build

dist/tracetest-docker-$(TAG).tar dist/tracetest-agent-docker-$(TAG).tar: $(CLI_SRC_FILES) $(SERVER_SRC_FILES) $(WEB_SRC_FILES)
	goreleaser release --clean --skip-announce --snapshot -f .goreleaser.dev.yaml
	docker save --output dist/tracetest-docker-$(TAG).tar "kubeshop/tracetest:$(TAG)"
	docker save --output dist/tracetest-agent-docker-$(TAG).tar "kubeshop/tracetest-agent-:$(TAG)"

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

view-open-api: ## Run SwaggerUI locally to see OpenAPI documentation
	@echo "Running SwaggerUI..."
	@echo "Open http://localhost:9002 after the message 'Configuration complete; ready for start up'"
	@echo ""
	@docker run --rm -p 9002:8080 -v $(shell pwd)/api:/api -e SWAGGER_JSON=/api/openapi.yaml swaggerapi/swagger-ui

.PHONY: run build build-go build-web build-docker
run: build-docker ## build and run tracetest using docker compose
	docker compose up
build-go: dist/tracetest dist/tracetest-server ## build all go code
build-web: web/build ## build web
build-docker: goreleaser-version web/build .goreleaser.dev.yaml dist/tracetest-docker-$(TAG).tar dist/tracetest-agent-docker-$(TAG).tar ## build and tag docker image as defined in .goreleaser.dev.yaml

.PHONY: generate generate-server generate-cli generate-web
generate: generate-server generate-cli generate-web ## generate code entities from openapi definitions for all parts of the code
generate-server: server/openapi ## generate code entities from openapi definitions for server
generate-cli: cli/openapi ## generate code entities from openapi definitions for cli
generate-web: web/src/types/Generated.types.ts ## generate code entities from openapi definitions for web

OPENAPI_SRC_FILES := $(shell find api -type f)
OPENAPI_GENERATOR_VER=v6.3.0
OPENAPI_GENERATOR_IMAGE=openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VER)
OPENAPI_GENERATOR_CLI=docker run --rm -u ${shell id -u}  -v "$(PROJECT_ROOT):/local" -w "/local" ${OPENAPI_GENERATOR_IMAGE}
OPENAPI_TARGET_DIR=openapi/
web/src/types/Generated.types.ts: $(OPENAPI_SRC_FILES)
	cd web; npm run types:generate

cli/openapi: $(OPENAPI_SRC_FILES)
	$(eval BASE := ./cli)
	mkdir -p $(BASE)/tmp
	rm -rf $(BASE)/$(OPENAPI_TARGET_DIR)
	mkdir -p $(BASE)/$(OPENAPI_TARGET_DIR)

	$(OPENAPI_GENERATOR_CLI) generate \
		-i api/openapi.yaml \
		-g go \
		-o $(BASE)/tmp \
		--generate-alias-as-model
	cp $(BASE)/tmp/*.go $(BASE)/$(OPENAPI_TARGET_DIR)
	chmod 644 $(BASE)/$(OPENAPI_TARGET_DIR)/*.go
	rm -rf $(BASE)/tmp

	cd $(BASE); go fmt ./...

server/openapi: $(OPENAPI_SRC_FILES)
	$(eval BASE := ./server)
	mkdir -p $(BASE)/tmp
	rm -rf $(BASE)/$(OPENAPI_TARGET_DIR)
	mkdir -p $(BASE)/$(OPENAPI_TARGET_DIR)

	$(OPENAPI_GENERATOR_CLI) generate \
		-i api/openapi.yaml \
		-g go-server \
		-o $(BASE)/tmp \
		--generate-alias-as-model
	cp $(BASE)/tmp/go/*.go $(BASE)/$(OPENAPI_TARGET_DIR)
	chmod 644 $(BASE)/$(OPENAPI_TARGET_DIR)/*.go
	rm -f $(BASE)/$(OPENAPI_TARGET_DIR)/api_api_service.go
	rm -rf $(BASE)/tmp

	cd $(BASE); go fmt ./...


.PHONY: clean
clean: ## cleans the build artifacts
	rm -rf dist
	rm -rf web/build
	rm -rf web/node_modules
	docker image rm "kubeshop/tracetest:$(TAG)"
