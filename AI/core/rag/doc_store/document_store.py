from pydantic import BaseModel

from configs.app_config import settings as app_settings
from core.rag.embedding import OAICompatibleEmbedding, OllamaDenseEmbeddingModel
from core.rag.embedding.sparse_embedding_model import SparseEmbeddingModel
from db.milvus_db import VectorEntity, client
from models.document import EmbeddingModelConfig
from utils import logger


class DocumentChunk(BaseModel):
    """Model representing a document chunk with its embedding."""

    content: str
    dataset_id: int


async def add_document_chunks(chunks: list[DocumentChunk], embedding_config: EmbeddingModelConfig) -> None:
    """Add document chunks to the Milvus database with configurable embedding type.

    Supports three embed_type values:
    - "dense": Generate dense embeddings only (sparse set to empty dict)
    - "sparse": Generate sparse embeddings only (dense set to zero vector)
    - "hybrid": Generate both dense and sparse embeddings (default)
    """
    if not chunks:
        logger.warning("No document chunks to add")
        return

    embed_type = embedding_config.embed_type
    logger.info(f"Adding {len(chunks)} document chunks with embed_type='{embed_type}'")

    contents = [chunk.content for chunk in chunks]
    dense_embeddings: list[list[float]] = []
    sparse_embeddings: list[dict[int, float]] = []

    # Generate dense embeddings if needed
    if embed_type in ("dense", "hybrid"):
        if embedding_config.provider_type == "ollama":
            dense_model = OllamaDenseEmbeddingModel(
                model_name=embedding_config.model_name,
                base_url=embedding_config.base_url,
            )
        else:
            dense_model = OAICompatibleEmbedding(
                model_name=embedding_config.model_name,
                base_url=embedding_config.base_url,
                api_key=embedding_config.api_key or "",
            )
        dense_embeddings = await dense_model.get_embeddings(contents)

    # Generate sparse embeddings if needed
    if embed_type in ("sparse", "hybrid"):
        sparse_model = SparseEmbeddingModel()
        sparse_embeddings = await sparse_model.get_embeddings(contents)

    # Build entities with appropriate vectors
    zero_dense = [0.0] * app_settings.MILVUS_DIM
    new_entities = []
    for i, chunk in enumerate(chunks):
        dense_vec = dense_embeddings[i] if i < len(dense_embeddings) else zero_dense
        sparse_vec = sparse_embeddings[i] if i < len(sparse_embeddings) else {}
        new_entities.append(
            VectorEntity(
                dense_vector=dense_vec,
                sparse_vector=sparse_vec,
                content=chunk.content,
                dataset_id=chunk.dataset_id,
            )
        )

    client.insert_entities(new_entities)
    logger.success(f"Inserted {len(new_entities)} chunks into Milvus (embed_type='{embed_type}')")


def delete_document_chunks_by_ids(entity_ids: list[int]) -> None:
    """Delete document chunks from the Milvus database by their IDs."""
    if not entity_ids:
        logger.warning("No entity IDs provided for deletion")
        return

    client.delete_entities_by_ids(entity_ids)


async def update_document_chunks(
    chunks: list[DocumentChunk], entity_ids: list[int], embedding_config: EmbeddingModelConfig
) -> None:
    """Update document chunks in the Milvus database by their IDs with configurable embedding type."""
    if not chunks or not entity_ids or len(chunks) != len(entity_ids):
        logger.warning("Chunks and entity IDs must be provided and have the same length for update")
        return

    embed_type = embedding_config.embed_type
    logger.info(f"Updating {len(chunks)} document chunks with embed_type='{embed_type}'")

    contents = [chunk.content for chunk in chunks]
    dense_embeddings: list[list[float]] = []
    sparse_embeddings: list[dict[int, float]] = []

    # Generate dense embeddings if needed
    if embed_type in ("dense", "hybrid"):
        if embedding_config.provider_type == "ollama":
            dense_model = OllamaDenseEmbeddingModel(
                model_name=embedding_config.model_name,
                base_url=embedding_config.base_url,
            )
        else:
            dense_model = OAICompatibleEmbedding(
                model_name=embedding_config.model_name,
                base_url=embedding_config.base_url,
                api_key=embedding_config.api_key or "",
            )
        dense_embeddings = await dense_model.get_embeddings(contents)

    # Generate sparse embeddings if needed
    if embed_type in ("sparse", "hybrid"):
        sparse_model = SparseEmbeddingModel()
        sparse_embeddings = await sparse_model.get_embeddings(contents)

    zero_dense = [0.0] * app_settings.MILVUS_DIM
    updated_entities = []
    for i, chunk in enumerate(chunks):
        dense_vec = dense_embeddings[i] if i < len(dense_embeddings) else zero_dense
        sparse_vec = sparse_embeddings[i] if i < len(sparse_embeddings) else {}
        updated_entities.append(
            VectorEntity(
                dense_vector=dense_vec,
                sparse_vector=sparse_vec,
                content=chunk.content,
                dataset_id=chunk.dataset_id,
            )
        )

    client.update_entities_by_ids(entity_ids, updated_entities)
