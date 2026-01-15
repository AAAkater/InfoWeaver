package middleware

import (
	"server/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			err := next(ctx)

			duration := time.Since(start)
			utils.Logger.Infof("HTTP Request: method=%s uri=%s ip=%s status=%d duration=%v error=%v",
				ctx.Request().Method,
				ctx.Request().RequestURI,
				ctx.RealIP(),
				ctx.Response().Status,
				duration,
				err,
			)

			return err
		}
	}
}
