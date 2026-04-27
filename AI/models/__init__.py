"""Models package."""

from pydantic import BaseModel, Field

from models.document import ProcessDocumentRequest, ProcessDocumentResponse
from models.search import SearchRequest, SearchResponse, SearchResult

__all__ = [
    "ProcessDocumentRequest",
    "ProcessDocumentResponse",
    "SearchRequest",
    "SearchResponse",
    "SearchResult",
]


class APIResponse[T](BaseModel):
    """Unified API response wrapper."""

    code: int = Field(default=0, description="Status code, 0 = success")
    msg: str = Field(default="success", description="Response message")
    data: T | None = Field(default=None, description="Response data payload")
