from pydantic import PostgresDsn, RedisDsn, computed_field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file="./.env",
        env_ignore_empty=True,
        extra="ignore",
    )
    PYTHONPATH: str | None = None

    # POSTGRESQL
    POSTGRESQL_USER: str = "postgres"
    POSTGRESQL_PASSWORD: str = ""
    POSTGRESQL_PORT: int = 5432
    POSTGRESQL_SERVER: str = "127.0.0.1"
    POSTGRESQL_DB: str = ""

    @computed_field
    @property
    def POSTGRESQL_URI(self) -> PostgresDsn:
        return PostgresDsn.build(
            scheme="postgresql+psycopg2",
            username=self.POSTGRESQL_USER,
            password=self.POSTGRESQL_PASSWORD,
            host=self.POSTGRESQL_SERVER,
            port=self.POSTGRESQL_PORT,
            path=self.POSTGRESQL_DB,
        )

    # REDIS
    REDIS_HOST: str = "127.0.0.1"
    REDIS_PORT: int = 6379
    REDIS_PASSWORD: str = ""
    REDIS_DB: int = 0
    REDIS_EXPIRE: int = 600

    @computed_field
    @property
    def REDIS_URL(self) -> RedisDsn:
        return RedisDsn.build(
            scheme="redis",
            host=self.REDIS_HOST,
            port=self.REDIS_PORT,
            password=self.REDIS_PASSWORD,
            path=f"/{self.REDIS_DB}",
        )

    # MILVUS
    MILVUS_HOST: str = "127.0.0.1"
    MILVUS_PORT: int = 19530
    MILVUS_DIM: int = 1024

    @computed_field
    @property
    def MILVUS_URL(self) -> str:
        return f"http://{self.MILVUS_HOST}:{self.MILVUS_PORT}"

    # MINIO
    MINIO_HOST: str = "127.0.0.1"
    MINIO_PORT: int = 9000
    MINIO_ACCESS_KEY: str = "minioadmin"
    MINIO_SECRET_KEY: str = "minioadmin"
    MINIO_BUCKET_NAME: str = "info-weaver"
    MINIO_USE_SSL: bool = False

    @computed_field
    @property
    def MINIO_URL(self) -> str:
        protocol = "https" if self.MINIO_USE_SSL else "http"
        return f"{protocol}://{self.MINIO_HOST}:{self.MINIO_PORT}"

    # RABBITMQ
    RABBITMQ_HOST: str = "127.0.0.1"
    RABBITMQ_PORT: int = 5672
    RABBITMQ_USER: str = "guest"
    RABBITMQ_PASSWORD: str = "guest"
    RABBITMQ_VHOST: str = "/"

    @computed_field
    @property
    def RABBITMQ_URL(self) -> str:
        return f"amqp://{self.RABBITMQ_USER}:{self.RABBITMQ_PASSWORD}@{self.RABBITMQ_HOST}:{self.RABBITMQ_PORT}{self.RABBITMQ_VHOST}"

    API_VER_STR: str = "/api/v1"


settings = Settings()
