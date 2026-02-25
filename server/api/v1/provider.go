package v1

import (
	"errors"
	"fmt"
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/common/response"
	"server/service"
	"server/utils"

	"github.com/labstack/echo/v5"
)

func SetProviderRouter(e *echo.Echo) {
	providerRouterGroup := e.Group(config.API_V1+"/provider", middleware.TokenMiddleware())
	providerHandler := &providerApi{}
	providerRouterGroup.POST("", providerHandler.createProvider)
	providerRouterGroup.GET("/list", providerHandler.getAllProviders)
	providerRouterGroup.GET("/info/:provider_id", providerHandler.getProviderByID)
	providerRouterGroup.POST("/update", providerHandler.updateProvider)
	providerRouterGroup.POST("/delete/:provider_id", providerHandler.deleteProvider)
}

type providerApi struct{}

// createProvider godoc
// @Summary      Create Provider
// @Description  Create a new provider
// @Tags         Provider
// @Accept       json
// @Produce      json
// @Param        body body models.ProviderCreateReq true "Create Provider Request Body"
// @Success      200 {object} response.ResponseBase[any] "Provider created successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden: Provider name already exists"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
// @Router       /provider [post]
func (this *providerApi) createProvider(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	args, err := utils.BindAndValidate[models.ProviderCreateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := providerService.CreateProvider(ctx.Request().Context(), currentUser.ID, args.Name, args.BaseURL, args.APIKey, args.Mode); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrDuplicatedKey):
		utils.Logger.Error(err)
		return response.ForbiddenWithMsg("Provider Name:" + args.Name + " already exists")
	default:
		utils.Logger.Error(err)
		return response.FailWithMsg(ctx, "Failed to create provider")
	}
}

// getAllProviders godoc
// @Summary      Get All Providers
// @Description  Get a list of all providers
// @Tags         Provider
// @Accept       json
// @Produce      json
// @Success      200 {object} response.ResponseBase[models.ProviderListResp] "Providers retrieved successfully"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
// @Router       /provider/list [get]
func (this *providerApi) getAllProviders(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	cows, providers, err := providerService.GetAllProviders(ctx.Request().Context(), currentUser.ID)
	if err != nil {
		utils.Logger.Error(err)
		return response.FailWithMsg(ctx, "Failed to get providers")
	}
	return response.OkWithData(ctx, models.ProviderListResp{
		Total:     cows,
		Providers: providers,
	})

}

// getProviderByID godoc
// @Summary      Get Provider by ID
// @Description  Get a provider by its ID
// @Tags         Provider
// @Accept       json
// @Produce      json
// @Param        provider_id path int true "Provider ID"
// @Success      200 {object} response.ResponseBase[models.ProviderInfo] "Provider retrieved successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid provider ID"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden: Unauthorized access to the provider"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
// @Router       /provider/info/{provider_id} [get]
func (this *providerApi) getProviderByID(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	args, err := utils.BindAndValidate[models.ProviderInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	provider, err := providerService.GetProviderByID(ctx.Request().Context(), args.ID, currentUser.ID)
	switch err {
	case nil:
		return response.OkWithData(ctx, provider)
	case service.ErrNotFound:
		return response.ForbiddenWithMsg(fmt.Sprintf("Unauthorized access to the model provider: %d", args.ID))
	default:
		utils.Logger.Errorf("Failed to model provider with ID %d: %v", args.ID, err)
		return response.FailWithMsg(ctx, "Unknown error")
	}
}

// updateProvider godoc
// @Summary      Update Provider
// @Description  Update an existing provider
// @Tags         Provider
// @Accept       json
// @Produce      json
// @Param        body body models.ProviderUpdateReq true "Update Provider Request Body"
// @Success      200 {object} response.ResponseBase[any] "Provider updated successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden: Provider does not exist"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
// @Router       /provider/update [post]
func (this *providerApi) updateProvider(ctx *echo.Context) error {

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	args, err := utils.BindAndValidate[models.ProviderUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := providerService.UpdateProvider(
		ctx.Request().Context(),
		args.ID,
		currentUser.ID,
		args.Name,
		args.BaseURL,
		args.APIKey,
		args.Mode); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		return response.ForbiddenWithMsg("Provider does not exist")
	default:
		utils.Logger.Error(err)
		return response.FailWithMsg(ctx, "Failed to update provider")
	}
}

// deleteProvider godoc
// @Summary      Delete Provider
// @Description  Delete a provider by ID
// @Tags         Provider
// @Accept       json
// @Produce      json
// @Param        provider_id path int true "Provider ID"
// @Param        body body models.ProviderInfoReq true "Delete Provider Request Body"
// @Success      200 {object} response.ResponseBase[any] "Provider deleted successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
// @Router       /provider/delete/{provider_id} [post]
func (this *providerApi) deleteProvider(ctx *echo.Context) error {

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	args, err := utils.BindAndValidate[models.ProviderInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	if _, err := providerService.GetProviderByID(ctx.Request().Context(), args.ID, currentUser.ID); err != nil {
		utils.Logger.Error(err)
		return response.ForbiddenWithMsg(fmt.Sprintf("Unauthorized access to delete provider:%d", args.ID))
	}

	if err := service.ProviderServiceApp.DeleteProvider(ctx.Request().Context(), args.ID); err != nil {
		utils.Logger.Error(err)
		return response.FailWithMsg(ctx, "Failed to delete provider")
	}

	return response.Ok(ctx)
}
