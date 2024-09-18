from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Initialize telemetry
from telemetry import init as telemetry_init, otlp_endpoint
tracer = telemetry_init() # run telemetry.init() before loading any other modules to capture any module-level telemetry

# from telemetry import heartbeat as telemetry_heartbeat
# telemetry_heartbeat(tracer)

import streamlit as st

from llm.providers import get_provider, get_providers

def read_default_text():
  with open("./example.txt") as f:
    return f.read()

@tracer.start_as_current_span("perform_summarization")
def perform_summarization(provider_type, source_text):
  provider = get_provider(provider_type)
  return provider.summarize(source_text)

############################
# UI App start
############################

# Streamlit app
st.subheader('Summarize Text')

st.divider()
st.text(f"OTel Collector endpoint: {otlp_endpoint}")
st.divider()

# Provider checkbox
provider_type = st.selectbox(
  "Choose a provider:",
  get_providers(),
)

# Get Source Text
def callback():
  st.session_state['source_text'] = read_default_text()

st.button("Add example text", on_click=callback)
source_text = st.text_area("Source Text", label_visibility="collapsed", height=250, key="source_text")

# If the 'Summarize' button is clicked
if st.button("Summarize"):
  with tracer.start_as_current_span("summarize button click") as span:
    # Validate inputs
    if not source_text.strip():
      st.error(f"Please provide the source text.")
    else:
      try:
        with st.spinner('Please wait...'):
          summary = perform_summarization(provider_type, source_text)
          st.success(summary)

          # Get trace ID from current span
          trace_id = span.get_span_context().trace_id

          # Convert trace_id to a hex string
          trace_id_hex = format(trace_id, '032x')
          st.text(f"Trace ID: {trace_id_hex}")

          # Add a hyperlink to the trace visualization tool
          trace_url = f"http://localhost:16686//trace/{trace_id_hex}"
          st.markdown(f"[View Trace]({trace_url})")
      except Exception as e:
        st.exception(f"An error occurred: {e}")
