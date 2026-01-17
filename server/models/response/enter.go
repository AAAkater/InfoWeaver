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

func NotFound(msg ...string) error {
	if len(msg) > 0 {
		return echo.NewHTTPError(http.StatusNotFound, msg[0])
	}
	return echo.NewHTTPError(http.StatusNotFound)
}

func NoAuth(msg ...string) error {
	if len(msg) > 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, msg[0])
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func BadRequest(msg ...string) error {
	if len(msg) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, msg[0])
	}
	return echo.NewHTTPError(http.StatusBadRequest)
}
