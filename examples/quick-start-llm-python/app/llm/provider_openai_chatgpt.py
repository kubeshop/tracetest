from langchain_community.docstore.document import Document
from langchain_core.prompts import ChatPromptTemplate
from langchain_text_splitters import CharacterTextSplitter
from langchain.chains.combine_documents import create_stuff_documents_chain

from langchain_openai import ChatOpenAI

import streamlit as st
import os

class OpenAIChatGPTProvider:
  def provider_name(self):
    return "OpenAI (ChatGPT)"

  def summarize(self, text):
    # Get OpenAI API key and URL to be summarized
    openai_api_key = os.getenv("OPENAI_API_KEY", "")

    if not openai_api_key.strip():
      raise ValueError("Please provide the OpenAI API Key on a .env file.")

    llm = ChatOpenAI(
      model="gpt-4o-mini",
      openai_api_key=openai_api_key
    )

    # Define prompt
    prompt = ChatPromptTemplate.from_messages(
        [("system", "Write a concise summary of the following:\\n\\n{context}")]
    )

    # Instantiate chain
    chain = create_stuff_documents_chain(llm, prompt)

    # Split the source text
    text_splitter = CharacterTextSplitter()
    texts = text_splitter.split_text(text)

    # Create Document objects for the texts (max 3 pages)
    docs = [Document(page_content=t) for t in texts[:3]]

    # Invoke chain
    return chain.invoke({"context": docs})
