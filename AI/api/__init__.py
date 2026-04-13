"""API routes package."""

from fastapi import APIRouter

from api.document_routes import router as document_router
from api.search_routes import router as search_router

v1_router = APIRouter()
# Include routers
v1_router.include_router(document_router, prefix="/api/v1/documents", tags=["documents"])
v1_router.include_router(search_router, prefix="/api/v1/search", tags=["search"])
