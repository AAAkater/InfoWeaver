package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

func Ok(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK,
		ResponseBase[any]{
			SUCCESS,
			"ok",
			nil,
		})
}

func OkWithData[T any](ctx echo.Context, data T) error {
	return ctx.JSON(
		http.StatusOK,
		ResponseBase[T]{
			SUCCESS,
			"ok",
			data,
		})
}

func Fail(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK,
		ResponseBase[any]{
			ERROR,
			"ok",
			nil,
		})
}

func FailWithMsg(ctx echo.Context, msg string) error {
	return ctx.JSON(
		http.StatusOK,
		ResponseBase[any]{
			ERROR,
			msg,
			nil,
		})
}

func NotFound(ctx echo.Context, msg string) error {
	return ctx.JSON(
		http.StatusNotFound,
		ResponseBase[any]{
			SUCCESS,
			msg,
			nil,
		})
}

func NoAuth(ctx echo.Context, msg string) error {
	return ctx.JSON(
		http.StatusUnauthorized,
		ResponseBase[any]{
			SUCCESS,
			msg,
			nil,
		})
}
