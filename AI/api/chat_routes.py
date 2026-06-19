"""Agentic RAG chat API routes."""

import json
from typing import AsyncGenerator

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import StreamingResponse
from sqlalchemy.orm import Session

from core.agent.agentic_rag import AgenticRAG
from db.postgresql_db import get_db_session
from models.chat import ChatRequest
from utils import logger

router = APIRouter(prefix="/chat", tags=["chat"])


@router.post(
    "/chat/stream",
    summary="Agentic RAG chat with streaming",
    description="Chat with an AI agent that can retrieve documents, with streaming response",
)
async def agentic_rag_chat_stream(
    req: ChatRequest,
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
            llm_config=req.llm_config,
            dataset_id=req.dataset_id,
            embedding_config=req.embedding_config,
            retrieval_config=req.retrieval_config,
            system_prompt=req.system_prompt,
        )

        async def generate_stream() -> AsyncGenerator[bytes, None]:
            async for chunk_type, chunk_text in agentic_rag.stream(req.query):
                payload = json.dumps({"type": chunk_type, "content": chunk_text}, ensure_ascii=False)
                yield f"{payload}\n".encode("utf-8")
    except Exception as e:
        logger.error(f"Agentic RAG chat stream failed: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Chat failed: {e}",
        )
        # Stream agent response

    return StreamingResponse(
        generate_stream(),
        media_type="text/event-stream",
    )
