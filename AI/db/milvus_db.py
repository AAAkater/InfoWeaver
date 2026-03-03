from pydantic import BaseModel
from pymilvus import DataType, MilvusClient

from configs.app_config import settings
from utils.logger import logger


class VectorEntity(BaseModel):
    id: int
    vector: list[float]
    content: str
    dataset_id: int


class MilvusDB:
    def __init__(self, uri: str, collection_name: str, db_name: str = "default"):
        self.client = MilvusClient(uri=uri, timeout=1000)
        self.db_name = db_name
        self.collection_name = collection_name
        if db_name != "default":
            databases = self.client.list_databases()
            if db_name not in databases:
                self.client.create_database(db_name=db_name)
                logger.info(f"Database '{db_name}' created successfully")
            self.client.use_database(db_name=db_name)
        self.create_collection()

    def create_collection(self):
        """Create the collection if not exists."""
        schema = MilvusClient.create_schema(
            auto_id=False,
            enable_dynamic_field=True,
        )

        schema.add_field(field_name="id", datatype=DataType.INT64, is_primary=True)
        schema.add_field(
            field_name="vector",
            datatype=DataType.FLOAT_VECTOR,
            dim=settings.MILVUS_DIM,
        )
        schema.add_field(field_name="content", datatype=DataType.VARCHAR, max_length=512)
        schema.add_field(field_name="dataset_id", datatype=DataType.INT64)

        index_params = self.client.prepare_index_params()
        index_params.add_index(field_name="id", index_type="AUTOINDEX")
        index_params.add_index(field_name="vector", index_type="AUTOINDEX", metric_type="COSINE")
        index_params.add_index(field_name="dataset_id", index_type="AUTOINDEX")

        self.client.create_collection(
            collection_name=self.collection_name,
            schema=schema,
            index_params=index_params,
        )

        result = self.client.get_collection_stats(collection_name=self.collection_name)
        logger.info(f"Collection '{self.collection_name}' created successfully. Stats: {result}")

    @property
    async def get_collection_list(self):
        return await self.client.list_collections()

    def insert_entities(self, new_vector_datas: list[VectorEntity]):
        self.client.insert(
            collection_name=self.collection_name,
            data=[data.model_dump() for data in new_vector_datas],
        )
        logger.info(f"Inserted {len(new_vector_datas)} records into collection '{self.collection_name}'")

    def delete_entities_by_id(self, entity_ids: list[int]):
        self.client.delete(
            collection_name=self.collection_name,
            ids=entity_ids,
        )
        logger.info(f"Deleted entities with IDs {entity_ids} from collection '{self.collection_name}'")

    def delete_entities_by_filter(self, expr: str):
        self.client.delete(collection_name=self.collection_name, filter_expression=expr)
        logger.info(f"Deleted entities with filter expression '{expr}' from collection '{self.collection_name}'")

    def search_entities(
        self,
        query_vector: list[float],
        dataset_id: int,
        top_k: int = 10,
    ) -> list[dict]:
        """Search for similar vectors in the collection.

        Args:
            query_vector: The query vector to search for.
            top_k: Number of results to return. Defaults to 10.
            dataset_id: Optional dataset ID to filter results.

        Returns:
            List of search results with distance and entity data.
        """

        results = self.client.search(
            collection_name=self.collection_name,
            data=[query_vector],
            limit=top_k,
            filter=f"dataset_id == {dataset_id}",
            output_fields=["id", "content"],
        )

        logger.info(f"Search returned {len(results[0])} results")
        return results[0] if results else []


milvus_db = MilvusDB(
    uri=settings.MILVUS_URI,
    db_name=settings.MILVUS_DB_NAME,
    collection_name=settings.MILVUS_COLLECTION_NAME,
)
