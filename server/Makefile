VERSION ?= "dev"
TRACETEST_ENV ?= "dev"
GO_LDFLAGS := $(shell echo \
		-X "'github.com/kubeshop/tracetest/server/app.Version=$(VERSION)'" \
		-X "'github.com/kubeshop/tracetest/server/app.Env=$(TRACETEST_ENV)'" \
		-X "'github.com/kubeshop/tracetest/server/analytics.SecretKey=$(ANALYTICS_BE_KEY)'" \
		-X "'github.com/kubeshop/tracetest/server/analytics.FrontendKey=$(ANALYTICS_FE_KEY)'" \
	| sed 's/ / /g')

help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

all: init-submodule proto generate

init-submodule:
	git submodule init
	git submodule update

test: # run go tests for this application
	go test -timeout 90s -coverprofile=coverage.out ./...

vet: ## run vet tool to analyze the code for suspicious, abnormal, or useless code
	go vet -structtag=false ./...

run: ## run server locally
	go run -ldflags="$(GO_LDFLAGS)" main.go serve

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

proto: ## generate code client from protobuf definitions
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
