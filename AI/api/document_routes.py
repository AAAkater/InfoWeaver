import tempfile
from pathlib import Path

from fastapi import APIRouter, BackgroundTasks, Depends, HTTPException, status
from llama_index.core.schema import BaseNode, Document
from sqlalchemy.orm import Session

from core.rag.extractor import common_extractor
from core.rag.splitter import common_splitter
from db.minio_client import minio_client
from db.postgresql_db import get_db_session
from models.document import (
    EmbedChunksData,
    EmbedChunksRequest,
    SplitDocumentData,
    SplitDocumentRequest,
)
from models.response import APIResponse
from services import document_service
from services.document_service import background_embed_chunks
from utils import logger

router = APIRouter(prefix="/documents", tags=["documents"])


@router.post(
    "/split",
    response_model=APIResponse[SplitDocumentData],
    status_code=status.HTTP_200_OK,
    summary="Split document from MinIO into chunks",
    description="Download a document from MinIO and split it into text chunks",
)
async def split_document(
    req: SplitDocumentRequest, db: Session = Depends(get_db_session)
) -> APIResponse[SplitDocumentData]:
    """
    Split a document stored in MinIO into chunks.

    This endpoint:
    1. Downloads the file from MinIO storage
    2. Splits the document into chunks
    3. Persists chunk contents to PostgreSQL
    4. Returns the chunk IDs

    Args:
        req: SplitDocumentRequest containing file metadata and MinIO path

    Returns:
        APIResponse[SplitDocumentData] with chunk IDs
    """
    file_name = Path(req.minio_path).stem
    suffix = Path(req.minio_path).suffix
    tmp_path: Path | None = None

    try:
        # Step 1: Download file from MinIO
        logger.info(f"Downloading file from MinIO: {req.minio_path}")
        file_data = minio_client.download_file(req.minio_path)
        with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp_file:
            tmp_file.write(file_data)
            tmp_file.flush()
            tmp_path = Path(tmp_file.name)

        # Step 2: Load and split document
        logger.info(f"Loading and splitting document: {file_name}")
        reader = common_extractor.extractor(str(tmp_path))
        splitter = common_splitter.splitter(
            chunk_size=req.chunk_size,
            chunk_overlap=req.chunk_overlap,
        )
        docs: list[Document] = await reader.aload_data()
        chunks: list[BaseNode] = await splitter.aget_nodes_from_documents(docs)

        chunk_contents = [chunk.get_content() for chunk in chunks]

        if not chunk_contents:
            logger.warning(f"No chunks extracted from document: {file_name}")
            return APIResponse(
                data=SplitDocumentData(
                    file_id=req.file_id,
                    dataset_id=req.dataset_id,
                    file_name=file_name,
                    chunks_count=0,
                ),
                msg="No chunks extracted from document",
            )

        # Step 3: Persist chunks to PostgreSQL
        chunk_ids = document_service.persist_chunks(db, chunk_contents, req.file_id)

        logger.success(
            f"Successfully split document: {file_name} into {len(chunk_ids)} chunks, "
            f"persisted to PostgreSQL with IDs: {chunk_ids}"
        )
        return APIResponse(
            data=SplitDocumentData(
                file_id=req.file_id,
                dataset_id=req.dataset_id,
                file_name=file_name,
                chunks_count=len(chunk_ids),
            ),
            msg="Document split successfully",
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error splitting document {req.minio_path}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to split document: {e}",
        )
    finally:
        if tmp_path:
            tmp_path.unlink(missing_ok=True)


@router.post(
    "/embedding",
    response_model=APIResponse[EmbedChunksData],
    status_code=status.HTTP_202_ACCEPTED,
    summary="Embed and store document chunks (background task)",
    description="Enqueue chunk embedding as a background task. Status can be queried via chunk status field.",
)
async def embed_chunks(
    req: EmbedChunksRequest,
    background_tasks: BackgroundTasks,
    db: Session = Depends(get_db_session),
) -> APIResponse[EmbedChunksData]:
    """
    Enqueue chunk embedding as a background task.

    This endpoint:
    1. Validates that all chunk IDs exist in PostgreSQL
    2. Enqueues a background task to generate embeddings and store in Milvus
    3. Returns immediately with 202 Accepted
    4. Chunk status transitions: pending → embedding → completed/failed

    Args:
        req: EmbedChunksRequest containing chunk IDs and embedding config

    Returns:
        APIResponse[EmbedChunksData] acknowledging the request
    """
    if not req.chunk_ids:
        return APIResponse(
            data=EmbedChunksData(chunk_ids=[], chunks_count=0),
            msg="No chunk IDs provided for embedding",
        )

    # Validate that all specified chunks exist
    pg_chunks = document_service.fetch_chunks_by_ids(db, req.chunk_ids)
    if not pg_chunks:
        logger.warning(f"No chunks found in PostgreSQL for IDs: {req.chunk_ids}")
        return APIResponse(
            data=EmbedChunksData(chunk_ids=[], chunks_count=0),
            msg="No chunks found for the given IDs",
        )

    # Enqueue background task
    background_tasks.add_task(background_embed_chunks, req.chunk_ids, req.embedding_config)

    logger.info(f"Enqueued background embedding for {len(req.chunk_ids)} chunks")
    return APIResponse(
        data=EmbedChunksData(chunk_ids=req.chunk_ids, chunks_count=len(pg_chunks)),
        msg="Chunk embedding enqueued as background task",
    )
