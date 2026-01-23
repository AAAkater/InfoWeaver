package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const (
	SUCCESS = iota
	ERROR
)

type ResponseBase[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func Ok(ctx *echo.Context) error {
	return ctx.JSON(http.StatusOK,
		ResponseBase[any]{
			SUCCESS,
			"ok",
			nil,
		})
}

func OkWithData[T any](ctx *echo.Context, data T) error {
	return ctx.JSON(
		http.StatusOK,
		ResponseBase[T]{
			SUCCESS,
			"ok",
			data,
		})
}

func Fail(ctx *echo.Context) error {
	return ctx.JSON(http.StatusOK,
		ResponseBase[any]{
			ERROR,
			"ok",
			nil,
		})
}

func FailWithMsg(ctx *echo.Context, msg string) error {
	return ctx.JSON(
		http.StatusOK,
		ResponseBase[any]{
			ERROR,
			msg,
			nil,
		})
}

func NotFound() error {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

func NotFoundWithMsg(msg string) error {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func NoAuth() error {
	return &echo.HTTPError{
		Code:    http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
	}
}

func NoAuthWithMsg(msg string) error {
	return &echo.HTTPError{
		Code:    http.StatusUnauthorized,
		Message: msg,
	}
}

func BadRequest() error {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}
}

func BadRequestWithMsg(msg string) error {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}
