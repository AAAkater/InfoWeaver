from core.config import settings
from utils.logger import logger

uri = settings.MILVUS_URI


logger.info(f"Milvus URI: {uri}")
