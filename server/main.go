package main

import (
	"server/api"
	"server/config"
	"server/db"
	"server/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	config.VP, config.Settings = config.InitViper(config.DEFAULT_ENV_FILENAME)
	utils.Logger = utils.InitLogger(config.Settings.SYSTEM_IS_DEV)
	defer utils.Logger.Sync()
	db.InitAllDB()
	e := echo.New()

	api.InitRouter(e)

	server_port := config.Settings.GetServerPort()
	utils.Logger.Fatal(e.Start(server_port))
}
