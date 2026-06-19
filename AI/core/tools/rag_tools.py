"""Retrieval tool for document search."""

from typing import Any

from llama_index.core.tools import FunctionTool, ToolMetadata

from core.rag.embedding import OAICompatibleEmbedding, OllamaDenseEmbeddingModel, SparseEmbeddingModel
from core.rag.retrieval.search import hybrid_search
from models.document import EmbeddingModelConfig
from models.search import RetrievalResult
from utils import logger


def _extract_query(tool_kwargs: dict[str, Any]) -> str:
    """Extract query string from tool kwargs, supporting both 'query' and 'input' keys.

    Many LLMs default to 'input' as the parameter name, while our function
    uses 'query'.  This helper makes the tool tolerant of either.

    Args:
        tool_kwargs: The kwargs passed to the tool function.

    Returns:
        The extracted query string.

    Raises:
        ValueError: If neither 'query' nor 'input' is found.
    """
    if "query" in tool_kwargs:
        return str(tool_kwargs["query"])
    if "input" in tool_kwargs:
        return str(tool_kwargs["input"])
    raise ValueError(
        f"Missing query parameter. Received kwargs: {list(tool_kwargs.keys())}. Expected 'query' or 'input'."
    )


async def _debug_callback(raw_output: Any) -> None:
    """Debug callback to log tool call results.

    Called by FunctionTool after execution.  Receives the raw return value of the
    tool function (NOT a ToolOutput object).
    """
    output_str = str(raw_output)
    if len(output_str) > 300:
        preview = output_str[:300] + "..."
    else:
        preview = output_str
    logger.info(f"Tool 'retrieve_documents' returned ({len(output_str)} chars): {preview}")


def create_retrieval_tool(
    dataset_id: int,
    embedding_config: EmbeddingModelConfig,
    top_k: int = 10,
) -> FunctionTool:
    """
    Create a retrieval tool for the agent to search documents.

    Args:
        dataset_id: Dataset ID to search in
        embedding_config: Embedding model configuration
        top_k: Number of results to retrieve

    Returns:
        FunctionTool for document retrieval
    """

    async def retrieve_documents(query: str = "", **kwargs: Any) -> RetrievalResult:
        """
        Retrieve relevant documents from the knowledge base.

        Args:
            query: The search query string.
            **kwargs: Catch-all for LLMs that send 'input' instead of 'query'.

        Returns:
            RetrievalResult with formatted text for the LLM (.text) and
            structured chunk metadata for the frontend (.sources).
        """
        # Support both 'query' and 'input' parameter names
        try:
            actual_query = _extract_query({"query": query, **kwargs})
        except ValueError as e:
            logger.error(f"Failed to extract query from tool call: {e}")
            return RetrievalResult(text=f"[Error] Invalid tool call parameters: {e}")

        logger.info(f"Agent retrieving documents for query: {actual_query}")

        try:
            # Initialize embedding models based on config
            sparse_embedding_model = SparseEmbeddingModel()

            if embedding_config.provider_type == "ollama":
                dense_embedding_model = OllamaDenseEmbeddingModel(
                    model_name=embedding_config.model_name,
                    base_url=embedding_config.base_url,
                )
            else:
                dense_embedding_model = OAICompatibleEmbedding(
                    model_name=embedding_config.model_name,
                    base_url=embedding_config.base_url,
                    api_key=embedding_config.api_key or "",
                )

            # Generate embeddings
            sparse_embeddings = await sparse_embedding_model.get_embeddings([actual_query])
            dense_embeddings = await dense_embedding_model.get_embeddings([actual_query])

            # Perform hybrid search
            search_results = hybrid_search(
                query_dense_embedding=dense_embeddings[0],
                query_sparse_embedding=sparse_embeddings[0],
                dataset_id=dataset_id,
                limit=top_k,
            )

            # Build sources for frontend (chunk IDs, scores, previews)
            sources: list[dict] = []
            for r in search_results:
                preview = r.content[:200] + "..." if len(r.content) > 200 else r.content
                sources.append(
                    {
                        "id": r.id,
                        "content": preview,
                        "score": round(r.distance, 4),
                    }
                )

            if not search_results:
                logger.debug(f"No relevant documents found for query: {actual_query}")
                return RetrievalResult(text="No relevant documents found.")

            logger.debug(f"Retrieved {len(search_results)} documents for query: {actual_query}")
            formatted_results: list[str] = []
            for i, result in enumerate(search_results, 1):
                content_preview = result.content[:200] + "..." if len(result.content) > 200 else result.content
                logger.debug(
                    f"[Document {i}/{len(search_results)}] distance={result.distance:.4f}, preview: {content_preview}"
                )
                formatted_results.append(f"[Document {i}]\n{result.content}\n")

            return RetrievalResult(
                text="\n".join(formatted_results),
                sources=sources,
            )

        except Exception as e:
            logger.exception(f"Hybrid search failed for query '{actual_query}': {e}")
            return RetrievalResult(text=f"[Error] Search failed: {e}")

    return FunctionTool(
        async_fn=retrieve_documents,
        metadata=ToolMetadata(
            name="retrieve_documents",
            description=(
                "Retrieve relevant documents from the knowledge base. "
                "Use this tool when you need to search for information to answer the user's question. "
                "Parameters: query (string, required) - the search query."
            ),
        ),
        async_callback=_debug_callback,
    )
