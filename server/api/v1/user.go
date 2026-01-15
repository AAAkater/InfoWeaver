package v1

import (
	"server/config"
	"server/models/response"

	"github.com/labstack/echo/v4"
)

func SetUserRouter(e *echo.Echo) {

	userRouterGroup := e.Group(config.API_V1 + "/user")
	userHandler := &userApi{}
	userRouterGroup.GET("/info", userHandler.getUserInfo)

}

type userApi struct{}

func (user *userApi) getUserInfo(c echo.Context) error {
	return response.OkWithData(c, "asfas")
}
