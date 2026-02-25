package v1

import (
	"errors"
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/common/response"
	"server/service"
	"server/utils"

	"github.com/labstack/echo/v5"
)

func SetDatasetRouter(e *echo.Echo) {

	datasetRouterGroup := e.Group(config.API_V1+"/dataset", middleware.TokenMiddleware())

	datasetHandler := &datasetApi{}
	datasetRouterGroup.POST("/create", datasetHandler.createDataset)
	datasetRouterGroup.GET("", datasetHandler.listDatasets)
	datasetRouterGroup.GET("/:dataset_id", datasetHandler.getDatasetInfo)
	datasetRouterGroup.POST("/update", datasetHandler.updateDatasetInfo)
	datasetRouterGroup.POST("/delete/:dataset_id", datasetHandler.deleteDataset)
}

type datasetApi struct{}

// createDataset godoc
// @Summary      Create Dataset
// @Description  Create a new dataset for the authenticated user
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        dataset body models.DatasetCreateReq true "Dataset creation request"
// @Success      200 {object} response.ResponseBase[any] "Dataset created successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request data"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, insufficient permissions"
// @Router       /dataset/create [post]
func (this *datasetApi) createDataset(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.DatasetCreateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := datasetService.CreateNewDataset(ctx.Request().Context(), args.Icon, args.Name, args.Description, currentUser.ID); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrDuplicatedKey):
		return response.ErrDatasetNameAlreadyExists()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// listDatasets godoc
// @Summary      List Datasets
// @Description  List all datasets owned by the authenticated user. If name query parameter is provided, filter datasets by name.
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        name query string false "Dataset name to filter by"
// @Success      200 {object} response.ResponseBase[models.DatasetListResp] "List of datasets"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request data"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, insufficient permissions"
// @Router       /dataset [get]
func (this *datasetApi) listDatasets(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	// Parse optional name query parameter
	args, err := utils.BindAndValidate[models.DatasetListReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	var total int64
	var datasets []models.DatasetInfo

	if args.Name != "" {
		total, datasets, err = datasetService.ListDatasetsByName(ctx.Request().Context(), currentUser.ID, args.Name)
	} else {
		total, datasets, err = datasetService.ListDatasetsByOwnerID(ctx.Request().Context(), currentUser.ID)
	}

	switch err {
	case nil:
		return response.OkWithData(ctx, models.DatasetListResp{
			Total:    total,
			Datasets: datasets,
		})
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// getDatasetInfo godoc
// @Summary      Get Dataset
// @Description  Get a specific dataset by ID
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        dataset_id path int true "Dataset ID"
// @Success      200 {object} response.ResponseBase[models.DatasetInfo] "Dataset details"
// @Failure      400 {object} response.ResponseBase[any] "Invalid dataset ID"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, insufficient permissions"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/{dataset_id} [get]
func (this *datasetApi) getDatasetInfo(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.DatasetInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch dbDataset, err := datasetService.GetDatasetInfoByID(ctx.Request().Context(), args.ID, currentUser.ID); {
	case err == nil:
		return response.OkWithData(ctx, dbDataset)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrDatasetNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// updateDatasetInfo godoc
// @Summary      Update Dataset
// @Description  Update an existing dataset
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        dataset body models.DatasetUpdateReq true "Dataset update request"
// @Success      200 {object} response.ResponseBase[any] "Dataset updated successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request data"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/update [post]
func (this *datasetApi) updateDatasetInfo(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.DatasetUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := datasetService.UpdateDataset(ctx.Request().Context(), args.ID, currentUser.ID, args.Icon, args.Name, args.Description); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrDatasetNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// deleteDataset godoc
// @Summary      Delete Dataset
// @Description  Delete a dataset by ID
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        dataset_id path int true "Dataset ID"
// @Success      200 {object} response.ResponseBase[any] "Dataset deleted successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid dataset ID"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/delete/{dataset_id} [post]
func (this *datasetApi) deleteDataset(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.DatasetInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := datasetService.DeleteDataset(ctx.Request().Context(), args.ID, currentUser.ID); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		Logger.Error(err)
		return response.ErrDatasetNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}
