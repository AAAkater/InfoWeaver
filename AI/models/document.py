"""Pydantic models for document processing API."""

from pydantic import BaseModel, Field


class EmbeddingModelConfig(BaseModel):
    """Configuration for embedding model."""

    model_name: str = Field(..., description="Embedding model name")
    base_url: str = Field(..., description="Base URL for embedding API")
    api_key: str | None = Field(None, description="API key for embedding (optional for Ollama)")
    provider_type: str = Field(..., description="Embedding provider type: 'openai' or 'ollama'")
    embed_type: str = Field(
        default="hybrid",
        description="Embedding type: 'dense' (dense only), 'sparse' (sparse only), 'hybrid' (both)",
    )


# ---- Request models ----


class ProcessDocumentRequest(BaseModel):
    """Request model for processing a document from MinIO."""

    file_id: int = Field(..., description="Unique file identifier")
    dataset_id: int = Field(..., description="Dataset identifier for grouping files")
    minio_path: str = Field(..., description="Path to file in MinIO storage")
    embedding_config: EmbeddingModelConfig = Field(..., description="Embedding model configuration")
    chunk_size: int = Field(default=512, description="Maximum size of each chunk in tokens")
    chunk_overlap: int = Field(default=50, description="Number of tokens to overlap between chunks")


class SplitDocumentRequest(BaseModel):
    """Request model for splitting a document from MinIO into chunks."""

    file_id: int = Field(..., description="Unique file identifier")
    dataset_id: int = Field(..., description="Dataset identifier for grouping files")
    minio_path: str = Field(..., description="Path to file in MinIO storage")
    chunk_size: int = Field(default=512, description="Maximum size of each chunk in tokens")
    chunk_overlap: int = Field(default=50, description="Number of tokens to overlap between chunks")


class EmbedChunksRequest(BaseModel):
    """Request model for embedding and storing document chunks."""

    chunk_ids: list[int] = Field(..., description="List of chunk IDs from PostgreSQL to embed")
    embedding_config: EmbeddingModelConfig = Field(..., description="Embedding model configuration")


# ---- Data payload models ----


class SplitDocumentData(BaseModel):
    """Data payload for document splitting response."""

    file_id: int
    dataset_id: int
    file_name: str
    chunks_count: int


class EmbedChunksData(BaseModel):
    """Data payload for chunk embedding response."""

    chunk_ids: list[int]
    chunks_count: int


class ProcessDocumentData(BaseModel):
    """Data payload for document processing response."""

    file_id: int
    dataset_id: int
    file_name: str
    chunks_count: int


# Response types are formed as APIResponse[PayloadData] at call sites.
# Example: APIResponse[SplitDocumentData], APIResponse[ProcessDocumentData]
