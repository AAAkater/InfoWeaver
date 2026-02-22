from core.config import Settings


class TestSettings:
    """Test cases for Settings class."""

    def test_settings_default_values(self):
        """Test that settings have correct default values."""
        settings = Settings(_env_file=".env.example", _env_file_encoding="utf-8")  # type: ignore

        # Test PostgreSQL defaults
        assert settings.POSTGRES_USERNAME == "postgres"
        assert settings.POSTGRES_PASSWORD == "postgres"
        assert settings.POSTGRES_PORT == 5432
        assert settings.POSTGRES_HOST == "localhost"
        assert settings.POSTGRES_DB == "InfoWeaver"

        # Test Redis defaults
        assert settings.REDIS_HOST == "localhost"
        assert settings.REDIS_PORT == 6379
        assert settings.REDIS_PASSWORD == ""
        assert settings.REDIS_DB == 0
        assert settings.REDIS_EXPIRE == 600

        # Test Milvus defaults
        assert settings.MILVUS_HOST == "localhost"
        assert settings.MILVUS_PORT == 19530
        assert settings.MILVUS_DIM == 1024

        # Test MinIO defaults
        assert settings.MINIO_HOST == "localhost"
        assert settings.MINIO_PORT == 9000
        assert settings.MINIO_ACCESS_KEY == "minioadmin"
        assert settings.MINIO_SECRET_KEY == "minioadmin"
        assert settings.MINIO_BUCKET_NAME == "info-weaver"
        assert settings.MINIO_USE_SSL is False

        # Test RabbitMQ defaults
        assert settings.RABBITMQ_HOST == "localhost"
        assert settings.RABBITMQ_PORT == 5672
        assert settings.RABBITMQ_USER == "guest"
        assert settings.RABBITMQ_PASSWORD == "guest"
        assert settings.RABBITMQ_VHOST == "/"
        assert settings.RABBITMQ_QUEUE == "info-weaver-file-queue"

        # Test API version
        assert settings.API_VER_STR == "/api/v1"

    def test_postgresql_dsn_computed_field(self):
        """Test PostgreSQL DSN is correctly computed."""
        settings = Settings(
            POSTGRES_USERNAME="test_user",
            POSTGRES_PASSWORD="test_pass",
            POSTGRES_HOST="localhost",
            POSTGRES_PORT=5433,
            POSTGRES_DB="test_db",
        )

        expected_dsn = "postgresql+psycopg2://test_user:test_pass@localhost:5433/test_db"
        assert str(settings.POSTGRESQL_DSN) == expected_dsn

    def test_redis_dsn_computed_field(self):
        """Test Redis DSN is correctly computed."""
        settings = Settings(
            REDIS_HOST="redis.example.com",
            REDIS_PORT=6380,
            REDIS_PASSWORD="redis_pass",
            REDIS_DB=1,
        )

        expected_dsn = "redis://:redis_pass@redis.example.com:6380//1"
        assert str(settings.REDIS_DSN) == expected_dsn

    def test_redis_dsn_without_password(self):
        """Test Redis DSN without password."""
        settings = Settings(
            REDIS_HOST="localhost",
            REDIS_PORT=6379,
            REDIS_PASSWORD="",
            REDIS_DB=0,
        )

        expected_dsn = "redis://localhost:6379//0"
        assert str(settings.REDIS_DSN) == expected_dsn

    def test_milvus_uri_computed_field(self):
        """Test Milvus URI is correctly computed."""
        settings = Settings(
            MILVUS_HOST="milvus.example.com",
            MILVUS_PORT=19531,
        )

        expected_uri = "http://milvus.example.com:19531"
        assert settings.MILVUS_URI == expected_uri

    def test_minio_url_computed_field_http(self):
        """Test MinIO URL with HTTP protocol."""
        settings = Settings(
            MINIO_HOST="minio.example.com",
            MINIO_PORT=9001,
            MINIO_USE_SSL=False,
        )

        expected_url = "http://minio.example.com:9001"
        assert settings.MINIO_URL == expected_url

    def test_minio_url_computed_field_https(self):
        """Test MinIO URL with HTTPS protocol."""
        settings = Settings(
            MINIO_HOST="minio.example.com",
            MINIO_PORT=9001,
            MINIO_USE_SSL=True,
        )

        expected_url = "https://minio.example.com:9001"
        assert settings.MINIO_URL == expected_url

    def test_rabbitmq_url_computed_field(self):
        """Test RabbitMQ URL is correctly computed."""
        settings = Settings(
            RABBITMQ_USER="rmq_user",
            RABBITMQ_PASSWORD="rmq_pass",
            RABBITMQ_HOST="rabbitmq.example.com",
            RABBITMQ_PORT=5673,
            RABBITMQ_VHOST="/test_vhost",
        )

        expected_url = "amqp://rmq_user:rmq_pass@rabbitmq.example.com:5673/test_vhost"
        assert settings.RABBITMQ_URL == expected_url

    def test_settings_extra_fields_ignored(self, monkeypatch):
        """Test that extra fields are ignored (extra='ignore')."""
        monkeypatch.setenv("UNKNOWN_FIELD", "should_be_ignored")

        # Should not raise an error
        settings = Settings()
        assert not hasattr(settings, "UNKNOWN_FIELD")

    def test_settings_postgresql_dsn_with_special_chars(self):
        """Test PostgreSQL DSN with special characters in password."""
        settings = Settings(
            POSTGRES_USERNAME="user",
            POSTGRES_PASSWORD="password123",
            POSTGRES_HOST="localhost",
            POSTGRES_PORT=5432,
            POSTGRES_DB="testdb",
        )

        # Pydantic should handle URL encoding
        dsn_str = str(settings.POSTGRESQL_DSN)
        assert "postgresql+psycopg2://" in dsn_str
        assert "localhost" in dsn_str
        assert "testdb" in dsn_str

    def test_settings_redis_dsn_with_special_chars(self):
        """Test Redis DSN with special characters in password."""
        settings = Settings(
            REDIS_HOST="localhost",
            REDIS_PORT=6379,
            REDIS_PASSWORD="p@ss:word",
            REDIS_DB=2,
        )

        dsn_str = str(settings.REDIS_DSN)
        assert "redis://" in dsn_str
        assert "localhost" in dsn_str
        assert "/2" in dsn_str

    def test_settings_minio_default_bucket(self):
        """Test MinIO default bucket name."""
        settings = Settings()
        assert settings.MINIO_BUCKET_NAME == "info-weaver"

    def test_settings_rabbitmq_default_queue(self):
        """Test RabbitMQ default queue name."""
        settings = Settings()
        assert settings.RABBITMQ_QUEUE == "info-weaver-file-queue"

    def test_settings_api_version(self):
        """Test API version string."""
        settings = Settings()
        assert settings.API_VER_STR == "/api/v1"

        # Test custom API version
        settings_custom = Settings(API_VER_STR="/api/v2")
        assert settings_custom.API_VER_STR == "/api/v2"

    def test_settings_postgresql_dsn_port_included(self):
        """Test that PostgreSQL DSN includes port number."""
        settings = Settings(
            POSTGRES_USERNAME="user",
            POSTGRES_PASSWORD="pass",
            POSTGRES_HOST="localhost",
            POSTGRES_PORT=5432,
            POSTGRES_DB="db",
        )

        dsn_str = str(settings.POSTGRESQL_DSN)
        assert ":5432" in dsn_str

    def test_settings_redis_expire_default(self):
        """Test Redis default expiration time."""
        settings = Settings()
        assert settings.REDIS_EXPIRE == 600

    def test_settings_redis_expire_custom(self):
        """Test Redis custom expiration time."""
        settings = Settings(REDIS_EXPIRE=3600)
        assert settings.REDIS_EXPIRE == 3600

    def test_settings_milvus_dimension_default(self):
        """Test Milvus default dimension."""
        settings = Settings()
        assert settings.MILVUS_DIM == 1024

    def test_settings_milvus_dimension_custom(self):
        """Test Milvus custom dimension."""
        settings = Settings(MILVUS_DIM=768)
        assert settings.MILVUS_DIM == 768

    def test_settings_minio_ssl_disabled(self):
        """Test MinIO SSL disabled by default."""
        settings = Settings()
        assert settings.MINIO_USE_SSL is False

    def test_settings_minio_ssl_enabled(self):
        """Test MinIO SSL enabled."""
        settings = Settings(MINIO_USE_SSL=True)
        assert settings.MINIO_USE_SSL is True
        assert settings.MINIO_URL.startswith("https://")

    def test_settings_rabbitmq_vhost_default(self):
        """Test RabbitMQ default vhost."""
        settings = Settings()
        assert settings.RABBITMQ_VHOST == "/"

    def test_settings_rabbitmq_vhost_custom(self):
        """Test RabbitMQ custom vhost."""
        settings = Settings(RABBITMQ_VHOST="/myapp")
        assert settings.RABBITMQ_VHOST == "/myapp"
        assert "/myapp" in settings.RABBITMQ_URL
