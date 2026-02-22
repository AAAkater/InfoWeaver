from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from core.config import settings

# Create SQLAlchemy engine
engine = create_engine(str(settings.POSTGRESQL_DSN), echo=True)


def get_db_session():
    """
    Dependency function to get database session.
    Use this in FastAPI route dependencies.
    """
    with Session(engine) as session:
        yield session
