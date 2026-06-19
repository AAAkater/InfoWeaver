"""Unified API response model."""

from pydantic import BaseModel, Field


class APIResponse[T](BaseModel):
    """Unified API response wrapper."""

    code: int = Field(default=0, description="Status code, 0 = success")
    msg: str = Field(default="success", description="Response message")
    data: T | None = Field(default=None, description="Response data payload")
