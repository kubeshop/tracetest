PROJECT_ROOT=${shell dirname $${PWD}}

OPENAPI_GENERATOR_VER=v5.4.0
OPENAPI_GENERATOR_IMAGE=openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VER)
OPENAPI_TARGET_DIR=./openapi/


generate: generate-cli generate-server

generate-cli:
	rm -rf $(OPENAPI_TARGET_DIR)
	mkdir -p ./tmp
	mkdir -p $(OPENAPI_TARGET_DIR)

	docker run --rm -u ${shell id -u}  -v "$(PROJECT_ROOT)/cli:/local" -w "/local" ${OPENAPI_GENERATOR_IMAGE} \
		generate \
		-i api/openapi.yaml \
		-g go \
		-o ./cli/tmp \
		--generate-alias-as-model
	cp ./tmp/*.go $(OPENAPI_TARGET_DIR)
	rm -rf ./tmp

	go fmt ./...; cd ..


generate-server:
	rm -rf $(OPENAPI_SERVER_TARGET_DIR)
	mkdir -p ./tmp

	docker run --rm -u ${shell id -u}  -v "$(PROJECT_ROOT)/server:/local" -w "/local" ${OPENAPI_GENERATOR_IMAGE} \
		generate \
		-i api/openapi.yaml \
		-g go-server \
		-o ./tmp \
		--generate-alias-as-model
	mv ./tmp/go $(OPENAPI_SERVER_TARGET_DIR)
	rm -f $(OPENAPI_SERVER_TARGET_DIR)/api_api_service.go
	rm -rf ./tmp

	go fmt ./...
