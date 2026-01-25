package v1

import (
	"server/config"
	"server/middleware"
	"server/models/request"
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
	userRouterGroup.POST("/restPassword", userHandler.resetUserPassword, middleware.TokenMiddleware())
}

type userApi struct{}

// register godoc
// @Summary      User Register
// @Description  Register a new user account;
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body request.RegisterReq true "Register Request"
// @Success      200 {object} map[string]interface{} "Register successful"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /user/register [post]
func (this *userApi) register(ctx *echo.Context) error {

	newUserInfo, err := utils.BindAndValidate[request.RegisterReq](ctx)
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

// login godoc
// @Summary      User Login
// @Description  Login with username and password
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body body request.LoginReq true "Login Request"
// @Success      200 {object} map[string]interface{} "Login successful"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      401 {object} map[string]interface{} "Invalid credentials"
// @Router       /user/login [post]
func (this *userApi) login(ctx *echo.Context) error {

	userInfo, err := utils.BindAndValidate[request.LoginReq](ctx)
	if err != nil {
		utils.Logger.Error(err.Error())
		return response.BadRequest()
	}

	db_user, err := userService.GetUserInfoByUsername(ctx.Request().Context(), userInfo.Username)

	if err != nil && db_user == nil {
		utils.Logger.Error(err.Error())
		return response.NoAuthWithMsg("User not found")
	}

	if !utils.BcryptCheck(userInfo.Password, db_user.Password) {
		return response.NoAuthWithMsg("Invalid password")
	}

	token, _ := utils.CreateToken(db_user.ID, db_user.Role == "admin")

	return response.OkWithData(ctx, response.TokenResult{
		Type:  "Bearer",
		Token: token,
	})
}

// getUserInfo godoc
// @Summary      Get User Info
// @Description  Get current user information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} response.UserInfoResult "User info"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Router       /user/info [get]
func (this *userApi) getUserInfo(ctx *echo.Context) error {

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg("token invalid or expired")
	}

	db_user, err := userService.GetUserInfoByID(ctx.Request().Context(), currentUser.ID)
	if err != nil || db_user == nil {
		return response.NoAuthWithMsg("user not found")
	}
	return response.OkWithData(ctx, response.UserInfoResult{
		ID:       db_user.ID,
		Username: db_user.Username,
		Email:    db_user.Email,
		Role:     db_user.Role,
	})
}

// resetUserPassword godoc
// @Summary      Reset User Password
// @Description  Reset the current user's password
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        body body request.ResetPasswordReq true "Reset Password Request"
// @Success      200 {object} map[string]interface{} "Password reset successful"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      401 {object} map[string]interface{} "Unauthorized"
// @Router       /user/restPassword [post]
func (this *userApi) resetUserPassword(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg("token invalid or expired")
	}

	password_info, err := utils.BindAndValidate[request.ResetPasswordReq](ctx)
	if err != nil {
		return response.BadRequest()
	}

	if err := userService.ResetUserPassword(ctx.Request().Context(), currentUser.ID, password_info.NewPassword); err != nil {
		return response.NoAuthWithMsg("Failed to reset password")
	}

	return response.Ok(ctx)
}

func (this *userApi) resetUserInfo(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg("token invalid or expired")
	}

}
