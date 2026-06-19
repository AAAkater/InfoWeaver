import ollama
from llama_index.embeddings.ollama import OllamaEmbedding

from utils import logger


class OllamaDenseEmbeddingModel:
    """Wrapper class for Ollama embedding model to provide a consistent interface."""

    def __init__(self, model_name: str, base_url: str, dims: int | None = 1024):
        self.dims = dims
        self.model_name = model_name
        self.base_url = base_url
        # self.embedding_model = OllamaEmbedding(model_name=model_name, base_url=base_url)
        self.embedding_model = ollama.AsyncClient(host=base_url)

    async def get_embeddings(self, texts: list[str]) -> list[list[float]]:
        """Get embeddings for multiple texts in batches."""
        if not texts:
            return []

        logger.info(f"Generating embeddings for {len(texts)} texts using model '{self.model_name}'")
        # embeddings = await self.embedding_model.aget_text_embedding_batch(texts)
        resp = await self.embedding_model.embed(self.model_name, input=texts, dimensions=self.dims)
        embeddings = [list(e) for e in resp.embeddings]
        logger.info(f"Generated {len(embeddings)} embeddings using model '{self.model_name}'")
        return embeddings
