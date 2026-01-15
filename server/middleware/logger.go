package middleware

import (
	"server/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)
			utils.Logger.Infof("HTTP Request: method=%s uri=%s ip=%s status=%d duration=%v error=%v",
				c.Request().Method,
				c.Request().RequestURI,
				c.RealIP(),
				c.Response().Status,
				duration,
				err,
			)

			return err
		}
	}
}
