generate-proto: clean-proto build-proto

build-proto:
	@protoc \
		--go_out=./ \
		--go_opt=paths=source_relative \
		--go-grpc_out=./ \
		--go-grpc_opt=paths=source_relative \
		proto/orchestrator.proto

clean-proto:
	@rm -f proto/*.go
