"""Pydantic models for search API."""

from dataclasses import dataclass, field

from pydantic import BaseModel, Field


class SearchRequest(BaseModel):
    """Request model for RAG search."""

    query: str = Field(..., description="User query text")
    dataset_id: int = Field(..., description="Dataset ID to search in")
    provider_id: int = Field(..., description="Provider ID for LLM")
    model_name: str = Field(..., description="Model name to use for generation")
    top_k: int = Field(default=10, description="Number of search results to retrieve")


class SearchResult(BaseModel):
    """Single search result."""

    id: int
    content: str
    distance: float


class SearchResponse(BaseModel):
    """Response model for search results."""

    results: list[SearchResult]
    total: int


@dataclass
class RetrievalResult:
    """Wrapper returned by retrieve_documents tool.

    __str__ returns LLM-friendly formatted text so the agent can read it,
    while ``sources`` carries structured chunk metadata for the frontend:

        {"id": <chunk_id>, "content": "<preview>", "score": <similarity>}
    """

    text: str
    sources: list[dict] = field(default_factory=list)

    def __str__(self) -> str:
        return self.text
