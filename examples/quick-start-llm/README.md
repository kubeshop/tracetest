## Quick Start LLM app

This is an example of a simple LLM app that uses the `langchain` library to summarize the content of a URL, based on [this example](https://github.com/alphasecio/langchain-examples/tree/main/url-summary)

### Starting new env from scratch

```bash

# create venv
python -m venv ./_venv

# activate env
source _venv/bin/activate

# install requirements
pip install -r app/requirements.txt

# install OTel auto-instrumentation
opentelemetry-bootstrap -a install

# add openai api key
echo "OPENAI_API_KEY={your-open-ai-api-key}" >> .env
```

### Run example

```bash
OTEL_SERVICE_NAME=quick-start-llm \
OTEL_TRACES_EXPORTER=console,otlp \
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=localhost:4317 \
opentelemetry-instrument \
  streamlit run ./app/streamlit_app.py
```

### References

- https://github.com/openlit/openlit?tab=readme-ov-file#-getting-started
- https://github.com/langchain-ai/langchain
- https://github.com/streamlit/streamlit
- https://docs.streamlit.io/develop/api-reference
- https://wandb.ai/site
