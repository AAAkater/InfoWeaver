from llama_index.embeddings.ollama import OllamaEmbedding

from configs.app_config import settings

embedding_model = OllamaEmbedding(
    model_name=settings.OLLAMA_EMBEDDING_MODEL,
    base_url=settings.OLLAMA_URL,
)
