from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(
        env_file=".env", env_file_encoding="utf-8"
    )

    DEEPSEEK_NAME: str = "deepseek-chat"
    DEEPSEEK_API_KEY: str = ""
    DEEPSEEK_API_BASE_URL: str = "https://api.deepseek.com"

    QWEN_NAME: str = "qwen-v1"
    QWEN_API_KEY: str = ""
    QWEN_API_BASE_URL: str = ""

    MILVUS_URI: str = ""
    MILVUS_DIM: int = 1024

    PYTHONPATH: str | None = None


settings = Settings()
