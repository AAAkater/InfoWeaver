package utils

import (
	"errors"
	"server/config"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	TokenExpired          = errors.New("token has expired")
	TokenNotValidYet      = errors.New("token is not yet valid")
	TokenMalformed        = errors.New("this is not a token")
	TokenSignatureInvalid = errors.New("invalid signature")
	TokenInvalid          = errors.New("unable to process this token")
)

var (
	jwtToolInstance *JwtInfo
	once            sync.Once
)

func JwtTool() *JwtInfo {
	init := func() {
		jwtToolInstance = &JwtInfo{
			config.Settings.GetJWTBufferTime(),
			config.Settings.GetJWTExpireTime(),
			[]byte(config.Settings.JWT_SIGNING_KEY),
		}
	}
	once.Do(init)
	return jwtToolInstance
}

type JwtCustomClaims struct {
	UserID     uint
	Username   string
	IsAdmin    bool
	BufferTime int64
	jwt.RegisteredClaims
}

type JwtInfo struct {
	bufferTime  time.Duration
	expiresTime time.Duration
	signingKey  []byte
}

func (this *JwtInfo) CreateToken(UserID uint, Username string, IsAdmin bool) (string, error) {

	claims := JwtCustomClaims{
		UserID:   UserID,
		Username: Username,
		IsAdmin:  IsAdmin,
		// Buffer time in seconds - during buffer period user gets new refresh token
		// User may have two valid tokens but frontend keeps only one
		BufferTime: int64(this.bufferTime / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{"GVA"},
			// NotBefore claim - token becomes valid from this time
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			// ExpiresAt claim - token expiration time from config
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(this.expiresTime)),
			Issuer:    config.Settings.JWT_ISSUER,
		},
	}

	// Create new token with ES256 signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	// Sign the token with signing key and return
	return token.SignedString(this.signingKey)
}

func (this *JwtInfo) ParseToken(token_string string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_string,
		&JwtCustomClaims{},
		func(token *jwt.Token) (i any, e error) {
			return this.signingKey, nil
		})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
