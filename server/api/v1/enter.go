package v1

import "server/service"

var (
	userService     = service.UserServiceApp
	fileService     = service.FileServiceApp
	datasetService  = service.DatasetServiceApp
	providerService = service.ProviderServiceApp
)
