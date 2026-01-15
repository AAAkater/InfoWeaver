package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	ERROR   = 7
	SUCCESS = 0
)

type ResponseBase[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func Ok(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, ResponseBase[any]{
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
