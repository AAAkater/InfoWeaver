"""Retrieval tool for document search."""

from llama_index.core.tools import FunctionTool, ToolMetadata

from core.rag.embedding import OAICompatibleEmbedding, OllamaDenseEmbeddingModel, SparseEmbeddingModel
from core.rag.retrieval.search import hybrid_search
from models.document import EmbeddingModelConfig
from utils import logger


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

    async def retrieve_documents(query: str) -> str:
        """
        Retrieve relevant documents from the knowledge base.

        Args:
            query: The search query string

        Returns:
            Concatenated content of retrieved documents
        """
        logger.info(f"Agent retrieving documents for query: {query}")

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
        sparse_embeddings = await sparse_embedding_model.get_embeddings([query])
        dense_embeddings = await dense_embedding_model.get_embeddings([query])

        # Perform hybrid search
        search_results = hybrid_search(
            query_dense_embedding=dense_embeddings[0],
            query_sparse_embedding=sparse_embeddings[0],
            dataset_id=dataset_id,
            limit=top_k,
        )

        # Format results for agent
        if not search_results:
            return "No relevant documents found."

        formatted_results = []
        for i, result in enumerate(search_results, 1):
            formatted_results.append(f"[Document {i}]\n{result.content}\n")

        return "\n".join(formatted_results)

    return FunctionTool(
        fn=retrieve_documents,
        metadata=ToolMetadata(
            name="retrieve_documents",
            description="Retrieve relevant documents from the knowledge base. Use this tool when you need to search for information to answer the user's question.",
        ),
    )
