OPENAPI_SERVER_TARGET_DIR=./openapi/

generate-client:
	$(eval TMPDIR := $(shell mktemp -d))
	rm -rf $(OPENAPI_SERVER_TARGET_DIR)

	openapi-generator-cli generate -i ../api/openapi.yaml -g go -o $(TMPDIR)
	mkdir -p $(OPENAPI_SERVER_TARGET_DIR)
	mv $(TMPDIR)/*.go $(OPENAPI_SERVER_TARGET_DIR)
	rm -rf $(TMPDIR)

	cd $(OPENAPI_SERVER_TARGET_DIR); go fmt ./...; cd ..;

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
