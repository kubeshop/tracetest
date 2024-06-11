# Dependencies:
# https://grpc.io/docs/protoc-installation/
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

REQUIRED_BINS := protoc protoc-gen-go protoc-gen-go-grpc

help: ## show list of commands
	@echo "Choose a command to run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

build-proto: ensure-dependencies clean-proto ## generate gRPC code from proto files
	@protoc \
		-I=../proto \
		--go_out=./proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./proto \
		--go-grpc_opt=paths=source_relative \
		../proto/paymentreceiver.proto

ensure-dependencies: ## check if required binaries are installed
	$(foreach bin,$(REQUIRED_BINS),\
		$(if $(shell command -v $(bin) 2> /dev/null),,$(error Please install `$(bin)` or run `make install-grpc-tools`)))

install-grpc-tools: ## install required binaries
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

clean-proto: ## remove generated gRPC code
	@rm -f proto/*.go
