"""Milvus database client initialization."""

from pydantic import BaseModel
from pymilvus import DataType, MilvusClient

from configs.app_config import settings
from utils.logger import logger


class VectorEntity(BaseModel):
    """Entity for hybrid vector insertion (dense + sparse)."""

    dense_vector: list[float]
    sparse_vector: dict[int, float]
    content: str
    dataset_id: int


def _create_database(client: MilvusClient, db_name: str) -> None:
    """Create database if it doesn't exist and switch to it."""
    if db_name == "default":
        return

    databases = client.list_databases()
    if db_name not in databases:
        client.create_database(db_name=db_name)
        logger.info(f"Database '{db_name}' created successfully")
    client.use_database(db_name=db_name)


def _create_collection(client: MilvusClient, collection_name: str) -> None:
    """Create the collection if not exists."""

    if client.has_collection(collection_name):
        logger.info(f"Collection '{collection_name}' already exists")
        return
    # Define schema with both dense and sparse vector fields, along with content and dataset_id
    schema = MilvusClient.create_schema(
        auto_id=True,
        enable_dynamic_field=True,
    )

    schema.add_field(
        field_name="id",
        datatype=DataType.INT64,
        is_primary=True,
        auto_id=True,
    )
    schema.add_field(
        field_name="dense_vector",
        datatype=DataType.FLOAT_VECTOR,
        dim=settings.MILVUS_DIM,
    )
    schema.add_field(field_name="sparse_vector", datatype=DataType.SPARSE_FLOAT_VECTOR)
    schema.add_field(field_name="content", datatype=DataType.VARCHAR, max_length=512)
    schema.add_field(field_name="dataset_id", datatype=DataType.INT64)
    # Create indexes for both dense and sparse vectors, and dataset_id for efficient querying
    index_params = client.prepare_index_params()

    index_params.add_index(
        field_name="dense_vector",
        index_name="dense_vector_index",
        index_type="AUTOINDEX",
        metric_type="IP",
    )

    index_params.add_index(
        field_name="sparse_vector",
        index_name="sparse_inverted_index",
        index_type="SPARSE_INVERTED_INDEX",
        metric_type="IP",
        params={"inverted_index_algo": "DAAT_MAXSCORE"},
    )

    index_params.add_index(field_name="dataset_id", index_type="AUTOINDEX")

    client.create_collection(
        collection_name=collection_name,
        schema=schema,
        index_params=index_params,
    )

    result = client.get_collection_stats(collection_name=collection_name)
    logger.info(f"Collection '{collection_name}' created successfully. Stats: {result}")


def _init_milvus_client(uri: str, db_name: str, collection_name: str) -> MilvusClient:
    """Initialize Milvus client with database and collection setup."""
    client = MilvusClient(uri=uri, timeout=1000)
    _create_database(client, db_name)
    _create_collection(client, collection_name)
    return client


# Export client instance
client = _init_milvus_client(
    uri=settings.MILVUS_URI,
    db_name=settings.MILVUS_DB_NAME,
    collection_name=settings.MILVUS_COLLECTION_NAME,
)
