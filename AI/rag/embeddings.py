from llama_index.embeddings.ollama import OllamaEmbedding

ollama_embedding_model = OllamaEmbedding(
    model_name="qwen3-embedding:0.6b",
    base_url="http://localhost:11434",
)
