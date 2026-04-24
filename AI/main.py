"""FastAPI application for document processing service."""

from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from api import v1_router
from configs.app_config import settings
from utils import logger


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager."""
    logger.info("Starting document processing API service...")
    logger.info(f"Service running with log level: {settings.LOG_LEVEL}")
    yield
    logger.info("Shutting down document processing API service...")


app = FastAPI(
    title="InfoWeaver Document Processing API",
    description="API for processing documents: upload, split, embed, and store in vector database",
    version="0.1.0",
    lifespan=lifespan,
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(v1_router)


@app.get("/health", tags=["health"])
async def health_check():
    """Health check endpoint."""
    return {"status": "healthy", "service": "document-processing"}
