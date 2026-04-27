"""Search API routes for RAG retrieval and generation."""

from fastapi import APIRouter, HTTPException, status

from core.rag.embedding import OllamaDenseEmbeddingModel, SparseEmbeddingModel
from core.rag.retrieval.search import hybrid_search
from models.search import SearchRequest, SearchResponse, SearchResult
from utils import logger

router = APIRouter(prefix="/search", tags=["search"])


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
