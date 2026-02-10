import asyncio

from llama_index.core import SimpleDirectoryReader, VectorStoreIndex
from llama_index.llms.deepseek import DeepSeek

from core.config import settings
from rag.embeddings import ollama_embedding_model
from rag.vector_store import storage_context
from utils.logger import logger

llm = DeepSeek(
    settings.DEEPSEEK_NAME,
    api_key=settings.DEEPSEEK_API_KEY,
)


async def main():
    documents = SimpleDirectoryReader(
        input_files=["./data/paul_graham_essay.txt"]
    ).load_data()

    logger.info(f"Loaded {len(documents)} documents")

    index = VectorStoreIndex.from_documents(
        documents,
        storage_context=storage_context,
        embed_model=ollama_embedding_model,
    )

    query_engine = index.as_query_engine(llm)

    response = await query_engine.aquery("What did the author learn?")

    logger.success(f"Response: {response}")


if __name__ == "__main__":
    asyncio.run(main())
