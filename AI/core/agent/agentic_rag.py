"""Agentic RAG module using llama-index FunctionAgent."""

import json
from typing import AsyncGenerator

from llama_index.core.agent import FunctionAgent
from llama_index.core.agent.workflow import (
    AgentOutput,
    AgentStream,
    ToolCallResult,
)
from llama_index.llms.openai_like import OpenAILike
from workflows.events import StopEvent, WorkflowFailedEvent

from core.tools import create_retrieval_tool
from models.chat import ModelConfig, RetrievalConfig, SamplingParams
from models.document import EmbeddingModelConfig
from models.search import RetrievalResult
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

    async def stream(self, query: str) -> AsyncGenerator[tuple[str, str], None]:
        """Stream agent response delta-by-delta.

        Yields token-level deltas as they are generated, including thinking
        content for models that support extended thinking (e.g. DeepSeek-R1,
        Qwen3).  Also emits ``("source", json_array)`` when the agent
        completes a retrieval call, carrying chunk IDs and scores for the
        frontend.

        Args:
            query: User query string.

        Yields:
            ("thinking", text) - model internal reasoning
            ("text", text)    - final response content
            ("source", json)  - retrieved chunk metadata (JSON array)
            ("error", text)   - stream-level error
        """
        try:
            handler = self.workflow.run(user_msg=query)
            async for event in handler.stream_events():
                match event:
                    case AgentStream():
                        if event.thinking_delta:
                            yield ("thinking", event.thinking_delta)
                        if event.delta:
                            yield ("text", event.delta)
                    case AgentOutput():
                        # Some reasoning models (DeepSeek-R1, Qwen3) put the
                        # final response in a ThinkingBlock rather than .content.
                        # Collect text from .content and all blocks.
                        text_parts: list[str] = []
                        if event.response and event.response.content:
                            text_parts.append(str(event.response.content))
                        for block in getattr(event.response, "blocks", []) or []:
                            block_text = getattr(block, "text", "") or getattr(block, "content", "") or ""
                            if isinstance(block_text, str) and block_text.strip():
                                text_parts.append(block_text)
                        if text_parts:
                            yield ("text", "".join(text_parts))
                        else:
                            logger.warning("AgentOutput has no text — nothing yielded to frontend")
                    case ToolCallResult(tool_name="retrieve_documents"):
                        # Extract chunk metadata from the retrieval tool's
                        # return value and forward it to the frontend.
                        raw = event.tool_output.raw_output
                        if isinstance(raw, RetrievalResult) and raw.sources:
                            yield ("source", json.dumps(raw.sources, ensure_ascii=False))
                    case StopEvent():
                        return
                    case WorkflowFailedEvent():
                        logger.error(
                            f"Workflow failed | step={event.step_name} | "
                            f"exception={type(event.exception).__name__}: {event.exception}"
                        )
                        yield ("error", f"LLM call failed in step '{event.step_name}': {event.exception}")
                        return

        except Exception as e:
            logger.error(f"Streaming generation failed: {e}")
            yield ("error", str(e))
