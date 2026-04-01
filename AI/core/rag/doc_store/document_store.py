from pydantic import BaseModel

from configs.app_config import settings
from core.rag.embedding import OllamaDenseEmbeddingModel
from core.rag.embedding.sparse_embedding_model import SparseEmbeddingModel
from db.milvus_db import VectorEntity, client
from utils import logger


class DocumentChunk(BaseModel):
    """Model representing a document chunk with its embedding."""

    content: str
    dataset_id: int


async def add_document_chunks(chunks: list[DocumentChunk]) -> None:
    """Add document chunks to the Milvus database."""
    if not chunks:
        logger.warning("No document chunks to add")
        return

    sparse_embedding_model = SparseEmbeddingModel()
    dense_embedding_model = OllamaDenseEmbeddingModel(
        model_name=settings.OLLAMA_DENSE_EMBEDDING_MODEL,
        base_url=settings.OLLAMA_URL,
    )

    sparse_embeddings = await sparse_embedding_model.get_embeddings([chunk.content for chunk in chunks])
    dense_embeddings = await dense_embedding_model.get_embeddings([chunk.content for chunk in chunks])

    new_entities = [
        VectorEntity(
            dense_vector=dense_embeddings[i],
            sparse_vector=sparse_embeddings[i],
            content=chunk.content,
            dataset_id=chunk.dataset_id,
        )
        for i, chunk in enumerate(chunks)
    ]

    client.insert_entities(new_entities)


def delete_document_chunks_by_ids(entity_ids: list[int]) -> None:
    """Delete document chunks from the Milvus database by their IDs."""
    if not entity_ids:
        logger.warning("No entity IDs provided for deletion")
        return

    client.delete_entities_by_ids(entity_ids)


async def update_document_chunks(chunks: list[DocumentChunk], entity_ids: list[int]) -> None:
    """Update document chunks in the Milvus database by their IDs."""
    if not chunks or not entity_ids or len(chunks) != len(entity_ids):
        logger.warning("Chunks and entity IDs must be provided and have the same length for update")
        return

    sparse_embedding_model = SparseEmbeddingModel()
    dense_embedding_model = OllamaDenseEmbeddingModel(
        model_name=settings.OLLAMA_DENSE_EMBEDDING_MODEL,
        base_url=settings.OLLAMA_URL,
    )

    sparse_embeddings = await sparse_embedding_model.get_embeddings([chunk.content for chunk in chunks])
    dense_embeddings = await dense_embedding_model.get_embeddings([chunk.content for chunk in chunks])

    updated_entities = [
        VectorEntity(
            dense_vector=dense_embeddings[i],
            sparse_vector=sparse_embeddings[i],
            content=chunk.content,
            dataset_id=chunk.dataset_id,
        )
        for i, chunk in enumerate(chunks)
    ]

    client.update_entities_by_ids(entity_ids, updated_entities)
