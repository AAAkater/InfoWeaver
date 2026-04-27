import tempfile
from pathlib import Path

from fastapi import APIRouter, Depends, HTTPException, status
from llama_index.core.schema import BaseNode, Document
from sqlalchemy.orm import Session

from core.rag.doc_store.document_store import DocumentChunk, add_document_chunks
from core.rag.extractor import common_extractor
from core.rag.splitter import common_splitter
from db.minio_client import minio_client
from db.postgresql_db import get_db_session
from models.document import (
    EmbedChunksData,
    EmbedChunksRequest,
    EmbedChunksResponse,
    ProcessDocumentData,
    ProcessDocumentRequest,
    ProcessDocumentResponse,
    SplitDocumentData,
    SplitDocumentRequest,
    SplitDocumentResponse,
)
from services import document_service
from utils import logger

router = APIRouter(prefix="/documents", tags=["documents"])


@router.post(
    "/split",
    response_model=SplitDocumentResponse,
    status_code=status.HTTP_200_OK,
    summary="Split document from MinIO into chunks",
    description="Download a document from MinIO and split it into text chunks",
)
async def split_document(req: SplitDocumentRequest, db: Session = Depends(get_db_session)) -> SplitDocumentResponse:
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
        SplitDocumentResponse with chunk IDs
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
            return SplitDocumentResponse(
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
        return SplitDocumentResponse(
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
    response_model=EmbedChunksResponse,
    status_code=status.HTTP_200_OK,
    summary="Embed and store document chunks",
    description="Generate embeddings for text chunks and store them in the vector database",
)
async def embed_chunks(req: EmbedChunksRequest, db: Session = Depends(get_db_session)) -> EmbedChunksResponse:
    """
    Embed document chunks and store in Milvus.

    This endpoint:
    1. Fetches chunk contents from PostgreSQL by chunk IDs
    2. Generates dense and sparse embeddings for each chunk
    3. Stores the chunks with embeddings in Milvus
    4. Updates vector_id in PostgreSQL

    Args:
        req: EmbedChunksRequest containing chunk IDs and embedding config

    Returns:
        EmbedChunksResponse with processing results
    """
    try:
        if not req.chunk_ids:
            return EmbedChunksResponse(
                data=EmbedChunksData(chunk_ids=[], chunks_count=0),
                msg="No chunk IDs provided for embedding",
            )

        # Step 1: Fetch chunks from PostgreSQL
        pg_chunks = document_service.fetch_chunks_by_ids(db, req.chunk_ids)
        if not pg_chunks:
            logger.warning(f"No chunks found in PostgreSQL for IDs: {req.chunk_ids}")
            return EmbedChunksResponse(
                data=EmbedChunksData(chunk_ids=[], chunks_count=0),
                msg="No chunks found for the given IDs",
            )

        logger.info(f"Embedding {len(pg_chunks)} chunks from PostgreSQL")
        contents = [chunk.content for chunk in pg_chunks]
        dataset_id = pg_chunks[0].file.dataset_id

        # Step 2: Generate embeddings and store in Milvus
        doc_entities = [
            DocumentChunk(
                content=content,
                dataset_id=dataset_id,
            )
            for content in contents
        ]
        await add_document_chunks(doc_entities, req.embedding_config)

        logger.success(f"Successfully embedded and stored {len(pg_chunks)} chunks to Milvus")
        return EmbedChunksResponse(
            data=EmbedChunksData(chunk_ids=req.chunk_ids, chunks_count=len(pg_chunks)),
            msg="Chunks embedded and stored successfully",
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error embedding chunks: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to embed chunks: {e}",
        )


@router.post(
    "/process",
    response_model=ProcessDocumentResponse,
    status_code=status.HTTP_200_OK,
    summary="Process document from MinIO (split + embed)",
    description="Download a document from MinIO, split it into chunks, embed and store in vector database",
)
async def process_document_from_minio(req: ProcessDocumentRequest) -> ProcessDocumentResponse:
    """
    Process a document stored in MinIO end-to-end.

    This endpoint combines split and embed:
    1. Downloads the file from MinIO storage
    2. Splits the document into chunks
    3. Generates embeddings for each chunk
    4. Stores the chunks with embeddings in Milvus

    Args:
        req: ProcessDocumentRequest containing file metadata and MinIO path

    Returns:
        ProcessDocumentResponse with processing results
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

        if not chunks:
            logger.warning(f"No chunks extracted from document: {file_name}")
            return ProcessDocumentResponse(
                data=ProcessDocumentData(
                    file_id=req.file_id,
                    dataset_id=req.dataset_id,
                    file_name=file_name,
                    chunks_count=0,
                ),
                msg="No chunks extracted from document",
            )

        logger.info(f"Extracted {len(chunks)} chunks from document: {file_name}")

        # Step 3: Store chunks with embeddings in Milvus
        logger.info(f"Storing {len(chunks)} chunks in Milvus")
        doc_entities = [
            DocumentChunk(
                content=chunk.get_content(),
                dataset_id=req.dataset_id,
            )
            for chunk in chunks
        ]
        await add_document_chunks(doc_entities, req.embedding_config)

        logger.success(f"Successfully processed document: {file_name} with {len(chunks)} chunks")
        return ProcessDocumentResponse(
            data=ProcessDocumentData(
                file_id=req.file_id,
                dataset_id=req.dataset_id,
                file_name=file_name,
                chunks_count=len(chunks),
            ),
            msg="Document processed successfully",
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Error processing document {req.minio_path}: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to process document: {e}",
        )
    finally:
        if tmp_path:
            tmp_path.unlink(missing_ok=True)
