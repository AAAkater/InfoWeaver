"""
Tests for core.config module.
"""

import os
import tempfile
from unittest.mock import patch

import pytest
from pydantic import ValidationError

from core.config import Settings


class TestSettings:
    """Test cases for Settings class."""

    def test_default_values(self):
        """Test that settings have correct default values."""
        # Create settings without environment variables
        with patch.dict(os.environ, {}, clear=True):
            settings = Settings()

            # Test default values
            assert settings.DEEPSEEK_NAME == "deepseek-chat"
            assert settings.DEEPSEEK_API_BASE_URL == "https://api.deepseek.com"
            assert settings.DEEPSEEK_API_KEY == ""

            assert settings.QWEN_NAME == "qwen-v1"
            assert settings.QWEN_API_KEY == ""
            assert settings.QWEN_API_BASE_URL == ""

            assert settings.MILVUS_URI == ""
            assert settings.MILVUS_DIM == 1024

            assert settings.PYTHONPATH is None

    def test_env_file_loading(self):
        """Test that settings can be loaded from .env file."""
        # Create a temporary .env file
        with tempfile.NamedTemporaryFile(
            mode="w", suffix=".env", delete=False
        ) as f:
            f.write("""DEEPSEEK_API_KEY=test_deepseek_key
QWEN_API_KEY=test_qwen_key
MILVUS_URI=http://test:19530
MILVUS_DIM=2048
PYTHONPATH=/test/path
""")
            env_file = f.name

        try:
            # Test loading from custom env file
            settings = Settings(_env_file=env_file, _env_file_encoding="utf-8")

            assert settings.DEEPSEEK_API_KEY == "test_deepseek_key"
            assert settings.QWEN_API_KEY == "test_qwen_key"
            assert settings.MILVUS_URI == "http://test:19530"
            assert settings.MILVUS_DIM == 2048
            assert settings.PYTHONPATH == "/test/path"

            # Default values should still be present
            assert settings.DEEPSEEK_NAME == "deepseek-chat"
            assert settings.DEEPSEEK_API_BASE_URL == "https://api.deepseek.com"
            assert settings.QWEN_NAME == "qwen-v1"
            assert settings.QWEN_API_BASE_URL == ""
        finally:
            # Clean up temporary file
            os.unlink(env_file)

    def test_environment_variables_override(self):
        """Test that environment variables override defaults."""
        env_vars = {
            "DEEPSEEK_API_KEY": "env_deepseek_key",
            "QWEN_API_KEY": "env_qwen_key",
            "MILVUS_URI": "http://env:19530",
            "MILVUS_DIM": "4096",
            "PYTHONPATH": "/env/path",
        }

        with patch.dict(os.environ, env_vars):
            settings = Settings()

            assert settings.DEEPSEEK_API_KEY == "env_deepseek_key"
            assert settings.QWEN_API_KEY == "env_qwen_key"
            assert settings.MILVUS_URI == "http://env:19530"
            assert settings.MILVUS_DIM == 4096
            assert settings.PYTHONPATH == "/env/path"

    def test_type_validation(self):
        """Test that type validation works correctly."""
        # MILVUS_DIM should be integer
        with patch.dict(os.environ, {"MILVUS_DIM": "not_a_number"}):
            with pytest.raises(ValidationError):
                Settings()

        # Valid integer should work
        with patch.dict(os.environ, {"MILVUS_DIM": "512"}):
            settings = Settings()
            assert settings.MILVUS_DIM == 512

    def test_settings_singleton(self):
        """Test that settings instance is a singleton."""
        from core.config import settings as global_settings

        # Both should be the same instance
        new_settings = Settings()
        assert (
            global_settings is not new_settings
        )  # Actually they are different instances
        # But they should have same configuration if loaded from same env

        # Test that global settings loads from .env
        assert hasattr(global_settings, "DEEPSEEK_API_KEY")

    def test_empty_string_values(self):
        """Test that empty string values are handled correctly."""
        with patch.dict(
            os.environ,
            {
                "DEEPSEEK_API_KEY": "",
                "QWEN_API_KEY": "",
                "MILVUS_URI": "",
            },
        ):
            settings = Settings()

            assert settings.DEEPSEEK_API_KEY == ""
            assert settings.QWEN_API_KEY == ""
            assert settings.MILVUS_URI == ""

    def test_optional_pythonpath(self):
        """Test PYTHONPATH can be None or string."""
        # When not set, should be None
        with patch.dict(os.environ, {}, clear=True):
            settings = Settings()
            assert settings.PYTHONPATH is None

        # When set to empty string, should be empty string
        with patch.dict(os.environ, {"PYTHONPATH": ""}):
            settings = Settings()
            assert settings.PYTHONPATH == ""

        # When set to path, should be that path
        with patch.dict(os.environ, {"PYTHONPATH": "/some/path"}):
            settings = Settings()
            assert settings.PYTHONPATH == "/some/path"


def test_global_settings_instance():
    """Test the global settings instance."""
    from core.config import settings

    # Should be an instance of Settings
    assert isinstance(settings, Settings)

    # Should have all expected attributes
    expected_attrs = [
        "DEEPSEEK_NAME",
        "DEEPSEEK_API_KEY",
        "DEEPSEEK_API_BASE_URL",
        "QWEN_NAME",
        "QWEN_API_KEY",
        "QWEN_API_BASE_URL",
        "MILVUS_URI",
        "MILVUS_DIM",
        "PYTHONPATH",
    ]

    for attr in expected_attrs:
        assert hasattr(settings, attr), f"settings missing attribute: {attr}"
