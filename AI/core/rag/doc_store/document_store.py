from pydantic import BaseModel

from db.milvus_db import VectorEntity, milvus_db
from utils.logger import logger


class DocumentChunk(BaseModel):
    """Model representing a document chunk with its embedding."""

    content: str
    vector: list[float]
    dataset_id: int


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
            vector=chunk.vector,
            content=chunk.content,
            dataset_id=chunk.dataset_id,
        )
        for chunk in chunks
    ]

    milvus_db.insert_entities(entities)
    logger.info(f"Added {len(chunks)} chunks to collection '{milvus_db.collection_name}'")
    return len(chunks)


def search_chunks(
    query_vector: list[float],
    dataset_id: int,
    top_k: int = 10,
) -> list[dict]:
    """
    Search for similar documents.

    Args:
        query_vector: Query embedding vector.
        top_k: Number of results to return.
        dataset_id: Optional dataset ID to filter results.

    Returns:
        list[dict]: List of search results with content and distance.
    """

    results = milvus_db.search_entities_filter(
        query_vector=query_vector,
        dataset_id=dataset_id,
        top_k=top_k,
    )

    if results and len(results) > 0:
        return [
            {
                "id": hit["id"],
                "content": hit["entity"]["content"],
                "dataset_id": hit["entity"]["dataset_id"],
                "distance": hit["distance"],
            }
            for hit in results[0]
        ]
    return []


def delete_chunks_by_dataset_id(dataset_id: int) -> None:
    """
    Delete all chunks associated with a dataset.

    Args:
        dataset_id: The dataset ID to delete chunks for.
    """
    milvus_db.delete_entities_by_filter(expr=f"dataset_id == {dataset_id}")
    logger.info(f"Deleted chunks for dataset_id {dataset_id}")


def delete_chunks_by_ids(ids: list[int]) -> None:
    """
    Delete chunks by their IDs.

    Args:
        ids: List of chunk IDs to delete.
    """
    milvus_db.delete_entities_by_id(entity_ids=ids)
    logger.info(f"Deleted chunks with IDs: {ids}")
