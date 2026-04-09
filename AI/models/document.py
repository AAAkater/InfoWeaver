"""Pydantic models for document processing API."""

from pydantic import BaseModel, Field


class ProcessDocumentRequest(BaseModel):
    """Request model for processing a document from MinIO."""

    file_id: int = Field(..., description="Unique file identifier")
    dataset_id: int = Field(..., description="Dataset identifier for grouping files")
    minio_path: str = Field(..., description="Path to file in MinIO storage")


class ProcessDocumentResponse(BaseModel):
    """Response model for document processing."""

    file_id: int
    dataset_id: int
    file_name: str
    chunks_count: int
    message: str
