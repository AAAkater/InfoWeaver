"""Agentic RAG module using llama-index FunctionAgent."""

from typing import AsyncGenerator

from llama_index.core.agent import FunctionAgent
from llama_index.core.agent.workflow import (
    AgentInput,
    AgentOutput,
    AgentStream,
    ToolCall,
    ToolCallResult,
)
from llama_index.llms.openai_like import OpenAILike
from workflows.events import StopEvent

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

    async def stream(self, query: str) -> AsyncGenerator[str, None]:
        """Stream agent response delta-by-delta.

        Yields token-level deltas as they are generated, including thinking
        content for models that support extended thinking (e.g. DeepSeek-R1,
        Qwen3).

        Args:
            query: User query string.

        Yields:
            Text chunks (deltas, thinking_deltas, or final response text).
        """
        logger.info(f"Streaming agent response for query: {query}")

        try:
            handler = self.workflow.run(user_msg=query)
            async for event in handler.stream_events():
                match event:
                    case AgentInput():
                        logger.debug(f"Agent input | agent={event.current_agent_name} | input={event.input}")
                    case AgentStream():
                        if event.delta:
                            yield event.delta
                        if event.thinking_delta:
                            yield event.thinking_delta
                    case AgentOutput():
                        logger.debug(
                            f"Agent output | response={event.response} "
                            f"| tool_calls={event.tool_calls} | raw={event.raw}"
                        )
                    case ToolCall():
                        logger.debug(f"Tool call | tool_name={event.tool_name} | tool_kwargs={event.tool_kwargs}")
                    case ToolCallResult():
                        logger.debug(
                            f"Tool call result | tool_name={event.tool_name} | tool_kwargs={event.tool_kwargs} | tool_output={event.tool_output}"
                        )
                    case StopEvent():
                        # Graceful termination — nothing more to stream.
                        logger.debug("Streaming completed successfully.")
                        logger.debug(f"Final agent output: {event.result}")
                        return
                    case _:
                        logger.debug(f"Unknown event type: {type(event).__name__}")

        except Exception as e:
            logger.error(f"Streaming generation failed: {e}")
            yield f"\n[Error: {e}]"


def _extract_agent_text(output: AgentOutput) -> str:
    """Extract human-readable text from an AgentOutput.

    Handles ChatMessage objects that may contain TextBlock and/or
    ThinkingBlock content blocks.
    """
    response = output.response
    text_parts: list[str] = []

    if hasattr(response, "blocks") and response.blocks:
        for block in response.blocks:
            if hasattr(block, "text") and block.text:
                text_parts.append(block.text)
            elif hasattr(block, "content") and block.content:
                text_parts.append(block.content)
    elif hasattr(response, "content") and response.content:
        text_parts.append(response.content)
    else:
        text_parts.append(str(response))

    return "\n".join(text_parts)
