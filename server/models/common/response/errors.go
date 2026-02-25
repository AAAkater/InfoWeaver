package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func ErrInvalidToken() error {
	return &echo.HTTPError{
		Code:    http.StatusUnauthorized,
		Message: "invalid or expired jwt",
	}
}

func ErrEmailAlreadyUsed() error {
	return &echo.HTTPError{
		Code:    http.StatusForbidden,
		Message: "this email has been already used",
	}
}

func ErrUnknownError() error {
	return &echo.HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "Unknown error",
	}
}

func ErrUserNotFound() error {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: "User does not exist",
	}
}

func ErrInvalidPassword() error {
	return &echo.HTTPError{
		Code:    http.StatusForbidden,
		Message: "Invalid password",
	}
}
