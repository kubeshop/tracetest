####################
## Configs        ##
####################
GORELEASER_VERSION=1.15.2-pro
OPENAPI_GENERATOR_VERSION=v6.3.0

####################
## Sources        ##
####################
DOCKER_TARGET := $(wildcard dist/goreleaserdocker*/Dockerfile)

OPENAPI_SRC_FILES := $(shell find api -type f)

WEB_SRC_FILES := $(shell find web -type f -not -path "*node_modules*" -not -path "*build*" -not -path "*cypress/videos*" -not -path "*cypress/screenshots*")

SERVER_SRC_FILES := $(shell find server -type f)
SERVER_OPENAPI_FILES := $(shell find server/openapi -type f -name '*.go')

CLI_SRC_FILES := $(shell find cli -type f)
CLI_OPENAPI_FILES := $(wildcard cli/openapi/*.go)

####################
## Helpers        ##
####################
help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

CURRENT_GORELEASER_VERSION := $(shell goreleaser --version | head -n 1 | cut -d' ' -f3-)
goreleaser-version: ## check if goreleaser version is ok
ifneq "$(CURRENT_GORELEASER_VERSION)" "$(GORELEASER_VERSION)"
	@printf "\033[0;31m Bad goreleaser version $(CURRENT_GORELEASER_VERSION), please install $(GORELEASER_VERSION)\033[0m\n\n"
	@printf "\033[0;31m Tracetest requires goreleaser pro installed (licence not necessary for local builds)\033[0m\n\n"
	@printf "\033[0;33m See https://goreleaser.com/install/ \033[0m\n\n"
	@exit 1
endif

####################
## ACTUAL TARGETS ##
####################

PROJECT_ROOT=${PWD}

OPENAPI_GENERATOR_IMAGE=openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VERSION)
OPENAPI_GENERATOR_CLI=docker run --rm -u ${shell id -u}  -v "$(PROJECT_ROOT):/local" -w "/local" ${OPENAPI_GENERATOR_IMAGE}
OPENAPI_TARGET_DIR=openapi/

## Generated cli openapi files
$(CLI_OPENAPI_FILES): $(OPENAPI_SRC_FILES)
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

dist/tracetest:$(CLI_SRC_FILES)
	goreleaser build --single-target --clean --snapshot --id cli
	find ./dist -name 'tracetest' -exec cp {} ./dist \;

## Generated server openapi files
$(SERVER_OPENAPI_FILES): $(OPENAPI_SRC_FILES)
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

dist/tracetest-server: $(SERVER_SRC_FILES)
	goreleaser build --single-target --clean --snapshot --id server
	find ./dist -name 'tracetest-server' -exec cp {} ./dist \;

## Generated web openapi files
web/src/types/Generated.types.ts: $(OPENAPI_SRC_FILES)
	cd web; npm run types:generate
web/node_modules: web/package.json web/package-lock.json
	cd web; npm install
dist/html: web/node_modules $(WEB_SRC_FILES)
	cd web; npm run build

dist/docker: Dockerfile dist/html dist/tracetest dist/tracetest-server
	docker build dist/ --file ./Dockerfile --tag kubeshop/tracetest:latest
	touch dist/docker

####################
## Helpful alias  ##
####################
all: build-docker ## build all targets

run: dist/docker ## build and run tracetest using docker compose
	docker compose up

serve-docs: ## serve documentation for Tracetest
	docker build -t tracetest-docs -f docs-Dockerfile .
	docker run --network host tracetest-docs
	sleep 1
	open http://localhost:8000

build-go: dist/tracetest dist/tracetest-server ## build all go code
build-web: dist/html ## build web
build-docker: dist/docker ## build and tag docker image as defined in .goreleaser.dev.yaml

generate: generate-server generate-cli generate-web ## generate code entities from openapi definitions for all parts of the code
generate-server: $(wildcard server/openapi/api.go) ## generate code entities from openapi definitions for server
generate-cli: $(wildcard cli/openapi/*.go) ## generate code entities from openapi definitions for cli
generate-web: web/src/types/Generated.types.ts ## generate code entities from openapi definitions for web
