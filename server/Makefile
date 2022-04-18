fmt:
	go fmt ./...

CONFIG_FILE = ./testdata/custom-config.yaml
ifeq ("$(wildcard $(CONFIG_FILE))","")
    CONFIG_FILE = ./testdata/local-config.yaml
endif
run-server:
	go build
	./server --config $(CONFIG_FILE)
