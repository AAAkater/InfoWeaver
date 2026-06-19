"""FastAPI application for document processing service."""

from fastapi import FastAPI

import middlewares
from api import v1_router
from db import lifespan

app = FastAPI(
    title="InfoWeaver Document Processing API",
    description="API for processing documents: upload, split, embed, and store in vector database",
    version="0.1.0",
    lifespan=lifespan,
)

# Add CORS middleware
middlewares.register_middlewares(app)

app.include_router(v1_router)

if __name__ == "__main__":
    import uvicorn

    from configs.app_config import settings

    uvicorn.run(
        "main:app",
        host=settings.SERVER_HOST,
        port=settings.SERVER_PORT,
        log_level=settings.SERVER_LOG_LEVEL.lower(),
    )
