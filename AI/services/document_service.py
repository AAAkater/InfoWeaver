"""Document service — database operations for document chunks."""

from sqlalchemy.orm import Session

from models.tables import Chunk
from utils import logger


def persist_chunks(db: Session, contents: list[str], file_id: int) -> list[int]:
    """Persist chunk contents to PostgreSQL and return their IDs.

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
