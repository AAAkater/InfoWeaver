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
	userRouterGroup.POST("/login", userHandler.login)
}

type userApi struct{}

func (this *userApi) register(ctx echo.Context) error {

	req, err := utils.BindAndValidate[request.RegisterBody](ctx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return response.BadRequest(err.Error())
	}

	if err := userService.CreateNewUser(ctx.Request().Context(), req.Username, req.Password, req.Email); err != nil {
		utils.Logger.Error(err)
		return response.NoAuth(err.Error())
	}
	return response.Ok(ctx)
}

func (this *userApi) login(ctx echo.Context) error {

	req, err := utils.BindAndValidate[request.LoginBody](ctx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return response.BadRequest(err.Error())
	}

	db_user, err := userService.GetUserByUsername(ctx.Request().Context(), req.Username)

	if err != nil && db_user == nil {
		utils.Logger.Error(err.Error())
		return response.NoAuth(err.Error())
	}

	if !utils.BcryptCheck(req.Password, db_user.Password) {
		return response.NoAuth("wrong password")
	}

	token, _ := utils.JwtTool().CreateToken(db_user.ID, db_user.Username, db_user.Role == "admin")

	return response.OkWithData(ctx, token)
}
