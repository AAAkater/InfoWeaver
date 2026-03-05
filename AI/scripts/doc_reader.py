import asyncio
import tempfile
from pathlib import Path

from core.rag.splitter.text_splitter import load_and_split_document

from configs.app_config import settings
from core.rag import doc_store, get_text_embeddings
from core.rag.doc_store.document_store import DocumentChunk
from db.minio_client import minio_client
from db.rabbitmq_client import FileUploadMessage, rabbitmq_client
from utils.logger import logger


async def process_document(msg: FileUploadMessage) -> None:
    """
    Process an uploaded document: download, split, embed, and store.

    Args:
        msg: FileUploadMessage containing file information.
    """
    file_name = Path(msg.minio_path).stem
    suffix = Path(msg.minio_path).suffix

    tmp_path: Path | None = None
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp_file:
        try:
            # Step 1: Download file from MinIO
            logger.info(f"Downloading file from MinIO: {msg.minio_path}")
            try:
                file_data = minio_client.download_file(msg.minio_path)
                tmp_file.write(file_data)
                tmp_file.flush()
                tmp_path = Path(tmp_file.name)
            except Exception as e:
                logger.error(f"Failed to download file from MinIO: {e}")
                raise

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
                    return

                logger.info(f"Extracted {len(chunks)} chunks from document: {file_name}")
            except Exception as e:
                logger.error(f"Failed to load and split document: {e}")
                raise

            # Step 3: Generate embeddings
            logger.info(f"Generating embeddings for {len(chunks)} chunks")
            try:
                embeddings = await get_text_embeddings(chunks)
            except Exception as e:
                logger.error(f"Failed to generate embeddings: {e}")
                raise

            # Step 4: Store chunks with embeddings in Milvus
            logger.info(f"Storing {len(chunks)} chunks in Milvus")
            try:
                doc_entities = [
                    DocumentChunk(
                        vector=emb,
                        content=chunk,
                        dataset_id=msg.dataset_id,
                    )
                    for chunk, emb in zip(chunks, embeddings)
                ]

                await doc_store.add_chunks(doc_entities)
            except Exception as e:
                logger.error(f"Failed to store chunks in Milvus: {e}")
                raise

            logger.success(f"Successfully processed document: {file_name} with {len(chunks)} chunks")

        except Exception as e:
            logger.error(f"Error processing document {msg.minio_path}: {e}")
            raise
        finally:
            # Clean up temporary file
            if tmp_path:
                tmp_path.unlink(missing_ok=True)


def on_message(msg: FileUploadMessage) -> None:
    """
    Callback for handling messages from RabbitMQ.

    Args:
        msg: FileUploadMessage containing file upload information.
    """
    logger.info(f"Received file upload message: {msg}")
    asyncio.run(process_document(msg))


def main() -> None:
    """Start the document reader service."""
    logger.info("Starting document reader service...")
    rabbitmq_client.consume(settings.RABBITMQ_QUEUE, on_message)


if __name__ == "__main__":
    main()
