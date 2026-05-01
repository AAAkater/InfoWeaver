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

func SetMcpRouter(e *echo.Echo) {
	mcpRouterGroup := e.Group(config.API_V1+"/mcp", middleware.TokenMiddleware())
	mcpHandler := &mcpApi{}
	mcpRouterGroup.POST("", mcpHandler.createMcp)
	mcpRouterGroup.GET("/list", mcpHandler.getAllMcps)
	mcpRouterGroup.GET("/info/:mcp_id", mcpHandler.getMcpInfoByID)
	mcpRouterGroup.POST("/update", mcpHandler.updateMcp)
	mcpRouterGroup.POST("/delete/:mcp_id", mcpHandler.deleteMcp)
}

type mcpApi struct{}

// createMcp godoc
//
//	@Summary		Create MCP Server
//	@Description	Create a new MCP server configuration
//	@Tags			MCP
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.McpCreateReq			true	"Create MCP Request Body"
//	@Success		200		{object}	response.ResponseBase[any]	"MCP server created successfully"
//	@Failure		400		{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403		{object}	response.ResponseBase[any]	"MCP server name already exists"
//	@Failure		500		{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/mcp [post]
func (a *mcpApi) createMcp(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.McpCreateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Check if an MCP server with the same name already exists for this owner
	if exist, err := mcpService.CheckMcpExistsByName(ctx.Request().Context(), currentUser.ID, args.Name); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	} else if exist {
		return response.ErrMcpNameAlreadyExists()
	}

	enabled := true
	if args.Enabled != nil {
		enabled = *args.Enabled
	}

	if err := mcpService.CreateMcp(
		ctx.Request().Context(),
		currentUser.ID,
		args.Name,
		args.Transport,
		args.Command,
		args.Args,
		args.URL,
		args.Headers,
		args.EnvVars,
		enabled,
	); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.Ok(ctx)
}

// getAllMcps godoc
//
//	@Summary		Get All MCP Servers
//	@Description	Get a list of all MCP server configurations
//	@Tags			MCP
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.ResponseBase[models.McpListResp]	"MCP servers retrieved successfully"
//	@Failure		401	{object}	response.ResponseBase[any]					"Invalid or expired token"
//	@Failure		500	{object}	response.ResponseBase[any]					"Internal server error"
//	@Router			/mcp/list [get]
func (a *mcpApi) getAllMcps(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	total, mcps, err := mcpService.GetAllMcps(ctx.Request().Context(), currentUser.ID)
	if err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.OkWithData(ctx, models.McpListResp{
		Total: total,
		Mcps:  mcps,
	})
}

// getMcpInfoByID godoc
//
//	@Summary		Get MCP Server by ID
//	@Description	Get an MCP server configuration by its ID
//	@Tags			MCP
//	@Accept			json
//	@Produce		json
//	@Param			mcp_id	path		int										true	"MCP Server ID"
//	@Success		200		{object}	response.ResponseBase[models.McpInfo]	"MCP server retrieved successfully"
//	@Failure		400		{object}	response.ResponseBase[any]				"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]				"Invalid or expired token"
//	@Failure		404		{object}	response.ResponseBase[any]				"MCP server not found"
//	@Failure		500		{object}	response.ResponseBase[any]				"Internal server error"
//	@Router			/mcp/info/{mcp_id} [get]
func (a *mcpApi) getMcpInfoByID(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.McpInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	mcp, err := mcpService.GetMcpByID(ctx.Request().Context(), args.ID, currentUser.ID)
	switch err {
	case nil:
		return response.OkWithData(ctx, mcp)
	case service.ErrNotFound:
		return response.ErrMcpNotFound()
	default:
		Logger.Errorf("Failed to get MCP server with ID %d: %v", args.ID, err)
		return response.ErrUnknownError()
	}
}

// updateMcp godoc
//
//	@Summary		Update MCP Server
//	@Description	Update an existing MCP server configuration
//	@Tags			MCP
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.McpUpdateReq			true	"Update MCP Request Body"
//	@Success		200		{object}	response.ResponseBase[any]	"MCP server updated successfully"
//	@Failure		400		{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403		{object}	response.ResponseBase[any]	"MCP server name already exists"
//	@Failure		404		{object}	response.ResponseBase[any]	"MCP server not found"
//	@Failure		500		{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/mcp/update [post]
func (a *mcpApi) updateMcp(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.McpUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	enabled := true
	if args.Enabled != nil {
		enabled = *args.Enabled
	}

	switch err := mcpService.UpdateMcp(
		ctx.Request().Context(),
		args.ID,
		currentUser.ID,
		args.Name,
		args.Transport,
		args.Command,
		args.Args,
		args.URL,
		args.Headers,
		args.EnvVars,
		enabled,
	); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrMcpNotFound()
	case errors.Is(err, service.ErrDuplicatedKey):
		return response.ErrMcpNameAlreadyExists()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// deleteMcp godoc
//
//	@Summary		Delete MCP Server
//	@Description	Delete an MCP server configuration by ID
//	@Tags			MCP
//	@Accept			json
//	@Produce		json
//	@Param			mcp_id	path		int							true	"MCP Server ID"
//	@Success		200		{object}	response.ResponseBase[any]	"MCP server deleted successfully"
//	@Failure		400		{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		404		{object}	response.ResponseBase[any]	"MCP server not found"
//	@Failure		500		{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/mcp/delete/{mcp_id} [post]
func (a *mcpApi) deleteMcp(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.McpInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := mcpService.DeleteMcp(ctx.Request().Context(), args.ID, currentUser.ID); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrMcpNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}
