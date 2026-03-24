import pytest
from pytest_mock import MockerFixture

from db.milvus_db import MilvusDB, VectorEntity


class TestMilvusDBInit:
    """Tests for MilvusDB initialization."""

    def test_init_default_database(self, mocker: MockerFixture) -> None:
        """Test initialization with default database."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        mock_milvus_client.assert_called_once_with(uri="http://localhost:19530", timeout=1000)
        assert db.db_name == "default"
        assert db.collection_name == "test_collection"

    def test_init_custom_database_creates_if_not_exists(self, mocker: MockerFixture) -> None:
        """Test initialization creates custom database if it doesn't exist."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection", db_name="custom_db")

        mock_client_instance.create_database.assert_called_once_with(db_name="custom_db")
        mock_client_instance.use_database.assert_called_once_with(db_name="custom_db")
        assert db.db_name == "custom_db"

    def test_init_custom_database_uses_existing(self, mocker: MockerFixture) -> None:
        """Test initialization uses existing custom database."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default", "existing_db"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        _db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection", db_name="existing_db")

        mock_client_instance.create_database.assert_not_called()
        mock_client_instance.use_database.assert_called_once_with(db_name="existing_db")


class TestMilvusDBCollection:
    """Tests for MilvusDB collection operations."""

    def test_create_collection(self, mocker: MockerFixture) -> None:
        """Test collection creation."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        mock_schema = mocker.MagicMock()
        mock_milvus_client.create_schema.return_value = mock_schema

        MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        mock_milvus_client.create_schema.assert_called_once_with(auto_id=True, enable_dynamic_field=True)
        mock_schema.add_field.assert_called()
        mock_client_instance.create_collection.assert_called_once()


class TestMilvusDBInsert:
    """Tests for MilvusDB insert operations."""

    def test_insert_entities(self, mocker: MockerFixture) -> None:
        """Test inserting entities."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        entities = [
            VectorEntity(vector=[0.1] * 1024, content="content 1", dataset_id=1),
            VectorEntity(vector=[0.2] * 1024, content="content 2", dataset_id=1),
        ]

        db.insert_entities(entities)

        mock_client_instance.insert.assert_called_once()
        call_args = mock_client_instance.insert.call_args
        assert call_args.kwargs["collection_name"] == "test_collection"
        assert len(call_args.kwargs["data"]) == 2

    def test_insert_empty_list(self, mocker: MockerFixture) -> None:
        """Test inserting empty list."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        db.insert_entities([])

        mock_client_instance.insert.assert_called_once_with(
            collection_name="test_collection",
            data=[],
        )


class TestMilvusDBDelete:
    """Tests for MilvusDB delete operations."""

    def test_delete_entities_by_id(self, mocker: MockerFixture) -> None:
        """Test deleting entities by ID."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        db.delete_entities_by_id([1, 2, 3])

        mock_client_instance.delete.assert_called_once_with(
            collection_name="test_collection",
            ids=[1, 2, 3],
        )

    def test_delete_entities_by_filter(self, mocker: MockerFixture) -> None:
        """Test deleting entities by filter expression."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        db.delete_entities_by_filter("dataset_id == 1")

        mock_client_instance.delete.assert_called_once_with(
            collection_name="test_collection",
            filter_expression="dataset_id == 1",
        )


class TestMilvusDBSearch:
    """Tests for MilvusDB search operations."""

    def test_search_entities_filter_basic(self, mocker: MockerFixture) -> None:
        """Test basic search with dataset_id filter."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        mock_search_result = [
            [
                {"id": 1, "distance": 0.1, "entity": {"content": "result 1"}},
                {"id": 2, "distance": 0.2, "entity": {"content": "result 2"}},
            ]
        ]
        mock_client_instance.search.return_value = mock_search_result

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        results = db.search_entities_filter(
            query_vector=[0.1] * 1024,
            dataset_id=1,
            top_k=10,
        )

        assert len(results) == 2
        mock_client_instance.search.assert_called_once()
        call_args = mock_client_instance.search.call_args
        assert call_args.kwargs["collection_name"] == "test_collection"
        assert call_args.kwargs["limit"] == 10
        assert call_args.kwargs["filter"] == "dataset_id == 1"
        assert call_args.kwargs["output_fields"] == ["id", "content"]

    def test_search_entities_filter_with_expr(self, mocker: MockerFixture) -> None:
        """Test search with additional filter expression."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        mock_client_instance.search.return_value = [[]]

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        db.search_entities_filter(
            query_vector=[0.1] * 1024,
            dataset_id=1,
            top_k=5,
            expr="content like '%test%'",
        )

        call_args = mock_client_instance.search.call_args
        assert call_args.kwargs["filter"] == "dataset_id == 1 and content like '%test%'"
        assert call_args.kwargs["limit"] == 5

    def test_search_entities_empty_result(self, mocker: MockerFixture) -> None:
        """Test search returning empty results."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        mock_client_instance.search.return_value = []

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        results = db.search_entities_filter(
            query_vector=[0.1] * 1024,
            dataset_id=1,
        )

        assert results == []


class TestMilvusDBGetCollectionList:
    """Tests for get_collection_list property."""

    @pytest.mark.asyncio
    async def test_get_collection_list(self, mocker: MockerFixture) -> None:
        """Test getting collection list."""
        mock_milvus_client = mocker.patch("db.milvus_db.MilvusClient")
        mock_settings = mocker.patch("db.milvus_db.settings")
        mock_settings.MILVUS_DIM = 1024

        mock_client_instance = mocker.MagicMock()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}
        mock_client_instance.list_collections = mocker.AsyncMock(return_value=["collection1", "collection2"])

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection")

        result = await db.get_collection_list

        assert result == ["collection1", "collection2"]
        mock_client_instance.list_collections.assert_awaited_once()
        mock_milvus_client.return_value = mock_client_instance
        mock_client_instance.list_databases.return_value = ["default", "existing_db"]
        mock_client_instance.get_collection_stats.return_value = {"row_count": 0}

        db = MilvusDB(uri="http://localhost:19530", collection_name="test_collection", db_name="existing_db")

        mock_client_instance.create_database.assert_not_called()
        mock_client_instance.use_database.assert_called_once_with(db_name="existing_db")
