from pydantic import BaseModel

from configs.app_config import settings
from db.milvus_db import VectorEntity, client
from utils.logger import logger


class DocumentChunk(BaseModel):
    """Model representing a document chunk with its embedding."""

    content: str
    vector: list[float]
    dataset_id: int


def insert_entities(entities: list[VectorEntity]) -> None:
    """Insert vector entities into the collection.

    Args:
        entities: List of VectorEntity objects to insert.
    """
    client.insert(
        collection_name=settings.MILVUS_COLLECTION_NAME,
        data=[data.model_dump() for data in entities],
    )
    logger.info(f"Inserted {len(entities)} records into collection '{settings.MILVUS_COLLECTION_NAME}'")


def delete_entities_by_id(entity_ids: list[int]) -> None:
    """Delete entities by their IDs.

    Args:
        entity_ids: List of entity IDs to delete.
    """
    client.delete(
        collection_name=settings.MILVUS_COLLECTION_NAME,
        ids=entity_ids,
    )
    logger.info(f"Deleted entities with IDs {entity_ids} from collection '{settings.MILVUS_COLLECTION_NAME}'")


def delete_entities_by_filter(expr: str) -> None:
    """Delete entities by filter expression.

    Args:
        expr: Filter expression for deletion.
    """
    client.delete(
        collection_name=settings.MILVUS_COLLECTION_NAME,
        filter_expression=expr,
    )
    logger.info(f"Deleted entities with filter expression '{expr}' from collection '{settings.MILVUS_COLLECTION_NAME}'")


async def add_chunks(chunks: list[DocumentChunk]) -> int:
    """
    Add document chunks to the store.

    Args:
        chunks: List of DocumentChunk objects to add.

    Returns:
        int: Number of chunks added.
    """
    if not chunks:
        logger.warning("No chunks to add")
        return 0

    entities = [
        VectorEntity(
            dense_vector=chunk.vector,
            sparse_vector={},  # TODO: add sparse vector support
            content=chunk.content,
            dataset_id=chunk.dataset_id,
        )
        for chunk in chunks
    ]

    insert_entities(entities)
    logger.info(f"Added {len(chunks)} chunks to collection")
    return len(chunks)


def delete_chunks_by_dataset_id(dataset_id: int) -> None:
    """
    Delete all chunks associated with a dataset.

    Args:
        dataset_id: The dataset ID to delete chunks for.
    """
    delete_entities_by_filter(expr=f"dataset_id == {dataset_id}")
    logger.info(f"Deleted chunks for dataset_id {dataset_id}")


def delete_chunks_by_ids(ids: list[int]) -> None:
    """
    Delete chunks by their IDs.

    Args:
        ids: List of chunk IDs to delete.
    """
    delete_entities_by_id(entity_ids=ids)
    logger.info(f"Deleted chunks with IDs: {ids}")
