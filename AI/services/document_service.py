"""Document service — database operations for document chunks."""

from sqlalchemy.orm import Session

from core.rag.doc_store.document_store import DocumentChunk, add_document_chunks
from models.document import EmbeddingModelConfig
from models.tables import Chunk
from utils import logger


def persist_chunks(db: Session, contents: list[str], file_id: int) -> list[int]:
    """Persist chunk contents to PostgreSQL and return their IDs.

    Chunks are created with status='pending'.

    Args:
        db: SQLAlchemy database session.
        contents: List of chunk content strings.
        file_id: The file ID to associate chunks with.

    Returns:
        List of auto-generated chunk IDs.
    """
    if not contents:
        logger.warning("No contents provided to persist_chunks")
        return []

    logger.info(f"Persisting {len(contents)} chunks to PostgreSQL for file_id={file_id}")
    db_chunks = [Chunk(content=content, file_id=file_id) for content in contents]
    db.add_all(db_chunks)
    db.commit()
    for chunk in db_chunks:
        db.refresh(chunk)
    chunk_ids = [chunk.id for chunk in db_chunks]
    logger.info(f"Persisted chunks with IDs: {chunk_ids}")
    return chunk_ids


def fetch_chunks_by_ids(db: Session, chunk_ids: list[int]) -> list[Chunk]:
    """Fetch chunks from PostgreSQL by their IDs.

    Args:
        db: SQLAlchemy database session.
        chunk_ids: List of chunk IDs to fetch.

    Returns:
        List of Chunk ORM objects. Empty list if none found.
    """
    if not chunk_ids:
        logger.warning("No chunk IDs provided to fetch_chunks_by_ids")
        return []

    logger.info(f"Fetching {len(chunk_ids)} chunks from PostgreSQL")
    chunks = db.query(Chunk).filter(Chunk.id.in_(chunk_ids)).all()
    logger.info(f"Found {len(chunks)} chunks")
    return chunks


def update_chunk_status(db: Session, chunk_ids: list[int], status: str) -> None:
    """Update the status of multiple chunks.

    Args:
        db: SQLAlchemy database session.
        chunk_ids: List of chunk IDs to update.
        status: New status value (pending | embedding | completed | failed).
    """
    if not chunk_ids:
        return
    logger.info(f"Updating status of {len(chunk_ids)} chunks to '{status}'")
    db.query(Chunk).filter(Chunk.id.in_(chunk_ids)).update({"status": status}, synchronize_session=False)
    db.commit()


async def background_embed_chunks(chunk_ids: list[int], embedding_config: EmbeddingModelConfig) -> None:
    """Run embedding in the background with status tracking.

    This function manages its own database sessions:
    1. Sets chunk status to 'embedding'
    2. Generates dense and sparse embeddings for each chunk
    3. Stores the chunks with embeddings in Milvus
    4. Sets chunk status to 'completed' or 'failed'

    Args:
        chunk_ids: List of chunk IDs to embed.
        embedding_config: Embedding model configuration.
    """
    if not chunk_ids:
        logger.warning("No chunk IDs provided for background embedding")
        return

    from db.postgresql_db import Session as DBSession
    from db.postgresql_db import engine

    with DBSession(engine) as db:
        try:
            # Step 1: Fetch chunks and set status to 'embedding'
            update_chunk_status(db, chunk_ids, "embedding")
            pg_chunks = fetch_chunks_by_ids(db, chunk_ids)
            if not pg_chunks:
                logger.warning(f"No chunks found for IDs: {chunk_ids}")
                return

            logger.info(f"Background embedding {len(pg_chunks)} chunks")
            contents = [chunk.content for chunk in pg_chunks]
            dataset_id = pg_chunks[0].file.dataset_id

            # Step 2: Generate embeddings and store in Milvus
            doc_entities = [DocumentChunk(content=content, dataset_id=dataset_id) for content in contents]
            await add_document_chunks(doc_entities, embedding_config)

            # Step 3: Mark as completed
            update_chunk_status(db, chunk_ids, "completed")
            logger.success(f"Background embedding completed for {len(chunk_ids)} chunks")

        except Exception as e:
            logger.error(f"Background embedding failed for chunk IDs {chunk_ids}: {e}")
            with DBSession(engine) as recovery_db:
                update_chunk_status(recovery_db, chunk_ids, "failed")
