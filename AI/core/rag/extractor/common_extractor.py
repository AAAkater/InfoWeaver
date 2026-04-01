from llama_index.core import SimpleDirectoryReader

from utils import logger

documents = SimpleDirectoryReader(input_dir="./data").load_data()


if __name__ == "__main__":
    logger.info(f"Loaded {len(documents)} documents.")
    for doc in documents:
        content = doc.get_content()
        logger.info(f"Document content preview: \n{content[:100]}")
