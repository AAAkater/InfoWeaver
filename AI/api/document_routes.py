"""Document processing API routes."""

import tempfile
from pathlib import Path

from fastapi import APIRouter, HTTPException, status

from core.rag.doc_store.document_store import DocumentChunk
from core.rag.splitter.text_splitter import load_and_split_document
from db.minio_client import minio_client
from models.document import ProcessDocumentRequest, ProcessDocumentResponse
from utils import logger

router = APIRouter()


@router.post(
    "/process",
    response_model=ProcessDocumentResponse,
    status_code=status.HTTP_200_OK,
    summary="Process document from MinIO",
    description="Download a document from MinIO, split it into chunks, embed and store in vector database",
)
async def process_document_from_minio(request: ProcessDocumentRequest) -> ProcessDocumentResponse:
    """
    Process a document stored in MinIO.

    This endpoint:
    1. Downloads the file from MinIO storage
    2. Splits the document into chunks
    3. Generates embeddings for each chunk
    4. Stores the chunks with embeddings in Milvus

    Args:
        request: ProcessDocumentRequest containing file metadata and MinIO path

    Returns:
        ProcessDocumentResponse with processing results

    Raises:
        HTTPException: If file download, processing, or storage fails
    """
    file_name = Path(request.minio_path).stem
    suffix = Path(request.minio_path).suffix

    tmp_path: Path | None = None
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp_file:
        try:
            # Step 1: Download file from MinIO
            logger.info(f"Downloading file from MinIO: {request.minio_path}")
            try:
                file_data = minio_client.download_file(request.minio_path)
                tmp_file.write(file_data)
                tmp_file.flush()
                tmp_path = Path(tmp_file.name)
            except Exception as e:
                logger.error(f"Failed to download file from MinIO: {e}")
                raise HTTPException(
                    status_code=status.HTTP_404_NOT_FOUND,
                    detail=f"Failed to download file from MinIO: {e}",
                )

            # Step 2: Load and split document
            logger.info(f"Loading and splitting document: {file_name}")
            try:
                chunks = load_and_split_document(
                    file_path=tmp_path,
                    chunk_size=512,
                    chunk_overlap=50,
                )

                if not chunks:
                    logger.warning(f"No chunks extracted from document: {file_name}")
                    return ProcessDocumentResponse(
                        file_id=request.file_id,
                        dataset_id=request.dataset_id,
                        file_name=file_name,
                        chunks_count=0,
                        message="No chunks extracted from document",
                    )

                logger.info(f"Extracted {len(chunks)} chunks from document: {file_name}")
            except Exception as e:
                logger.error(f"Failed to load and split document: {e}")
                raise HTTPException(
                    status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
                    detail=f"Failed to process document: {e}",
                )

            # Step 3: Store chunks with embeddings in Milvus
            logger.info(f"Storing {len(chunks)} chunks in Milvus")
            try:
                from core.rag.doc_store.document_store import add_document_chunks

                doc_entities = [
                    DocumentChunk(
                        content=chunk,
                        dataset_id=request.dataset_id,
                    )
                    for chunk in chunks
                ]

                await add_document_chunks(doc_entities)
            except Exception as e:
                logger.error(f"Failed to store chunks in Milvus: {e}")
                raise HTTPException(
                    status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                    detail=f"Failed to store chunks in vector database: {e}",
                )

            logger.success(f"Successfully processed document: {file_name} with {len(chunks)} chunks")

            return ProcessDocumentResponse(
                file_id=request.file_id,
                dataset_id=request.dataset_id,
                file_name=file_name,
                chunks_count=len(chunks),
                message="Document processed successfully",
            )

        except HTTPException:
            raise
        except Exception as e:
            logger.error(f"Error processing document {request.minio_path}: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail=f"Unexpected error processing document: {e}",
            )
        finally:
            # Clean up temporary file
            if tmp_path:
                tmp_path.unlink(missing_ok=True)
