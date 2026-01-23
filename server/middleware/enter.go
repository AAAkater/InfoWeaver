package middleware

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func InitMiddleWares(e *echo.Echo) {
	e.HTTPErrorHandler = CustomHTTPErrorHandler(false)
	e.Use(middleware.RequestLogger())
}
