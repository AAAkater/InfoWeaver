package utils

import (
	"errors"
	"server/config"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

var (
	TokenExpired          = errors.New("token has expired")
	TokenNotValidYet      = errors.New("token is not yet valid")
	TokenMalformed        = errors.New("this is not a token")
	TokenSignatureInvalid = errors.New("invalid signature")
	TokenInvalid          = errors.New("unable to process this token")
)

type JwtCustomClaims struct {
	ID      uint `json:"name_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func CreateToken(UserID uint, IsAdmin bool) (string, error) {

	expireTime := config.Settings.GetJWTExpireTime()
	signingKey := config.Settings.GetJWTSigningKey()
	claims := JwtCustomClaims{
		ID:      UserID,
		IsAdmin: IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{"GVA"},
			// NotBefore claim - token becomes valid from this time
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			// ExpiresAt claim - token expiration time from config
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
		},
	}

	// Create new token with HS256 signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with signing key and return
	return token.SignedString(signingKey)
}

func GetCurrentUser(ctx *echo.Context) (*JwtCustomClaims, error) {
	token, err := echo.ContextGet[*jwt.Token](ctx, "user")
	if err != nil {
		Logger.Errorf("Failed to get token from context: %v", err)
		return nil, TokenExpired
	}
	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok {
		return nil, TokenInvalid
	}
	return claims, nil
}
