package middleware

import "github.com/labstack/echo/v4"

func InitMiddleWares(e *echo.Echo) {
	e.Use(LoggerMiddleware())
}
