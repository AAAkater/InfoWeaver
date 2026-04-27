"""API routes package."""

from fastapi import APIRouter

from api.chat_routes import router as chat_router
from api.document_routes import router as document_router
from api.search_routes import router as search_router

v1_router = APIRouter(prefix="/api/v1", tags=["v1"])
# Include routers
v1_router.include_router(document_router)
v1_router.include_router(search_router)
v1_router.include_router(chat_router)
