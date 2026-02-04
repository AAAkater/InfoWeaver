package v1

import (
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/response"
	"server/utils"

	"github.com/labstack/echo/v5"
)

func SetDatasetRouter(e *echo.Echo) {

	datasetRouterGroup := e.Group(config.API_V1+"/dataset", middleware.TokenMiddleware())

	datasetHandler := &datasetApi{}
	datasetRouterGroup.POST("/create", datasetHandler.createDataset)
	datasetRouterGroup.GET("/list", datasetHandler.listDatasets)
	datasetRouterGroup.GET("/:id", datasetHandler.getDatasetInfo)
	datasetRouterGroup.POST("/update", datasetHandler.updateDatasetInfo)
	datasetRouterGroup.POST("/delete/:id", datasetHandler.deleteDataset)
}

type datasetApi struct{}

// createDataset godoc
// @Summary      Create Dataset
// @Description  Create a new dataset for the authenticated user
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        dataset body models.DatasetCreateReq true "Dataset creation request"
// @Success      200 {object} response.ResponseBase[models.DatasetCreateResp] "Dataset created successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request data"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, insufficient permissions"
// @Router       /dataset/create [post]
func (this *datasetApi) createDataset(ctx *echo.Context) error {
	req, err := utils.BindAndValidate[models.DatasetCreateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	if err := datasetService.CreateNewDataset(ctx.Request().Context(), req.Name, req.Description, currentUser.ID); err != nil {
		return response.FailWithMsg(ctx, err.Error())
	}

	return response.Ok(ctx)
}

// listDatasets godoc
// @Summary      List Datasets
// @Description  List all datasets owned by the authenticated user
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Success      200 {object} response.ResponseBase[models.DatasetListResp] "List of datasets"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, insufficient permissions"
// @Router       /dataset/list [get]
func (this *datasetApi) listDatasets(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	total, datasets, err := datasetService.ListDatasetsByOwnerID(ctx.Request().Context(), currentUser.ID)
	if err != nil {
		return response.FailWithMsg(ctx, err.Error())
	}

	return response.OkWithData(ctx, models.DatasetListResp{
		Total:    total,
		Datasets: datasets,
	})
}

// getDatasetInfo godoc
// @Summary      Get Dataset
// @Description  Get a specific dataset by ID
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        id path int true "Dataset ID"
// @Success      200 {object} response.ResponseBase[models.DatasetInfo] "Dataset details"
// @Failure      400 {object} response.ResponseBase[any] "Invalid dataset ID"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, insufficient permissions"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/{id} [get]
func (this *datasetApi) getDatasetInfo(ctx *echo.Context) error {
	req, err := utils.BindAndValidate[models.DatasetInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}
	// Get user ID from token context
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return response.NoAuthWithMsg("Invalid user authentication")
	}

	dbDataset, err := datasetService.GetDatasetInfoByID(ctx.Request().Context(), req.ID, userID)
	if err != nil {
		return response.FailWithMsg(ctx, err.Error())
	}

	return response.OkWithData(ctx, dbDataset)
}

// updateDatasetInfo godoc
// @Summary      Update Dataset
// @Description  Update an existing dataset
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        dataset body models.DatasetUpdateReq true "Dataset update request"
// @Success      200 {object} response.ResponseBase[models.DatasetUpdateResp] "Dataset updated successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request data"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/update [post]
func (this *datasetApi) updateDatasetInfo(ctx *echo.Context) error {
	req, err := utils.BindAndValidate[models.DatasetUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	if err := datasetService.UpdateDataset(ctx.Request().Context(), req.ID, currentUser.ID, req.Name, req.Description); err != nil {
		return response.FailWithMsg(ctx, err.Error())
	}

	return response.Ok(ctx)
}

// deleteDataset godoc
// @Summary      Delete Dataset
// @Description  Delete a dataset by ID
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        id path int true "Dataset ID"
// @Success      200 {object} response.ResponseBase[models.DatasetDeleteResp] "Dataset deleted successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid dataset ID"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/delete/{id} [post]
func (this *datasetApi) deleteDataset(ctx *echo.Context) error {
	req, err := utils.BindAndValidate[models.DatasetUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}
	if err := datasetService.DeleteDataset(ctx.Request().Context(), req.ID, currentUser.ID); err != nil {
		return response.FailWithMsg(ctx, err.Error())
	}

	return response.Ok(ctx)
}
