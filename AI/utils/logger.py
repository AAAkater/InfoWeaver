import sys

from loguru import logger

logger.remove()
logger.add(
    "./logs/app_{time:YYYY-MM-DD}.log",
    rotation="500 MB",
    retention="10 days",
    format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {message}",
    level="DEBUG",
    enqueue=True,
)
logger.add(sys.stderr, level="INFO")
