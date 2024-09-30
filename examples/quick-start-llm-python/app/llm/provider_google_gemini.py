from langchain.chains.summarize import load_summarize_chain
from langchain_community.docstore.document import Document
from langchain_text_splitters import CharacterTextSplitter

from langchain_google_genai import ChatGoogleGenerativeAI

class GoogleGeminiProvider:
  def provider_name(self):
    return "Google (Gemini)"

  def summarize(self, text):
    chat = ChatGoogleGenerativeAI(model="gemini-pro")

    # Split the source text
    text_splitter = CharacterTextSplitter()
    texts = text_splitter.split_text(text)

    # Create Document objects for the texts (max 3 pages)
    docs = [Document(page_content=t) for t in texts[:3]]

    # Load and run the summarize chain
    chain = load_summarize_chain(chat, chain_type="map_reduce")
    return chain.run(docs)
