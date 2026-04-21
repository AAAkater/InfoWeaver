from datetime import datetime
from typing import Any

from sqlalchemy import DateTime, ForeignKey, Integer, String, Text, func
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


class User(Base):
    """User table - represents system users with role-based permissions"""

    __tablename__ = "users"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    username: Mapped[str] = mapped_column(String(255), nullable=False)
    email: Mapped[str] = mapped_column(String(255), unique=True, nullable=False, index=True)
    password: Mapped[str] = mapped_column(String(255), nullable=False)
    role: Mapped[str] = mapped_column(String(50), nullable=False, default="user")  # "user" or "admin"


class Provider(Base):
    """Provider table - represents AI model providers (OpenAI, Gemini, Anthropic, Ollama, etc.)"""

    __tablename__ = "providers"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    name: Mapped[str] = mapped_column(String(255), nullable=False, index=True)
    mode: Mapped[str] = mapped_column(
        String(255), nullable=False
    )  # "openai", "openai_response", "gemini", "anthropic", "ollama"
    base_url: Mapped[str] = mapped_column(String(500), nullable=False)
    api_key: Mapped[str] = mapped_column(String(500), nullable=False)
    available_models: Mapped[dict[str, bool] | None] = mapped_column(JSONB, nullable=True)  # model_id -> enabled
    owner_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    user: Mapped["User"] = relationship("User", backref="providers", lazy="select")


class Dataset(Base):
    """Dataset table - represents a collection of files owned by a user"""

    __tablename__ = "datasets"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    name: Mapped[str] = mapped_column(String(255), nullable=False)
    icon: Mapped[str | None] = mapped_column(String(50), nullable=True)  # Emoji icon (e.g., 🚀, ❤️)
    description: Mapped[str | None] = mapped_column(Text, nullable=True)
    search_type: Mapped[str] = mapped_column(String(50), nullable=False, default="dense")  # "sparse", "dense", "hybrid"
    embedding_model: Mapped[str] = mapped_column(String(255), nullable=False)
    provider_id: Mapped[int] = mapped_column(Integer, ForeignKey("providers.id", ondelete="CASCADE"), nullable=False)
    owner_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    user: Mapped["User"] = relationship("User", backref="datasets", lazy="select")
    provider: Mapped["Provider"] = relationship("Provider", backref="datasets", lazy="select")


class File(Base):
    """File table - represents uploaded files stored in MinIO"""

    __tablename__ = "files"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    name: Mapped[str] = mapped_column(String(500), nullable=False)  # Original filename
    minio_path: Mapped[str] = mapped_column(String(1000), unique=True, nullable=False)  # MinIO object path (bucket/key)
    size: Mapped[int] = mapped_column(Integer, nullable=False)  # File size in bytes
    type: Mapped[str] = mapped_column(String(255), nullable=False)  # MIME type
    dataset_id: Mapped[int] = mapped_column(Integer, ForeignKey("datasets.id", ondelete="CASCADE"), nullable=False)
    user_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    user: Mapped["User"] = relationship("User", backref="files", lazy="select")
    dataset: Mapped["Dataset"] = relationship("Dataset", backref="files", lazy="select")


class Chunk(Base):
    """Chunk table - represents document chunks for RAG"""

    __tablename__ = "chunks"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        server_default=func.now(),
        onupdate=func.now(),
        nullable=False,
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    content: Mapped[str] = mapped_column(Text, nullable=False)
    chunk_metadata: Mapped[dict[str, Any] | None] = mapped_column(
        JSONB, nullable=True
    )  # Additional metadata (source, type, etc.)
    vector_id: Mapped[str] = mapped_column(String(255), unique=True, nullable=False, index=True)
    file_id: Mapped[int] = mapped_column(Integer, ForeignKey("files.id", ondelete="CASCADE"), nullable=False)
    file: Mapped["File"] = relationship("File", backref="chunks", lazy="select")


class ChatSession(Base):
    """ChatSession table - represents a conversation session with an AI model"""

    __tablename__ = "chat_sessions"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    title: Mapped[str] = mapped_column(String(500), nullable=False)
    owner_id: Mapped[int] = mapped_column(Integer, ForeignKey("users.id", ondelete="CASCADE"), nullable=False)
    user: Mapped["User"] = relationship("User", backref="chat_sessions", lazy="select")


# Association table for many-to-many relationship between Memory and Chunk
memory_documents = Base.metadata.tables.get("memory_documents")
if memory_documents is None:
    from sqlalchemy import Column, Table

    memory_documents = Table(
        "memory_documents",
        Base.metadata,
        Column("memory_id", Integer, ForeignKey("memories.id", ondelete="CASCADE"), primary_key=True),
        Column("chunk_id", Integer, ForeignKey("chunks.id", ondelete="CASCADE"), primary_key=True),
    )


class Memory(Base):
    """Memory table - stores a single chat message within a session"""

    __tablename__ = "memories"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, autoincrement=True)
    created_at: Mapped[datetime] = mapped_column(DateTime(timezone=True), server_default=func.now(), nullable=False)
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True), server_default=func.now(), onupdate=func.now(), nullable=False
    )
    deleted_at: Mapped[datetime | None] = mapped_column(DateTime(timezone=True), nullable=True)

    session_id: Mapped[int] = mapped_column(
        Integer, ForeignKey("chat_sessions.id", ondelete="CASCADE"), nullable=False, index=True
    )
    content: Mapped[str] = mapped_column(Text, nullable=False)
    role: Mapped[str] = mapped_column(String(50), nullable=False, default="user")  # "user", "assistant", "system"
    chat_session: Mapped["ChatSession"] = relationship("ChatSession", backref="memories", lazy="select")
    retrieved_documents: Mapped[list["Chunk"]] = relationship(
        "Chunk", secondary=memory_documents, backref="memories", lazy="select"
    )
