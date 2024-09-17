help: Makefile ## show list of commands
	@echo "Choose a command to run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

build-images: ## build images used by docker compose file
	@docker compose build

start/with-docker: build-images ## build and run app using docker compose
	@docker compose up

start/with-docker/only-observability: ## run observability stack using docker compose in debug mode
	@docker compose up otel-collector jaeger

start/with-local: ## run app using docker compose
	@docker compose up -d otel-collector jaeger
	@OTEL_SERVICE_NAME=quick-start-llm OTEL_TRACES_EXPORTER=otlp OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4317 opentelemetry-instrument streamlit run ./app/streamlit_app.py

stop: ## stop all running containers
	@docker compose down
