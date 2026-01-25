package api

import (
	v1 "server/api/v1"

	"github.com/labstack/echo/v5"
)

func InitRouter(e *echo.Echo) {
	v1.SetUserRouter(e)
	v1.SetSwaggerRouter(e)
}
