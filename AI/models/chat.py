"""Pydantic models for chat API."""

from pydantic import BaseModel, Field

from models.document import EmbeddingConfig


class SamplingConfig(BaseModel):
    """Configuration for LLM sampling/generation parameters."""

    temperature: float = Field(default=0.7, ge=0.0, le=2.0, description="Sampling temperature (0.0-2.0)")
    top_p: float = Field(default=0.9, ge=0.0, le=1.0, description="Nucleus sampling probability (0.0-1.0)")
    max_tokens: int | None = Field(default=None, ge=1, description="Maximum tokens to generate (optional)")
    presence_penalty: float = Field(default=0.0, ge=-2.0, le=2.0, description="Presence penalty (-2.0 to 2.0)")
    frequency_penalty: float = Field(default=0.0, ge=-2.0, le=2.0, description="Frequency penalty (-2.0 to 2.0)")


class RetrievalConfig(BaseModel):
    """Configuration for retrieval/search parameters."""

    top_k: int = Field(default=10, ge=1, description="Number of search results to retrieve")


class ModelConfig(BaseModel):
    """Configuration for LLM model."""

    model_name: str = Field(..., description="Model name to use for generation")
    api_key: str = Field(..., description="API key for the model provider")
    base_url: str = Field(..., description="Base URL for the model API")
    provider_type: str = Field(..., description="Provider type: 'openai', 'anthropic', 'ollama', etc")
    sampling_config: SamplingConfig = Field(
        default_factory=SamplingConfig, description="Sampling/generation parameters"
    )


class ChatRequest(BaseModel):
    """Request model for agentic RAG chat."""

    query: str = Field(..., description="User query text")
    dataset_id: int = Field(..., description="Dataset ID to search in")
    session_id: int | None = Field(None, description="Chat session ID (optional for new session)")
    llm_config: ModelConfig = Field(..., description="LLM model configuration")
    embedding_config: EmbeddingConfig = Field(..., description="Embedding model configuration")
    retrieval_config: RetrievalConfig = Field(
        default_factory=RetrievalConfig, description="Retrieval/search parameters"
    )
    system_prompt: str | None = Field(
        None,
        description="Custom system prompt (optional)",
    )


class ChatMessage(BaseModel):
    """Single chat message."""

    role: str = Field(..., description="Message role: 'user', 'assistant', or 'system'")
    content: str = Field(..., description="Message content")


class ChatResponse(BaseModel):
    """Response model for chat (non-streaming)."""

    session_id: int = Field(..., description="Chat session ID")
    message: ChatMessage = Field(..., description="Assistant response message")
    retrieved_chunks: list[int] = Field(default_factory=list, description="IDs of retrieved chunks")


class ChatSessionCreate(BaseModel):
    """Request model for creating a new chat session."""

    title: str = Field(..., description="Session title")
    owner_id: int = Field(..., description="Owner user ID")


class ChatSessionResponse(BaseModel):
    """Response model for chat session."""

    id: int
    title: str
    owner_id: int
