package v1

import (
	"server/service"
	"server/utils"
)

var Logger = utils.Logger
var (
	userService     = service.UserServiceApp
	fileService     = service.FileServiceApp
	datasetService  = service.DatasetServiceApp
	providerService = service.ProviderServiceApp
)
