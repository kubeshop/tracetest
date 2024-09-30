## Quick Start LLM app

This is an example of a simple LLM app that uses the `langchain` library to summarize the content of a URL, based on [this example](https://github.com/alphasecio/langchain-examples/tree/main/url-summary)

### Running example with docker

```bash
make start/on-docker
```

### Running example with locally

#### Setting up the environment

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
# add google gemini api key (optional)
echo "GOOGLE_API_KEY={your-google-gemini-api-key}" >> .env

# add tracetest agent keys
echo "TRACETEST_API_KEY={your-tracetest-api-key}" >> .env
echo "TRACETEST_ENVIRONMENT_ID={your-tracetest-env-id}" >> .env

# add tracetest token for playwright tests with TS lib
echo "TRACETEST_API_TOKEN={your-tracetest-token-for-ts-libs}" >> ./tests/.env
```

#### Running the apps

```bash

# start the app running the UI locally
make start/local-ui

# start the app running the API locally
make start/local-api

# start the app everything on docker
make start/on-docker

```
