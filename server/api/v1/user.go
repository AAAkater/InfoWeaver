package v1

import (
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/response"
	"server/utils"

	"github.com/labstack/echo/v5"
)

func SetUserRouter(e *echo.Echo) {

	userRouterGroup := e.Group(config.API_V1 + "/user")
	userHandler := &userApi{}
	userRouterGroup.POST("/register", userHandler.register)
	userRouterGroup.POST("/login", userHandler.login)
	userRouterGroup.GET("/info", userHandler.getUserInfo, middleware.TokenMiddleware())
	userRouterGroup.POST("/resetPassword", userHandler.resetUserPassword, middleware.TokenMiddleware())
	userRouterGroup.POST("/updateInfo", userHandler.updateUserInfo, middleware.TokenMiddleware())
}

type userApi struct{}

// register godoc
// @Summary      User Register
// @Description  Register a new user account with username, password and email
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body models.RegisterReq true "Register Request Body"
// @Success      200 {object} response.ResponseBase[any] "Register successful"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "User already exists or registration failed"
// @Router       /user/register [post]
func (this *userApi) register(ctx *echo.Context) error {

	newUserInfo, err := utils.BindAndValidate[models.RegisterReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	if err := userService.CreateNewUser(ctx.Request().Context(), newUserInfo.Username, newUserInfo.Password, newUserInfo.Email); err != nil {
		utils.Logger.Error(err)
		return response.NoAuthWithMsg(err.Error())
	}
	return response.Ok(ctx)
}

// login godoc
// @Summary      User Login
// @Description  Login with username and password to get authentication token
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body models.UserLoginReq true "Login Request Body"
// @Success      200 {object} response.ResponseBase[models.UserLoginResp] "Login successful, returns Bearer token"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Invalid password or credentials"
// @Failure      404 {object} response.ResponseBase[any] "User not found"
// @Router       /user/login [post]
func (this *userApi) login(ctx *echo.Context) error {

	userInfo, err := utils.BindAndValidate[models.UserLoginReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	dbUser, err := userService.GetUserInfoByUsername(ctx.Request().Context(), userInfo.Username)

	if err != nil && dbUser == nil {
		return response.NotFoundWithMsg(err.Error())
	}

	if !utils.BcryptCheck(userInfo.Password, dbUser.Password) {
		return response.NoAuthWithMsg("Invalid password")
	}

	token, err := utils.CreateToken(dbUser.ID, dbUser.Role == "admin")
	if err != nil {
		utils.Logger.Error(err)
		return response.NoAuthWithMsg("Failed to generate token")
	}

	return response.OkWithData(ctx, models.UserLoginResp{
		Type:  "Bearer",
		Token: token,
	})
}

// getUserInfo godoc
// @Summary      Get User Info
// @Description  Get current authenticated user information including id, username, email and role
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.ResponseBase[models.UserInfoResp] "User information retrieved successfully"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized or invalid token"
// @Failure      404 {object} response.ResponseBase[any] "User not found"
// @Router       /user/info [get]
func (this *userApi) getUserInfo(ctx *echo.Context) error {

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	dbUser, err := userService.GetUserInfoByID(ctx.Request().Context(), currentUser.ID)
	if err != nil || dbUser == nil {
		return response.NoAuthWithMsg(err.Error())
	}
	return response.OkWithData(ctx, models.UserInfoResp{
		ID:       dbUser.ID,
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Role:     dbUser.Role,
	})
}

// resetUserPassword godoc
// @Summary      Reset User Password
// @Description  Reset the current authenticated user's password
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        body body models.ResetPasswordReq true "Reset Password Request Body"
// @Success      200 {object} response.ResponseBase[any] "Password reset successful"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized or invalid token"
// @Router       /user/resetPassword [post]
func (this *userApi) resetUserPassword(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	req, err := utils.BindAndValidate[models.ResetPasswordReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	if err := userService.ResetUserPassword(ctx.Request().Context(), currentUser.ID, req.FirstPassword); err != nil {
		return response.NoAuthWithMsg("Failed to reset password")
	}

	return response.Ok(ctx)
}

// updateUserInfo godoc
// @Summary      Update User Info
// @Description  Update the current authenticated user's username and/or email
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        body body models.UpdateUserInfoReq true "Update User Info Request Body"
// @Success      200 {object} response.ResponseBase[any] "User information updated successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized or invalid token"
// @Router       /user/updateInfo [post]
func (this *userApi) updateUserInfo(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	newUserInfo, err := utils.BindAndValidate[models.UpdateUserInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	if err := userService.UpdateUserInfo(ctx.Request().Context(), currentUser.ID, newUserInfo.Username, newUserInfo.Email); err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	return response.Ok(ctx)

}
