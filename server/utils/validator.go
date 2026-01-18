package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validatorInstance = validator.New()

func BindAndValidate[T any](ctx echo.Context) (*T, error) {

	var data T

	if err := ctx.Bind(&data); err != nil {
		return nil, err
	}

	if err := validatorInstance.Struct(data); err != nil {
		return nil, err
	}

	return &data, nil
}
