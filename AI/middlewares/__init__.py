def register_middlewares(app):
    """Add middlewares to the FastAPI application."""
    from .cors import setup_cors

    setup_cors(app)
