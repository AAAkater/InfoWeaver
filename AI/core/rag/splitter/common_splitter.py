from llama_index.core.node_parser import SentenceSplitter


def splitter(
    chunk_size: int = 512,
    chunk_overlap: int = 50,
) -> SentenceSplitter:

    splitter = SentenceSplitter(
        chunk_size=chunk_size,
        chunk_overlap=chunk_overlap,
    )
    return splitter
