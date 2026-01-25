package middleware

import (
	"server/utils"
	"time"

	"github.com/labstack/echo/v5"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			start := time.Now()

			err := next(ctx)

			duration := time.Since(start)
			req := ctx.Request()

			utils.Logger.Infof("%s uri=%s ip=%s duration=%v error=%v",
				req.Method,
				req.RequestURI,
				ctx.RealIP(),
				duration,
				err,
			)
			return err
		}
	}
}
