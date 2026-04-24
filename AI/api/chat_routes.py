"""Agentic RAG chat API routes."""

from typing import AsyncGenerator

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import StreamingResponse
from sqlalchemy.orm import Session

from core.agent.agentic_rag import AgenticRAG
from db.postgresql_db import get_db_session
from models.chat import ChatRequest
from utils import logger

router = APIRouter()


@router.post(
    "/chat/stream",
    summary="Agentic RAG chat with streaming",
    description="Chat with an AI agent that can retrieve documents, with streaming response",
)
async def agentic_rag_chat_stream(
    request: ChatRequest,
    db: Session = Depends(get_db_session),
) -> StreamingResponse:
    """
    Agentic RAG chat with streaming response.

    This endpoint:
    1. Creates an agent with retrieval tools
    2. The agent decides when to retrieve documents
    3. Streams the agent's response

    Args:
        request: ChatRequest containing query and configuration
        db: Database session

    Returns:
        StreamingResponse with agent's response

    Raises:
        HTTPException: If chat fails
    """
    try:
        # Create AgenticRAG instance
        agentic_rag = AgenticRAG(
            llm_config=request.llm_config,
            dataset_id=request.dataset_id,
            embedding_config=request.embedding_config,
            retrieval_config=request.retrieval_config,
            system_prompt=request.system_prompt,
        )

        # Stream agent response
        async def generate_stream() -> AsyncGenerator[bytes, None]:
            """Generate streaming response from agent."""
            async for chunk in agentic_rag.stream(request.query):
                yield chunk.encode("utf-8")

        return StreamingResponse(
            generate_stream(),
            media_type="text/plain",
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"Agentic RAG chat stream failed: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Chat failed: {e}",
        )
