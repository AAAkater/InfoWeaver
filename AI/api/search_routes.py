"""Search API routes for RAG retrieval and generation."""

from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.responses import StreamingResponse
from sqlalchemy import select
from sqlalchemy.orm import Session

from core.rag.embedding import OAICompatibleEmbedding, OllamaDenseEmbeddingModel, SparseEmbeddingModel
from core.rag.retrieval.search import hybrid_search
from db.postgresql_db import get_db_session
from models.search import SearchRequest, SearchResponse, SearchResult
from models.tables import Provider
from utils import logger

router = APIRouter()


@router.post(
    "/search",
    response_model=SearchResponse,
    status_code=status.HTTP_200_OK,
    summary="Search documents",
    description="Search documents using hybrid retrieval (dense + sparse vectors)",
)
async def search_documents(request: SearchRequest) -> SearchResponse:
    """
    Search documents using hybrid retrieval.

    This endpoint:
    1. Generates dense and sparse embeddings for the query
    2. Performs hybrid search in Milvus
    3. Returns search results

    Args:
        request: SearchRequest containing query and search parameters

    Returns:
        SearchResponse with search results

    Raises:
        HTTPException: If search fails
    """
    try:
        # Generate embeddings for the query
        logger.info(f"Generating embeddings for query: {request.query}")

        # Use default embedding models for search
        sparse_embedding_model = SparseEmbeddingModel()
        dense_embedding_model = OllamaDenseEmbeddingModel(
            model_name="qwen3-embedding:0.6b",
            base_url="http://localhost:11434",
        )

        sparse_embeddings = await sparse_embedding_model.get_embeddings([request.query])
        dense_embeddings = await dense_embedding_model.get_embeddings([request.query])

        # Perform hybrid search
        logger.info(f"Performing hybrid search in dataset {request.dataset_id}")
        search_results = hybrid_search(
            query_dense_embedding=dense_embeddings[0],
            query_sparse_embedding=sparse_embeddings[0],
            dataset_id=request.dataset_id,
            limit=request.top_k,
        )

        # Convert to response format
        results = [
            SearchResult(
                id=result.id,
                content=result.content,
                distance=result.distance,
            )
            for result in search_results
        ]

        logger.info(f"Search completed with {len(results)} results")

        return SearchResponse(
            results=results,
            total=len(results),
        )

    except Exception as e:
        logger.error(f"Search failed: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Search failed: {e}",
        )


@router.post(
    "/rag",
    summary="RAG query with streaming response",
    description="Retrieve documents and generate streaming response using LLM",
)
async def rag_query(
    request: SearchRequest,
    db: Session = Depends(get_db_session),
) -> StreamingResponse:
    """
    RAG query with streaming response.

    This endpoint:
    1. Retrieves provider configuration from database
    2. Generates embeddings for the query
    3. Performs hybrid search
    4. Constructs context from search results
    5. Streams LLM response

    Args:
        request: SearchRequest containing query and LLM parameters
        db: Database session

    Returns:
        StreamingResponse with LLM generated text

    Raises:
        HTTPException: If provider not found or generation fails
    """
    try:
        # Step 1: Get provider configuration from database
        logger.info(f"Fetching provider with ID {request.provider_id}")
        stmt = select(Provider).where(Provider.id == request.provider_id, Provider.deleted_at.is_(None))
        provider = db.execute(stmt).scalar_one_or_none()

        if not provider:
            logger.error(f"Provider with ID {request.provider_id} not found")
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Provider with ID {request.provider_id} not found",
            )

        logger.info(f"Found provider: {provider.name} (kind: {provider.kind})")

        # Step 2: Generate embeddings for the query
        logger.info(f"Generating embeddings for query: {request.query}")
        sparse_embedding_model = SparseEmbeddingModel()

        # Use appropriate embedding model based on provider kind
        if provider.kind == "ollama":
            dense_embedding_model = OllamaDenseEmbeddingModel(
                model_name="qwen3-embedding:0.6b",
                base_url=provider.base_url,
            )
        else:
            dense_embedding_model = OAICompatibleEmbedding(
                model_name="text-embedding-3-small",
                base_url=provider.base_url,
            )

        sparse_embeddings = await sparse_embedding_model.get_embeddings([request.query])
        dense_embeddings = await dense_embedding_model.get_embeddings([request.query])

        # Step 3: Perform hybrid search
        logger.info(f"Performing hybrid search in dataset {request.dataset_id}")
        search_results = hybrid_search(
            query_dense_embedding=dense_embeddings[0],
            query_sparse_embedding=sparse_embeddings[0],
            dataset_id=request.dataset_id,
            limit=request.top_k,
        )

        # Step 4: Construct context from search results
        context = "\n\n".join([result.content for result in search_results])
        logger.info(f"Constructed context with {len(search_results)} documents")

        # Step 5: Stream LLM response
        logger.info(f"Starting streaming generation with model {request.model_name}")

        # Use OpenAI-compatible LLM for all providers
        from llama_index.llms.openai_like import OpenAILike

        llm = OpenAILike(
            model=request.model_name,
            api_base=provider.base_url,
            api_key=provider.api_key,
            is_chat_model=True,
        )

        # Construct prompt with context
        prompt = f"""Based on the following context, please answer the user's question.

Context:
{context}

Question: {request.query}

Answer:"""

        # Stream the response
        async def generate_stream():
            """Generate streaming response from LLM."""
            try:
                response_stream = await llm.astream_complete(prompt)
                async for response in response_stream:
                    if response.delta:
                        yield response.delta.encode("utf-8")

            except Exception as e:
                logger.error(f"Streaming generation failed: {e}")
                yield f"\n\nError: Failed to generate response - {e}".encode("utf-8")

        return StreamingResponse(
            generate_stream(),
            media_type="text/event-stream",
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.error(f"RAG query failed: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"RAG query failed: {e}",
        )
