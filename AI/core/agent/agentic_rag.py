"""Agentic RAG module using llama-index ReActAgent."""

from typing import AsyncGenerator

from llama_index.core.agent import FunctionAgent, ReActAgent
from llama_index.core.tools import FunctionTool
from llama_index.llms.openai_like import OpenAILike

from core.tools import create_retrieval_tool
from models.chat import ModelConfig, RetrievalConfig
from models.document import EmbeddingModelConfig
from utils import logger


def create_llm(llm_config: ModelConfig) -> OpenAILike:
    """
    Create LLM instance from configuration.

    Args:
        llm_config: LLM model configuration

    Returns:
        OpenAILike LLM instance
    """
    sampling = llm_config.sampling_params
    return OpenAILike(
        model=llm_config.model_name,
        api_base=llm_config.base_url,
        api_key=llm_config.api_key,
        is_chat_model=True,
        temperature=sampling.temperature,
        max_tokens=sampling.max_tokens,
    )


def create_rag_agent(
    llm_config: ModelConfig,
    tools: list[FunctionTool],
    system_prompt: str | None = None,
) -> FunctionAgent:
    """
    Create a llama-index FunctionAgent with retrieval tools.

    Args:
        llm_config: LLM model configuration
        tools: List of tools for the agent
        system_prompt: Custom system prompt (optional)

    Returns:
        ReActAgent instance
    """
    # Create LLM using OpenAI-compatible interface
    llm = create_llm(llm_config)

    # Default system prompt for RAG agent
    default_system_prompt = """You are a helpful AI assistant with access to a knowledge base.
When answering questions, you should:
1. First use the retrieve_documents tool to search for relevant information
2. Base your answers on the retrieved documents when available
3. If no relevant documents are found, you can use your general knowledge but should mention that the information may not be from the knowledge base
4. Be concise and helpful in your responses"""

    final_system_prompt = system_prompt or default_system_prompt

    # Create ReActAgent from tools
    agent = FunctionAgent(
        tools=tools,
        llm=llm,
        system_prompt=final_system_prompt,
        verbose=True,
    )

    return agent


async def run_agent_query(agent: FunctionAgent, query: str) -> str:
    """
    Run agent query and return response content.

    Args:
        agent: ReActAgent instance
        query: User query string

    Returns:
        Response content string
    """
    logger.info(f"Running agent for query: {query}")
    response = await agent.run(query)

    # Extract response content
    response_content = str(response.response) if hasattr(response, "response") else str(response)
    return response_content


async def stream_agent_response(agent: FunctionAgent, query: str) -> AsyncGenerator[str, None]:
    """
    Stream agent response for query.

    Args:
        agent: FunctionAgent instance
        query: User query string

    Yields:
        Response chunks as strings
    """
    logger.info(f"Streaming agent response for query: {query}")

    try:
        # Use streaming chat
        response_stream = await agent.run(query, stream=True)

        for response in response_stream:
            if hasattr(response, "delta") and response.delta:
                yield response.delta
            elif hasattr(response, "response"):
                # For non-delta responses, yield the content
                content = str(response.response)
                if content:
                    yield content

    except Exception as e:
        logger.error(f"Streaming generation failed: {e}")
        yield f"\n[Error: {e}]"


class AgenticRAG:
    """Agentic RAG class for managing agent-based retrieval and generation."""

    agent: FunctionAgent
    llm_config: ModelConfig
    dataset_id: int
    embedding_config: EmbeddingModelConfig
    retrieval_config: RetrievalConfig
    system_prompt: str | None

    def __init__(
        self,
        llm_config: ModelConfig,
        dataset_id: int,
        embedding_config: EmbeddingModelConfig,
        retrieval_config: RetrievalConfig | None = None,
        system_prompt: str | None = None,
    ) -> None:
        """
        Initialize AgenticRAG.

        Args:
            llm_config: LLM model configuration
            dataset_id: Dataset ID to search in
            embedding_config: Embedding model configuration
            retrieval_config: Retrieval/search parameters (optional, defaults to top_k=10)
            system_prompt: Custom system prompt (optional)
        """
        self.llm_config = llm_config
        self.dataset_id = dataset_id
        self.embedding_config = embedding_config
        self.retrieval_config = retrieval_config or RetrievalConfig()
        self.system_prompt = system_prompt

        # Create retrieval tool
        retrieval_tool = create_retrieval_tool(
            dataset_id=dataset_id,
            embedding_config=embedding_config,
            top_k=self.retrieval_config.top_k,
        )

        # Create agent
        self.agent = create_rag_agent(
            llm_config=llm_config,
            tools=[retrieval_tool],
            system_prompt=system_prompt,
        )

    async def query(self, query: str) -> str:
        """
        Run query and return response.

        Args:
            query: User query string

        Returns:
            Response content string
        """
        return await run_agent_query(self.agent, query)

    async def stream(self, query: str) -> AsyncGenerator[str, None]:
        """
        Stream response for query.

        Args:
            query: User query string

        Yields:
            Response chunks as strings
        """
        async for chunk in stream_agent_response(self.agent, query):
            yield chunk
