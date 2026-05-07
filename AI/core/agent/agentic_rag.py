"""Agentic RAG module using llama-index AgentWorkflow."""

from typing import AsyncGenerator

from llama_index.core.agent import FunctionAgent
from llama_index.core.agent.workflow import AgentWorkflow
from llama_index.core.tools import FunctionTool
from llama_index.llms.openai_like import OpenAILike

from core.tools import create_retrieval_tool
from models.chat import ModelConfig, RetrievalConfig, SamplingParams
from models.document import EmbeddingModelConfig
from utils import logger


class AgenticRAG:
    """Agentic RAG class for managing agent-based retrieval and generation.

    Encapsulates LLM creation, agent construction, query execution, and response
    streaming in a single class. Uses the modern AgentWorkflow API
    (workflow.run(user_msg=...)) instead of the deprecated positional
    agent.run(...).
    """

    workflow: FunctionAgent
    llm_config: ModelConfig
    dataset_id: int
    embedding_config: EmbeddingModelConfig
    retrieval_config: RetrievalConfig
    system_prompt: str | None
    # --- Constructor --------------------------------------------------------

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
        # Build LLM and agent, then wrap in AgentWorkflow
        sampling = llm_config.sampling_params or SamplingParams()
        llm = OpenAILike(
            model=llm_config.model_name,
            api_base=llm_config.base_url,
            api_key=llm_config.api_key,
            is_chat_model=True,
            is_function_calling_model=True,
            temperature=sampling.temperature,
            max_tokens=sampling.max_tokens,
        )

        default_system_prompt = """
You are a helpful AI assistant with access to a knowledge base.
When answering questions, you should:
1. First use the retrieve_documents tool to search for relevant information
2. Base your answers on the retrieved documents when available
3. If no relevant documents are found, you can use your general knowledge but should mention that the information may not be from the knowledge base
4. Be concise and helpful in your responses"""

        final_system_prompt = system_prompt or default_system_prompt

        self.workflow = FunctionAgent(
            tools=[retrieval_tool],
            llm=llm,
            system_prompt=final_system_prompt,
            verbose=True,
        )

    async def query(self, query: str) -> str:
        """
        Run agent query and return response content.

        Args:
            query: User query string

        Returns:
            Response content string
        """
        logger.info(f"Running agent for query: {query}")
        response = await self.workflow.run(user_msg=query)

        # Extract response content
        return str(response.response) if hasattr(response, "response") else str(response)

    async def stream(self, query: str) -> AsyncGenerator[str, None]:
        """
        Stream agent response for query.

        Args:
            query: User query string

        Yields:
            Response chunks as strings
        """
        logger.info(f"Streaming agent response for query: {query}")

        try:
            response_stream = await self.workflow.run(user_msg=query)

            for response in response_stream:
                if hasattr(response, "delta") and response.delta:
                    yield response.delta
                elif hasattr(response, "response"):
                    content = str(response.response)
                    if content:
                        yield content

        except Exception as e:
            logger.error(f"Streaming generation failed: {e}")
            yield f"\n[Error: {e}]"
