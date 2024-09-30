from .provider_google_gemini import GoogleGeminiProvider
from .provider_openai_chatgpt import OpenAIChatGPTProvider

_providers = [
  GoogleGeminiProvider(),
  OpenAIChatGPTProvider()
]

def get_providers():
  providers = []

  for provider in _providers:
    if provider.enabled():
      providers.append(provider.provider_name())

  return providers

def get_provider(provider_name):
  for provider in _providers:
    if provider.provider_name() == provider_name:
      return provider
  return None


