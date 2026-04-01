"""Hybrid search module combining dense and sparse vector retrieval."""

from pydantic import BaseModel
from pymilvus import AnnSearchRequest, WeightedRanker

from configs.app_config import settings
from db.milvus_db import client
from utils import logger


class SearchResult(BaseModel):
    """Search result model."""

    id: int
    content: str
    distance: float


def dense_search(
    query_dense_embedding: list[float],
    dataset_id: int,
    limit: int = 10,
    expr: str | None = None,
) -> list[SearchResult]:
    """Search using dense vectors only.

    Args:
        query_dense_embedding: Dense query vector.
        dataset_id: Dataset ID to filter results.
        limit: Number of results to return.
        expr: Additional filter expression.

    Returns:
        List of SearchResult objects.
    """

    final_expr = f"dataset_id == {dataset_id}" + (f" and {expr}" if expr else "")
    search_params = {"metric_type": "IP", "params": {}}

    results = client.search(
        collection_name=settings.MILVUS_COLLECTION_NAME,
        data=[query_dense_embedding],
        anns_field="dense_vector",
        limit=limit,
        filter=final_expr,
        output_fields=["id", "content"],
        param=search_params,
    )

    logger.info(f"Dense search returned {len(results[0]) if results else 0} results")
    return _parse_results(results[0] if results else [])


def sparse_search(
    query_sparse_embedding: dict[int, float],
    dataset_id: int,
    limit: int = 10,
    expr: str | None = None,
) -> list[SearchResult]:
    """Search using sparse vectors only.

    Args:
        query_sparse_embedding: Sparse query vector (dict of index -> value).
        dataset_id: Dataset ID to filter results.
        limit: Number of results to return.
        expr: Additional filter expression.

    Returns:
        List of SearchResult objects.
    """

    final_expr = f"dataset_id == {dataset_id}" + (f" and {expr}" if expr else "")
    search_params = {"metric_type": "IP", "params": {}}

    results = client.search(
        collection_name=settings.MILVUS_COLLECTION_NAME,
        data=[query_sparse_embedding],
        anns_field="sparse_vector",
        limit=limit,
        filter=final_expr,
        output_fields=["id", "content"],
        param=search_params,
    )

    logger.info(f"Sparse search returned {len(results[0]) if results else 0} results")
    return _parse_results(results[0] if results else [])


def hybrid_search(
    query_dense_embedding: list[float],
    query_sparse_embedding: dict[int, float],
    dataset_id: int,
    sparse_weight: float = 1.0,
    dense_weight: float = 1.0,
    limit: int = 10,
    expr: str | None = None,
) -> list[SearchResult]:
    """Hybrid search combining dense and sparse vectors.

    Uses weighted reranking to combine results from both vector types.

    Args:
        query_dense_embedding: Dense query vector.
        query_sparse_embedding: Sparse query vector (dict of index -> value).
        dataset_id: Dataset ID to filter results.
        sparse_weight: Weight for sparse vector search (default: 1.0).
        dense_weight: Weight for dense vector search (default: 1.0).
        limit: Number of results to return.
        expr: Additional filter expression.

    Returns:
        List of SearchResult objects with reranked scores.
    """

    final_expr = f"dataset_id == {dataset_id}" + (f" and {expr}" if expr else "")

    dense_search_params = {"metric_type": "IP", "params": {}}
    dense_req = AnnSearchRequest(
        data=[query_dense_embedding],
        anns_field="dense_vector",
        param=dense_search_params,
        limit=limit,
        expr=final_expr,
    )

    sparse_search_params = {"metric_type": "IP", "params": {}}
    sparse_req = AnnSearchRequest(
        data=[query_sparse_embedding],
        anns_field="sparse_vector",
        param=sparse_search_params,
        limit=limit,
        expr=final_expr,
    )

    rerank = WeightedRanker(sparse_weight, dense_weight)

    results = client.hybrid_search(
        collection_name=settings.MILVUS_COLLECTION_NAME,
        reqs=[sparse_req, dense_req],
        ranker=rerank,
        limit=limit,
        output_fields=["id", "content"],
    )

    logger.info(f"Hybrid search returned {len(results[0]) if results else 0} results")
    return _parse_results(results[0] if results else [])


def _parse_results(raw_results: list[dict]) -> list[SearchResult]:
    """Parse raw Milvus results into SearchResult objects.

    Args:
        raw_results: Raw results from Milvus search.

    Returns:
        List of SearchResult objects.
    """
    search_results: list[SearchResult] = []
    for hit in raw_results:
        entity = hit.get("entity", {})
        search_results.append(
            SearchResult(
                id=entity.get("id", hit.get("id")),
                content=entity.get("content", ""),
                distance=hit.get("distance", 0.0),
            )
        )
    logger.info(f"Parsed {len(search_results)} search results")
    return search_results
