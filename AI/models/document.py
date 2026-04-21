"""Pydantic models for document processing API."""

from pydantic import BaseModel, Field


class EmbeddingConfig(BaseModel):
    """Configuration for embedding model."""

    model_name: str = Field(..., description="Embedding model name")
    base_url: str = Field(..., description="Base URL for embedding API")
    api_key: str | None = Field(None, description="API key for embedding (optional for Ollama)")
    provider_type: str = Field(..., description="Embedding provider type: 'openai' or 'ollama'")


class ProcessDocumentRequest(BaseModel):
    """Request model for processing a document from MinIO."""

    file_id: int = Field(..., description="Unique file identifier")
    dataset_id: int = Field(..., description="Dataset identifier for grouping files")
    minio_path: str = Field(..., description="Path to file in MinIO storage")
    embedding_config: EmbeddingConfig = Field(..., description="Embedding model configuration")


class ProcessDocumentResponse(BaseModel):
    """Response model for document processing."""

    file_id: int
    dataset_id: int
    file_name: str
    chunks_count: int
    message: str
