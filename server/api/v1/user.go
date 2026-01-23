package v1

import (
	"server/config"
	"server/middleware"
	"server/models/request"
	"server/models/response"
	"server/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func SetUserRouter(e *echo.Echo) {

	userRouterGroup := e.Group(config.API_V1 + "/user")
	userHandler := &userApi{}
	userRouterGroup.POST("/register", userHandler.register)
	userRouterGroup.POST("/login", userHandler.login)
	userRouterGroup.GET("/info", userHandler.getUserInfo, middleware.TokenMiddleware())
}

type userApi struct{}

func (this *userApi) register(ctx *echo.Context) error {

	newUserInfo, err := utils.BindAndValidate[request.RegisterBody](ctx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return response.BadRequest()
	}

	if err := userService.CreateNewUser(ctx.Request().Context(), newUserInfo.Username, newUserInfo.Password, newUserInfo.Email); err != nil {
		utils.Logger.Error(err)
		return response.NoAuth()
	}
	return response.Ok(ctx)
}

func (this *userApi) login(ctx *echo.Context) error {

	userInfo, err := utils.BindAndValidate[request.LoginBody](ctx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return response.BadRequest()
	}

	db_user, err := userService.GetUserByUsername(ctx.Request().Context(), userInfo.Username)

	if err != nil && db_user == nil {
		utils.Logger.Error(err.Error())
		return response.NoAuth()
	}

	if !utils.BcryptCheck(userInfo.Password, db_user.Password) {
		return response.NoAuth()
	}

	token, _ := utils.CreateToken(db_user.ID, db_user.Role == "admin")

	return response.OkWithData(ctx, token)
}

func (this *userApi) getUserInfo(ctx *echo.Context) error {

	user, err := echo.ContextGet[*jwt.Token](ctx, "user")
	if err != nil {
		return response.NoAuthWithMsg("token invalid or expired")
	}

	claims := user.Claims.(utils.JwtCustomClaims)
	userID := claims.UserID
	return response.OkWithData(ctx, userID)
}
