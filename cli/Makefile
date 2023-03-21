
help: Makefile ## show list of commands
	@echo "Choose a command run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

test: ## execute unit tests
	@go test -coverprofile=coverage.out ./...

docgen: ## generate CLI documentation
	@go run main.go docgen -o ../docs/docs/cli/reference/ -d ../docs/docs
