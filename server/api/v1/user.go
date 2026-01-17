package v1

import (
	"server/config"
	"server/models/request"
	"server/models/response"
	"server/utils"

	"github.com/labstack/echo/v4"
)

func SetUserRouter(e *echo.Echo) {

	userRouterGroup := e.Group(config.API_V1 + "/user")
	userHandler := &userApi{}
	userRouterGroup.POST("/register", userHandler.register)
}

type userApi struct{}

func (this *userApi) register(ctx echo.Context) error {

	var req request.RegisterBody
	if err := ctx.Bind(&req); err != nil {
		utils.Logger.Error(err.Error())
		return response.BadRequest()
	}

	if err := userService.CreateNewUser(ctx.Request().Context(), req.Username, req.Password, req.Email); err != nil {
		utils.Logger.Error(err)
		return response.NoAuth(err.Error())
	}
	return response.Ok(ctx)
}
