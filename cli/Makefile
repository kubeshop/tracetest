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

	cd openapi; pwd; go fmt ./...; cd ..

build: generate-client
	go build -o tracetest main.go

test: mockserver
	@go test -coverprofile=coverage.out ./...
	@make --no-print-directory kill-mock

mockserver:
	@make --no-print-directory kill-mock
	@echo "Starting prism mock server"
	@prism mock ../api/openapi.yaml &
	@sleep 2

kill-mock:
	@echo "Stopping prism mock server"
	@ps -ef | grep /bin/prism | grep -v grep | awk '{print $$2}' | xargs -r kill -9
