"""Milvus database client initialization."""

from pydantic import BaseModel
from pymilvus import DataType, MilvusClient

from configs.app_config import settings
from utils import logger


class VectorEntity(BaseModel):
    """Entity for hybrid vector insertion (dense + sparse)."""

    dense_vector: list[float]
    sparse_vector: dict[int, float]
    content: str
    dataset_id: int


class MilvusDB:
    """Milvus database client wrapper for vector operations."""

    def __init__(
        self,
        uri: str,
        collection_name: str,
        db_name: str = "default",
    ):
        self.client = MilvusClient(uri=uri, db_name=db_name, timeout=1000)
        self.collection_name = collection_name
        self._create_collection()

    def _create_collection(self) -> None:
        """Create the collection if not exists."""

        if self.client.has_collection(self.collection_name):
            logger.info(f"Collection '{self.collection_name}' already exists")
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
        index_params = self.client.prepare_index_params()

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

        self.client.create_collection(
            collection_name=self.collection_name,
            schema=schema,
            index_params=index_params,
        )

        result = self.client.get_collection_stats(collection_name=self.collection_name)
        logger.info(f"Collection '{self.collection_name}' created successfully. Stats: {result}")

    def insert_entities(self, entities: list[VectorEntity]) -> None:
        """Insert vector entities into the collection.

        Args:
            entities: List of VectorEntity objects to insert.
        """
        self.client.insert(
            collection_name=self.collection_name,
            data=[data.model_dump() for data in entities],
        )
        logger.info(f"Inserted {len(entities)} records into collection '{self.collection_name}'")

    def delete_entities_by_ids(self, entity_ids: list[int]) -> None:
        """Delete entities by their IDs.

        Args:
            entity_ids: List of entity IDs to delete.
        """
        self.client.delete(
            collection_name=self.collection_name,
            ids=entity_ids,
        )
        logger.info(f"Deleted entities with IDs {entity_ids} from collection '{self.collection_name}'")

    def delete_entities_by_filter(self, expr: str) -> None:
        """Delete entities by filter expression.

        Args:
            expr: Filter expression for deletion.
        """
        self.client.delete(
            collection_name=self.collection_name,
            filter_expression=expr,
        )
        logger.info(f"Deleted entities with filter expression '{expr}' from collection '{self.collection_name}'")

    def update_entities_by_ids(self, entity_ids: list[int], updated_data: list[VectorEntity]) -> None:
        """Update entities by their IDs.

        Args:
            entity_ids: List of entity IDs to update.
            updated_data: List of VectorEntity objects with updated data.
        """
        self.client.upsert(
            collection_name=self.collection_name,
            ids=entity_ids,
            data=[data.model_dump() for data in updated_data],
        )
        logger.info(f"Updated entities with IDs {entity_ids} in collection '{self.collection_name}'")


client = MilvusDB(
    uri=settings.MILVUS_URI,
    collection_name=settings.MILVUS_COLLECTION_NAME,
    db_name=settings.MILVUS_DB_NAME,
)
