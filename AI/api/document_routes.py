import tempfile
from pathlib import Path

from fastapi import APIRouter, HTTPException, status
from llama_index.core.schema import BaseNode, Document

from core.rag.doc_store.document_store import DocumentChunk, add_document_chunks
from core.rag.extractor import common_extractor
from core.rag.splitter import common_splitter
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
async def process_document_from_minio(req: ProcessDocumentRequest) -> ProcessDocumentResponse:
    """
    Process a document stored in MinIO.

    This endpoint:
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
                file_id=req.file_id,
                dataset_id=req.dataset_id,
                file_name=file_name,
                chunks_count=0,
                message="No chunks extracted from document",
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
            file_id=req.file_id,
            dataset_id=req.dataset_id,
            file_name=file_name,
            chunks_count=len(chunks),
            message="Document processed successfully",
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
