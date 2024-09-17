## Quick Start LLM app

This is an example of a simple LLM app that uses the `langchain` library to summarize the content of a URL, based on [this example](https://github.com/alphasecio/langchain-examples/tree/main/url-summary)

### Starting new env from scratch

```bash

# create venv
python -m venv ./_venv

# activate env
source _venv/bin/activate

# install requirements
pip install -r requirements.txt

# add openai api key
echo "OPENAI_API_KEY={your-open-ai-api-key}" >> .env
```

### Run example

```bash
streamlit run ./app/streamlit_app.py
```

### References

- https://github.com/openlit/openlit?tab=readme-ov-file#-getting-started
- https://github.com/langchain-ai/langchain
- https://github.com/streamlit/streamlit
- https://docs.streamlit.io/develop/api-reference
- https://wandb.ai/site
