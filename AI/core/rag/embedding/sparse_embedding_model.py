from FlagEmbedding import BGEM3FlagModel


class SparseEmbeddingModel:
    def __init__(self, model_name: str = "BAAI/bge-m3"):
        self.embedding_model = BGEM3FlagModel(
            model_name,
            use_fp16=False,
            devices="cpu",
        )

    async def get_embeddings(self, texts: list[str]) -> list[dict[int, float]]:
        """Get sparse embeddings for a list of texts.

        Returns:
            A list of sparse vectors, each represented as a dict mapping token indices to weights.
        """
        if not texts:
            return []

        # BGEM3FlagModel.encode returns a dict with 'lexical_weights' for sparse vectors
        # lexical_weights contains token indices mapped to their weights
        result = self.embedding_model.encode(
            texts,
            return_dense=False,
            return_sparse=True,
            return_colbert_vecs=False,
        )

        # The 'lexical_weights' key contains sparse vectors as list[dict[int, float]]
        # Each dict maps token indices to their weights
        sparse_vectors: list[dict[int, float]] = result["lexical_weights"]

        return sparse_vectors
