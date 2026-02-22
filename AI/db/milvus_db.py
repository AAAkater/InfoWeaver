from pydantic import BaseModel
from pymilvus import DataType, MilvusClient

from core.config import settings
from utils.logger import logger


class VectorEntity(BaseModel):
    id: int
    vector: list[float]
    content: str


class MilvusDB:
    def __init__(self):
        self.client = MilvusClient(uri=str(settings.MILVUS_URI), timeout=1000)

    def create_collections(self, new_collection_name: str):
        # Create collections for each vector type
        schema = MilvusClient.create_schema(
            auto_id=False,
            enable_dynamic_field=True,
        )

        # 3.2. Add fields to schema
        schema.add_field(field_name="id", datatype=DataType.INT64, is_primary=True)
        schema.add_field(
            field_name="vector",
            datatype=DataType.FLOAT_VECTOR,
            dim=settings.MILVUS_DIM,
        )
        schema.add_field(field_name="content", datatype=DataType.VARCHAR, max_length=512)
        index_params = self.client.prepare_index_params()
        index_params.add_index(field_name="id", index_type="AUTOINDEX")
        index_params.add_index(field_name="vector", index_type="AUTOINDEX", metric_type="COSINE")

        # 3.3. Create collection
        self.client.create_collection(
            collection_name=new_collection_name,
            schema=schema,
            index_params=index_params,
        )

        result = self.client.get_collection_stats(collection_name=new_collection_name)

        logger.info(f"Collection '{new_collection_name}' created successfully. Stats: {result}")

    @property
    async def get_collection_list(self):
        return await self.client.list_collections()

    def insert_entities(self, collection_name: str, new_vector_datas: list[VectorEntity]):

        # 4.2. Insert data into collection
        self.client.insert(
            collection_name=collection_name,
            data=[data.model_dump() for data in new_vector_datas],
        )

        logger.info(f"Inserted {len(new_vector_datas)} records into collection '{collection_name}'")

    def delete_entities_by_id(self, collection_name: str, entity_ids: list[int]):
        # 4.3. Delete entity by ID
        self.client.delete(
            collection_name=collection_name,
            ids=entity_ids,
        )

        logger.info(f"Deleted entities with IDs {entity_ids} from collection '{collection_name}'")

    def delete_entities_by_filter(self, collection_name: str, expr: str):
        # 4.4. Delete entity by filter expression
        self.client.delete(collection_name=collection_name, filter_expression=expr)

        logger.info(f"Deleted entities with filter expression '{expr}' from collection '{collection_name}'")
