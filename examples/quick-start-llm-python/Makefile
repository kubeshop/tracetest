help: Makefile ## show list of commands
	@echo "Choose a command to run:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# Added to skip list to avoid warnings
DISABLED_INSTRUMENTATIONS=aleph_alpha_client,chromadb,cohere,groq,haystack-ai,lancedb,llama-index,marqo,milvus,mistralai,pinecone_client,qdrant_client,replicate,together,google_cloud_aiplatform,ibm-watson-machine-learning,weaviate_client

build/docker: ## build images used by docker compose file
	@docker compose build

start/on-docker: build/docker ## build and run app using docker compose
	@docker compose up

start/on-docker/only-observability: ## run observability stack using docker compose in debug mode
	@docker compose up -d otel-collector jaeger

start/local-ui: start/on-docker/only-observability ## run UI app using docker compose
	@OTEL_SERVICE_NAME=quick-start-llm \
	 OTEL_TRACES_EXPORTER=otlp \
	 OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4317 \
	 OTEL_PYTHON_DISABLED_INSTRUMENTATIONS=$(DISABLED_INSTRUMENTATIONS) \
	 opentelemetry-instrument streamlit run ./app/streamlit_app.py

start/local-api: start/on-docker/only-observability ## run UI app using docker compose
	@OTEL_SERVICE_NAME=quick-start-llm \
	 OTEL_TRACES_EXPORTER=otlp \
	 OTEL_METRICS_EXPORTER=none \
	 OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4317 \
	 OTEL_PYTHON_DISABLED_INSTRUMENTATIONS=$(DISABLED_INSTRUMENTATIONS) \
	 opentelemetry-instrument python ./app/flask_app.py

stop: ## stop all running containers
	@docker compose down

test: ## run e2e tests
	@tracetest run -f ./tests/run-gemini.yaml
