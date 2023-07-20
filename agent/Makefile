# Dependencies:
# https://grpc.io/docs/protoc-installation/

REQUIRED_BINS := protoc

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
