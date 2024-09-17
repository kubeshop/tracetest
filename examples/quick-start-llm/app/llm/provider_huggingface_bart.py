
from langchain_community.docstore.document import Document
from langchain_core.prompts import ChatPromptTemplate
from langchain_text_splitters import CharacterTextSplitter
from langchain.chains.combine_documents import create_stuff_documents_chain

from langchain_huggingface import HuggingFaceEndpoint, ChatHuggingFace
from langchain import HuggingFaceHub

import os

class HuggingFaceBartProvider:
  def provider_name(self):
    return "Hugging Face (Bart)"

  def summarize(self, text):
    huggingfacehub_api_token = os.getenv("HUGGINGFACE_HUB_API_TOKEN", "")

    if not huggingfacehub_api_token.strip():
      raise ValueError("Please provide the HuggingFace API Token on a .env file.")

    llm = HuggingFaceHub(
      repo_id="facebook/bart-large-cnn",
      task="summarization",
      model_kwargs={
        "temperature":0,
        "max_length":180,
        'max_new_tokens' : 120,
        'top_k' : 10,
        'top_p': 0.95,
        'repetition_penalty':1.03
      },
      huggingfacehub_api_token=huggingfacehub_api_token
    )

    chat = ChatHuggingFace(llm=llm, verbose=True)

    # Define prompt
    prompt = ChatPromptTemplate.from_messages(
      [
        ("system", "Write a concise summary of the following:\\n\\n{context}"),
        ("human", "")
      ]
    )

    # Instantiate chain
    chain = create_stuff_documents_chain(chat, prompt)

    # Split the source text
    text_splitter = CharacterTextSplitter()
    texts = text_splitter.split_text(text)

    # Create Document objects for the texts (max 3 pages)
    docs = [Document(page_content=t) for t in texts[:3]]

    # Invoke chain
    return chain.invoke({"context": docs})
