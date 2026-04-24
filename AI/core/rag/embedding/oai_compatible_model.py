from llama_index.embeddings.openai import OpenAIEmbedding

from utils import logger


class OAICompatibleEmbedding:
    """Wrapper class for OpenAI-compatible embedding models to provide a consistent interface."""

    def __init__(self, model_name: str, base_url: str, api_key: str):
        self.model_name = model_name
        self.base_url = base_url
        self.embedding_model = OpenAIEmbedding(model=model_name, api_base=base_url, api_key=api_key)

    async def get_embedding(self, text: str) -> list[float]:
        """Get embedding for a single text."""
        return await self.embedding_model.aget_text_embedding(text)

    async def get_embeddings(self, texts: list[str], batch_size: int = 32) -> list[list[float]]:
        """Get embeddings for multiple texts in batches."""
        if not texts:
            return []

        logger.info(f"Generating embeddings for {len(texts)} texts using model '{self.model_name}'")
        all_embeddings: list[list[float]] = []

        for i in range(0, len(texts), batch_size):
            batch = texts[i : i + batch_size]
            embeddings = await self.embedding_model.aget_text_embedding_batch(batch)
            all_embeddings.extend(embeddings)

        logger.info(f"Generated {len(all_embeddings)} embeddings using model '{self.model_name}'")
        return all_embeddings
