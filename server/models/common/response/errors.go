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

func ErrDatasetNameAlreadyExists() error {
	return &echo.HTTPError{
		Code:    http.StatusForbidden,
		Message: "dataset with the same name already exists",
	}
}

func ErrDatasetNotFound() error {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: "Dataset not found",
	}
}

func ErrProviderNameAlreadyExists() error {
	return &echo.HTTPError{
		Code:    http.StatusForbidden,
		Message: "provider with the same name already exists",
	}
}

func ErrProviderNotFound() error {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: "Model provider not found",
	}
}

func ErrMissDatasetID() error {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Failed to get dataset ID",
	}
}

func ErrMissFile() error {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Failed to get file from form data",
	}
}

func ErrFileNumberLimited() error {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Maximum 5 files allowed per upload",
	}
}

func ErrNoFileUploaded() error {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: "No files uploaded",
	}
}

func ErrFileNotFound() error {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: "File not found",
	}
}
