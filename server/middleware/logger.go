package middleware

import (
	"server/utils"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			combinedInfo := v.Host + " " + v.URI
			if v.Error == nil {
				utils.Logger.Infof("method=%s %s status=%d latency=%v",
					v.Method, combinedInfo, v.Status, v.Latency)
				return nil
			}
			utils.Logger.Errorf("method=%s %s status=%d latency=%v bytes_in=%s bytes_out=%d user_agent=%s request_id=%s error=%s",
				v.Method, combinedInfo, v.Status, v.Latency, v.ContentLength, v.ResponseSize, v.UserAgent, v.RequestID, v.Error.Error())
			return nil
		},
	})
}
