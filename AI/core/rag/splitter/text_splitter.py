from pathlib import Path

from langchain_community.document_loaders import UnstructuredFileLoader
from langchain_text_splitters import RecursiveCharacterTextSplitter

from utils.logger import logger


def _get_splitter(chunk_size: int, chunk_overlap: int) -> RecursiveCharacterTextSplitter:
    """Create a text splitter with the given parameters."""
    return RecursiveCharacterTextSplitter.from_tiktoken_encoder(
        chunk_size=chunk_size,
        chunk_overlap=chunk_overlap,
    )


def split_documents(
    documents: list[str],
    chunk_size: int = 512,
    chunk_overlap: int = 50,
) -> list[str]:
    """
    Split multiple documents into chunks.

    Args:
        documents: List of document texts to split.
        chunk_size: Maximum size of each chunk in tokens.
        chunk_overlap: Number of tokens to overlap between chunks.

    Returns:
        list[str]: List of all text chunks from all documents.
    """
    all_chunks: list[str] = []
    splitter = _get_splitter(chunk_size, chunk_overlap)

    for i, doc in enumerate(documents):
        chunks = splitter.split_text(doc)
        logger.info(f"Document {i + 1}: split into {len(chunks)} chunks")
        all_chunks.extend(chunks)

    logger.info(f"Total chunks from all documents: {len(all_chunks)}")
    return all_chunks


def load_document(file_path: str | Path) -> list[str]:
    """
    Load document from file path using UnstructuredFileLoader.

    Args:
        file_path: Path to the document file.

    Returns:
        list[str]: List of document text contents.
    """
    file_path = Path(file_path)
    if not file_path.exists():
        raise FileNotFoundError(f"File not found: {file_path}")

    loader = UnstructuredFileLoader(str(file_path))
    documents = loader.load()

    logger.info(f"Loaded {len(documents)} documents from {file_path}")
    return [doc.page_content for doc in documents]


def load_and_split_document(
    file_path: str | Path,
    chunk_size: int = 512,
    chunk_overlap: int = 50,
) -> list[str]:
    """
    Load a document and split it into chunks.

    This is a convenience function that combines load_document and split_documents.

    Args:
        file_path: Path to the document file.
        chunk_size: Maximum size of each chunk in tokens.
        chunk_overlap: Number of tokens to overlap between chunks.

    Returns:
        list[str]: List of text chunks.
    """
    documents = load_document(file_path)
    return split_documents(documents, chunk_size, chunk_overlap)


def split_text(
    text: str,
    chunk_size: int = 512,
    chunk_overlap: int = 50,
) -> list[str]:
    """
    Split text into chunks using RecursiveCharacterTextSplitter.

    Args:
        text: Text to split.
        chunk_size: Maximum size of each chunk in tokens.
        chunk_overlap: Number of tokens to overlap between chunks.

    Returns:
        list[str]: List of text chunks.
    """
    splitter = _get_splitter(chunk_size, chunk_overlap)
    chunks = splitter.split_text(text)
    logger.info(f"Split text into {len(chunks)} chunks")
    return chunks
