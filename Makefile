PROJECT_ROOT=${PWD}

CLI_SRC_FILES := $(shell find cli -type f)
dist/tracetest: generate-cli $(CLI_SRC_FILES)
	goreleaser build --single-target --clean --snapshot --id cli
	find ./dist -name 'tracetest' -exec cp {} ./dist \;

SERVER_SRC_FILES := $(shell find server -type f)
dist/tracetest-server: generate-server $(SERVER_SRC_FILES)
	goreleaser build --single-target --clean --snapshot --id server
	find ./dist -name 'tracetest-server' -exec cp {} ./dist \;

web/node_modules: web/package.json web/package-lock.json
	cd web; npm ci

WEB_SRC_FILES := $(shell find web -type f -not -path "*node_modules*" -not -path "*build*" -not -path "*cypress/videos*" -not -path "*cypress/screenshots*")
web/build: generate-web web/node_modules $(WEB_SRC_FILES)
	cd web; npm run build

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: run build build-go build-web build-docker
run: build-docker ## build and run tracetest using docker compose
	docker compose up
build-go: dist/tracetest dist/tracetest-server ## build all go code
build-web: web/build ## build web
build-docker: build-go build-web .goreleaser.dev.yaml ## build and tag docker image as defined in .goreleaser.dev.yaml
	VERSION=latest goreleaser release --clean --skip-announce --snapshot -f .goreleaser.dev.yaml

.PHONY: generate generate-server generate-cli generate-web
generate: generate-server generate-cli generate-web ## generate code entities from openapi definitions for all parts of the code
generate-server: $(wildcard server/openapi/*.go) ## generate code entities from openapi definitions for server
generate-cli: $(wildcard cli/openapi/*.go) ## generate code entities from openapi definitions for cli
generate-web: web/src/types/Generated.types.ts ## generate code entities from openapi definitions for web

OPENAPI_SRC_FILES := $(shell find api -type f)
OPENAPI_GENERATOR_VER=v6.3.0
OPENAPI_GENERATOR_IMAGE=openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VER)
OPENAPI_GENERATOR_CLI=docker run --rm -u ${shell id -u}  -v "$(PROJECT_ROOT):/local" -w "/local" ${OPENAPI_GENERATOR_IMAGE}
OPENAPI_TARGET_DIR=openapi/
web/src/types/Generated.types.ts: $(OPENAPI_SRC_FILES)
	cd web; npm run types:generate

$(wildcard cli/openapi/*.go): $(OPENAPI_SRC_FILES)
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

$(wildcard server/openapi/*.go): $(OPENAPI_SRC_FILES)
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

serve-docs: ## serve documentation for Tracetest
	docker build -t tracetest-docs -f docs-Dockerfile .
	docker run --network host tracetest-docs
	sleep 1
	open http://localhost:8000
