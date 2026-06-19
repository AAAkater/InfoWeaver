from llama_index.core import SimpleDirectoryReader


def extractor(input_files: str) -> SimpleDirectoryReader:
    """Load documents from files or directory.

    Args:
        input_files: List of file paths to load.
    Returns:
        List of loaded documents.
    """

    return SimpleDirectoryReader(input_files=[input_files])
