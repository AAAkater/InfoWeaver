package main

import (
	"server/api"
	"server/config"
	"server/db"
	_ "server/docs"
	"server/middleware"
	"server/utils"

	"github.com/labstack/echo/v5"
)

// @title InfoWeaver API
// @version 1.0
// @description Multi-modal RAG system API
// @BasePath /api/v1
func main() {
	config.VP, config.Settings = config.InitViper(config.DEFAULT_ENV_FILENAME)
	db.InitAllDB()
	e := echo.New()

	middleware.InitMiddleWares(e)
	api.InitRouter(e)

	server_port := config.Settings.GetServerPort()
	utils.Logger.Fatal(e.Start(server_port))
}
