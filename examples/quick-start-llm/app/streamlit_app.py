from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Initialize OpenLit
from telemetry import init as telemetry_init, heartbeat as telemetry_heartbeat, otlp_endpoint
tracer = telemetry_init() # run telemetry.init() before loading any other modules to capture any module-level telemetry
telemetry_heartbeat(tracer)

import streamlit as st

from llm.providers import get_provider, get_providers

def read_default_text():
  with open("./example.txt") as f:
    return f.read()

############################
# App start
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

st.button("Add default text", on_click=callback)
source_text = st.text_area("Source Text", label_visibility="collapsed", height=250, key="source_text")

# If the 'Summarize' button is clicked
if st.button("Summarize"):
  # Validate inputs
  if not source_text.strip():
    st.error(f"Please provide the source text.")
  else:
    try:
      with st.spinner('Please wait...'):
        provider = get_provider(provider_type)
        summary = provider.summarize(source_text)
        st.success(summary)
    except Exception as e:
      st.exception(f"An error occurred: {e}")
