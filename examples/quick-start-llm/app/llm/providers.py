from .provider_google_gemini import GoogleGeminiProvider
from .provider_openai_chatgpt import OpenAIChatGPTProvider

_providers = [
  GoogleGeminiProvider(),
  OpenAIChatGPTProvider()
]

def get_providers():
  return list(map(lambda p: p.provider_name(), _providers))

def get_provider(provider_name):
  for provider in _providers:
    if provider.provider_name() == provider_name:
      return provider
  return None


