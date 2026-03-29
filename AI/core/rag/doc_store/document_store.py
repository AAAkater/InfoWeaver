from pydantic import BaseModel

from db.milvus_db import VectorEntity, client
from utils.logger import logger


class DocumentChunk(BaseModel):
    """Model representing a document chunk with its embedding."""

    content: str
    dataset_id: int


async def add_document_chunks(chunks: list[DocumentChunk]) -> None:
    """Add document chunks to the Milvus database."""
    if not chunks:
        logger.warning("No document chunks to add")
        return

    new_entities = [
        VectorEntity(
            dense_vector=[],  # Placeholder for dense vector, to be filled after embedding generation
            sparse_vector=[],  # Placeholder for sparse vector, to be filled after embedding generation
            content=chunk.content,
            dataset_id=chunk.dataset_id,
        )
        for chunk in chunks
    ]

    client.insert_entities(new_entities)


def delete_document_chunks_by_ids(entity_ids: list[int]) -> None:
    """Delete document chunks from the Milvus database by their IDs."""
    if not entity_ids:
        logger.warning("No entity IDs provided for deletion")
        return

    client.delete_entities_by_ids(entity_ids)


def update_document_chunks(chunks: list[DocumentChunk], entity_ids: list[int]) -> None:
    """Update document chunks in the Milvus database by their IDs."""
    if not chunks or not entity_ids or len(chunks) != len(entity_ids):
        logger.warning("Chunks and entity IDs must be provided and have the same length for update")
        return

    updated_entities = [
        VectorEntity(
            dense_vector=[],  # Placeholder for dense vector, to be filled after embedding generation
            sparse_vector=[],  # Placeholder for sparse vector, to be filled after embedding generation
            content=chunk.content,
            dataset_id=chunk.dataset_id,
        )
        for chunk in chunks
    ]

    client.update_entities_by_ids(entity_ids, updated_entities)
