all: init-submodule proto generate

init-submodule:
	git submodule init
	git submodule update

test:
	go tes -timeout 90s -coverprofile=coverage.out ./...

vet:
	go vet -structtag=false ./...

run:
	go run main.go

OPENAPI_GENERATOR_VER=v5.4.0
OPENAPI_GENERATOR_IMAGE=openapitools/openapi-generator-cli:$(OPENAPI_GENERATOR_VER)
OPENAPI_GENERATOR_CLI=docker run --rm -u ${shell id -u}  -v "${PWD}:/app" -v "${PWD}/../api:/app/api" -w "/app" ${OPENAPI_GENERATOR_IMAGE}
OPENAPI_SERVER_TARGET_DIR=./openapi
generate:
	rm -rf $(OPENAPI_SERVER_TARGET_DIR)
	mkdir -p ./tmp

	$(OPENAPI_GENERATOR_CLI)  generate \
		-i api/openapi.yaml \
		-g go-server \
		-o ./tmp \
		--generate-alias-as-model
	mv ./tmp/go $(OPENAPI_SERVER_TARGET_DIR)
	rm -f $(OPENAPI_SERVER_TARGET_DIR)/api_api_service.go
	rm -rf ./tmp

	go fmt ./...

PROTOC_VER=0.3.1
UNAME_P := $(shell uname -p)
ifeq ($(UNAME_P),x86_64)
	PROTOC_IMAGE=jaegertracing/protobuf:$(PROTOC_VER)
endif
ifeq ($(UNAME_P),i386)
	PROTOC_IMAGE=jaegertracing/protobuf:$(PROTOC_VER)
endif
ifneq ($(filter arm%,$(UNAME_P)),)
	PROTOC_IMAGE=schoren/protobuf:$(PROTOC_VER)
endif
PROTOC=docker run --rm -u ${shell id -u} -v "${PWD}:${PWD}" -w ${PWD} ${PROTOC_IMAGE} --proto_path=${PWD}



PROTO_INCLUDES := \
	-I/usr/include/github.com/gogo/protobuf \
	-Iexternal/jaeger-idl/proto \
	-Iexternal/opentelemetry-proto \
	-Iexternal/opentelemetry-proto/opentelemetry/proto

PROTO_GOGO_MAPPINGS := $(shell echo \
		Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/types, \
		Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types, \
		Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types, \
		Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types, \
		Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api, \
		Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api, \
		Mopentelemetry/proto/trace/v1/trace.proto=go.opentelemetry.io/proto/otlp/trace/v1, \
		Mtrace/v1/trace.proto=go.opentelemetry.io/proto/otlp/trace/v1, \
		Mmodel.proto=github.com/jaegertracing/jaeger/model \
	| sed 's/ //g')

PROTO_GEN_GO_DIR ?= internal/proto-gen-go

PROTOC_WITH_GRPC := $(PROTOC) \
		$(PROTO_INCLUDES) \
		--gogo_out=plugins=grpc,$(PROTO_GOGO_MAPPINGS):$(PWD)/${PROTO_GEN_GO_DIR}

PROTOC_INTERNAL := $(PROTOC) \
		$(PROTO_INCLUDES)

proto:
	rm -rf ./$(PROTO_GEN_GO_DIR)
	mkdir -p ${PROTO_GEN_GO_DIR}

	# API v3
	$(PROTOC_WITH_GRPC) \
		external/jaeger-idl/proto/api_v3/query_service.proto

	$(PROTOC_INTERNAL) \
		google/api/annotations.proto \
		google/api/http.proto \
		gogoproto/gogo.proto

	$(PROTOC_WITH_GRPC) \
		external/tempo-idl/tempo.proto
	mv $(PROTO_GEN_GO_DIR)/external/tempo-idl $(PROTO_GEN_GO_DIR)/tempo-idl
	cp external/tempo-idl/prealloc.go.tmp $(PROTO_GEN_GO_DIR)/tempo-idl/prealloc.go
