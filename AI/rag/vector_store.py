from llama_index.core import StorageContext
from llama_index.vector_stores.milvus import MilvusVectorStore

from core.config import settings


def init_vector_store():
    return MilvusVectorStore(
        uri=settings.MILVUS_URI,
        dim=settings.MILVUS_DIM,
        collection_name="test_collection",
        embedding_field="embedding",
        id_field="id",
        similarity_metric="COSINE",
        consistency_level="Strong",
        overwrite=True,  # To overwrite the collection if it already exists
    )


milvus_vector_store = init_vector_store()


storage_context = StorageContext.from_defaults(vector_store=milvus_vector_store)
