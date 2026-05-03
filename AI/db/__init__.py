"""Database initialization and connection health checks.

Provides ``lifespan`` context manager for FastAPI startup/shutdown,
testing all database / external service connections on startup.
"""

from contextlib import asynccontextmanager

from fastapi import FastAPI

from configs.app_config import settings
from utils import logger


def check_postgresql() -> bool:
    """Test PostgreSQL connection by executing a simple query."""
    from sqlalchemy import text

    from db.postgresql_db import engine

    try:
        with engine.connect() as conn:
            conn.execute(text("SELECT 1"))
            conn.commit()
        logger.info("✓ PostgreSQL connected successfully")
        return True
    except Exception as e:
        logger.error(f"✗ PostgreSQL connection failed: {e}")
        return False


def check_redis() -> bool:
    """Test Redis connection with ping."""
    from db.redis_client import redis_client

    try:
        if redis_client.ping():
            logger.info("✓ Redis connected successfully")
            return True
        logger.error("✗ Redis ping failed")
        return False
    except Exception as e:
        logger.error(f"✗ Redis connection failed: {e}")
        return False


async def check_milvus() -> bool:
    """Test Milvus connection by listing collections."""
    from pymilvus import MilvusClient

    try:
        client = MilvusClient(uri=settings.MILVUS_URI, db_name=settings.MILVUS_DB_NAME, timeout=10)
        client.list_collections()  # type: ignore
        logger.info("✓ Milvus connected successfully")
        return True
    except Exception as e:
        logger.error(f"✗ Milvus connection failed: {e}")
        return False


def check_minio() -> bool:
    """Test MinIO connection by listing buckets."""
    from minio import Minio

    try:
        client = Minio(
            endpoint=settings.MINIO_ENDPOINT,
            access_key=settings.MINIO_ACCESS_KEY,
            secret_key=settings.MINIO_SECRET_KEY,
            secure=settings.MINIO_USE_SSL,
        )
        buckets = client.list_buckets()
        logger.info(f"✓ MinIO connected successfully (buckets: {[b.name for b in buckets]})")
        return True
    except Exception as e:
        logger.error(f"✗ MinIO connection failed: {e}")
        return False


def check_rabbitmq() -> bool:
    """Test RabbitMQ connection by opening and closing a connection."""
    from pika import BlockingConnection, URLParameters

    try:
        url = f"amqp://{settings.RABBITMQ_USER}:{settings.RABBITMQ_PASSWORD}@{settings.RABBITMQ_HOST}:{settings.RABBITMQ_PORT}{settings.RABBITMQ_VHOST}"
        conn = BlockingConnection(URLParameters(url))
        conn.close()
        logger.info("✓ RabbitMQ connected successfully")
        return True
    except Exception as e:
        logger.error(f"✗ RabbitMQ connection failed: {e}")
        return False


async def check_all_db_connections() -> None:
    """Run all database connection health checks and log a summary."""
    logger.info("Checking database connections...")

    results = {
        "PostgreSQL": check_postgresql(),
        "Redis": check_redis(),
        "Milvus": await check_milvus(),
        "MinIO": check_minio(),
        "RabbitMQ": check_rabbitmq(),
    }

    success = sum(1 for v in results.values() if v)
    total = len(results)
    logger.info(f"DB connection check complete: {success}/{total} passed")

    if success < total:
        failed = [name for name, ok in results.items() if not ok]
        logger.warning(f"Some connections failed: {', '.join(failed)}")


def close_db_connections() -> None:
    """Close / dispose all database connections on shutdown."""
    logger.info("Closing database connections...")

    # PostgreSQL engine disposal
    try:
        from db.postgresql_db import engine

        engine.dispose()
        logger.info("✓ PostgreSQL engine disposed")
    except Exception as e:
        logger.warning(f"Failed to dispose PostgreSQL engine: {e}")

    # Redis
    try:
        from db.redis_client import redis_client

        redis_client.close()
        logger.info("✓ Redis connection closed")
    except Exception as e:
        logger.warning(f"Failed to close Redis connection: {e}")


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan: check DB connections on startup, clean up on shutdown."""
    logger.info("Starting document processing API service...")
    logger.info(f"Service running with log level: {settings.SERVER_LOG_LEVEL}")

    await check_all_db_connections()

    yield

    close_db_connections()
    logger.info("Shutting down document processing API service...")
