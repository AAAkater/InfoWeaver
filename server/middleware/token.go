package middleware

import (
	"server/config"
	"server/utils"

	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func TokenMiddleware() echo.MiddlewareFunc {
	config := echoJwt.Config{
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(utils.JwtCustomClaims)
		},
		SigningKey: config.Settings.GetJWTSigningKey(),
	}
	return echoJwt.WithConfig(config)
}
