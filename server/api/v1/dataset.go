package v1

import (
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/response"
	"server/service"
	"server/utils"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func SetDatasetRouter(e *echo.Echo) {

	datasetRouterGroup := e.Group(config.API_V1+"/dataset", middleware.TokenMiddleware())

	datasetHandler := &datasetApi{}
	datasetRouterGroup.POST("/create", datasetHandler.createDataset)
	datasetRouterGroup.GET("/list", datasetHandler.listDatasets)
	datasetRouterGroup.GET("/:id", datasetHandler.getDataset)
	datasetRouterGroup.PUT("/update", datasetHandler.updateDataset)
	datasetRouterGroup.DELETE("/:id", datasetHandler.deleteDataset)

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

	if err := service.DatasetServiceApp.CreateNewDataset(ctx.Request().Context(), req.Name, req.Description, currentUser.ID); err != nil {
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
// @Router       /dataset/list [get]
func (this *datasetApi) listDatasets(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	datasets, err := service.DatasetServiceApp.ListDatasetsByOwner(ctx.Request().Context(), currentUser.ID)
	if err != nil {
		return response.FailWithMsg(ctx, err.Error())
	}

	// Convert to response format
	respDatasets := make([]models.DatasetResp, 0, len(datasets))
	for _, dataset := range datasets {
		respDatasets = append(respDatasets, models.DatasetResp{
			ID:          dataset.ID,
			Name:        dataset.Name,
			Description: dataset.Description,
			OwnerID:     dataset.OwnerID,
			CreatedAt:   dataset.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   dataset.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return response.OkWithData(ctx, models.DatasetListResp{Datasets: respDatasets})
}

// getDataset godoc
// @Summary      Get Dataset
// @Description  Get a specific dataset by ID
// @Tags         Dataset
// @Accept       json
// @Produce      json
// @Param        id path int true "Dataset ID"
// @Success      200 {object} response.ResponseBase[models.DatasetResp] "Dataset details"
// @Failure      400 {object} response.ResponseBase[any] "Invalid dataset ID"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Router       /dataset/{id} [get]
func (this *datasetApi) getDataset(ctx *echo.Context) error {
	id, err := ctx.ParamInt("id")
	if err != nil {
		return response.BadRequestWithMsg("Invalid dataset ID")
	}

	// Get user ID from token context
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return response.NoAuthWithMsg("Invalid user authentication")
	}

	dataset, err := service.DatasetServiceApp.GetDatasetByID(ctx.Request().Context(), uint(id), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFoundWithMsg("Dataset not found")
		}
		utils.Logger.Error(err)
		return response.FailWithMsg("Failed to get dataset")
	}

	respDataset := models.DatasetResp{
		ID:          dataset.ID,
		Name:        dataset.Name,
		Description: dataset.Description,
		OwnerID:     dataset.OwnerID,
		CreatedAt:   dataset.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   dataset.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response.OkWithData(ctx, respDataset)
}

// updateDataset godoc
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
// @Router       /dataset/update [put]
func (this *datasetApi) updateDataset(ctx *echo.Context) error {
	req, err := utils.BindAndValidate[models.DatasetUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get user ID from token context
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return response.NoAuthWithMsg("Invalid user authentication")
	}

	// Get dataset ID from query parameter or request body
	id, err := ctx.ParamInt("id")
	if err != nil {
		// Try to get from request body if not in URL
		id = int(req.ID)
	}

	if id <= 0 {
		return response.BadRequestWithMsg("Invalid dataset ID")
	}

	if err := service.DatasetServiceApp.UpdateDataset(ctx.Request().Context(), uint(id), userID, req.Name, req.Description); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFoundWithMsg("Dataset not found")
		}
		utils.Logger.Error(err)
		return response.FailWithMsg("Failed to update dataset")
	}

	return response.OkWithData(ctx, models.DatasetUpdateResp{ID: uint(id)})
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
// @Router       /dataset/{id} [delete]
func (this *datasetApi) deleteDataset(ctx *echo.Context) error {
	id, err := ctx.ParamInt("id")
	if err != nil {
		return response.BadRequestWithMsg("Invalid dataset ID")
	}

	// Get user ID from token context
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return response.NoAuthWithMsg("Invalid user authentication")
	}

	if err := service.DatasetServiceApp.DeleteDataset(ctx.Request().Context(), uint(id), userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.NotFoundWithMsg("Dataset not found")
		}
		utils.Logger.Error(err)
		return response.FailWithMsg("Failed to delete dataset")
	}

	return response.OkWithData(ctx, models.DatasetDeleteResp{ID: uint(id)})
}
