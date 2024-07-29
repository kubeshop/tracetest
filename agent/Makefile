# Dependencies:
# https://grpc.io/docs/protoc-installation/
# go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

REQUIRED_BINS := protoc protoc-gen-go protoc-gen-go-grpc

build-proto: ensure-dependencies clean-proto
	@protoc \
		--go_out=./ \
		--go_opt=paths=source_relative \
		--go-grpc_out=./ \
		--go-grpc_opt=paths=source_relative \
		proto/orchestrator.proto

ensure-dependencies:
	$(foreach bin,$(REQUIRED_BINS),\
		$(if $(shell command -v $(bin) 2> /dev/null),,$(error Please install `$(bin)`)))

clean-proto:
	@rm -f proto/*.go

test:
	go test -timeout 150s -coverprofile=coverage.out ./...
