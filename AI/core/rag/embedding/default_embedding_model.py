from llama_index.embeddings.ollama import OllamaEmbedding

from configs.app_config import settings
from utils.logger import logger

embedding_model = OllamaEmbedding(
    model_name=settings.OLLAMA_EMBEDDING_MODEL,
    base_url=settings.OLLAMA_URL,
)


async def get_text_embedding(text: str) -> list[float]:
    """Get embedding for a single text."""
    return await embedding_model.aget_text_embedding(text)


async def get_text_embeddings(texts: list[str], batch_size: int = 32) -> list[list[float]]:
    """Get embeddings for multiple texts in batches."""
    if not texts:
        return []

    logger.info(f"Generating embeddings for {len(texts)} texts")
    all_embeddings: list[list[float]] = []

    for i in range(0, len(texts), batch_size):
        batch = texts[i : i + batch_size]
        embeddings = await embedding_model.aget_text_embedding_batch(batch)
        all_embeddings.extend(embeddings)

    logger.info(f"Generated {len(all_embeddings)} embeddings")
    return all_embeddings


def get_text_embedding_sync(text: str) -> list[float]:
    """Synchronously get embedding for a single text."""
    return embedding_model.get_text_embedding(text)


def get_text_embeddings_sync(texts: list[str], batch_size: int = 32) -> list[list[float]]:
    """Synchronously get embeddings for multiple texts in batches."""
    if not texts:
        return []

    logger.info(f"Generating embeddings for {len(texts)} texts")
    all_embeddings: list[list[float]] = []

    for i in range(0, len(texts), batch_size):
        batch = texts[i : i + batch_size]
        embeddings = embedding_model.get_text_embedding_batch(batch)
        all_embeddings.extend(embeddings)

    logger.info(f"Generated {len(all_embeddings)} embeddings")
    return all_embeddings