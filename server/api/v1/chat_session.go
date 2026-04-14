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

func SetChatSessionRouter(e *echo.Echo) {
	chatSessionRouterGroup := e.Group(config.API_V1+"/chat-session", middleware.TokenMiddleware())
	chatSessionHandler := &chatSessionApi{}
	chatSessionRouterGroup.POST("/create", chatSessionHandler.createChatSession)
	chatSessionRouterGroup.GET("", chatSessionHandler.listChatSessions)
	chatSessionRouterGroup.GET("/:session_id", chatSessionHandler.getChatSessionInfo)
	chatSessionRouterGroup.POST("/update", chatSessionHandler.updateChatSessionInfo)
	chatSessionRouterGroup.POST("/delete/:session_id", chatSessionHandler.deleteChatSession)
}

type chatSessionApi struct{}

// createChatSession godoc
//
//	@Summary		Create Chat Session
//	@Description	Create a new chat session for the authenticated user
//	@Tags			ChatSession
//	@Accept			json
//	@Produce		json
//	@Param			session	body		models.ChatSessionCreateReq	true	"Chat session creation request"
//	@Success		200		{object}	response.ResponseBase[any]	"Chat session created successfully"
//	@Failure		400		{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403		{object}	response.ResponseBase[any]	"Chat session title already exists"
//	@Failure		500		{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat-session/create [post]
func (this *chatSessionApi) createChatSession(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatSessionCreateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Check if chat session with same title already exists
	if exist, err := chatSessionService.CheckChatSessionExistsByTitle(ctx.Request().Context(), currentUser.ID, args.Title); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	} else if exist {
		return response.ErrChatSessionTitleAlreadyExists()
	}

	if err := chatSessionService.CreateNewChatSession(ctx.Request().Context(),
		args.Title,
		currentUser.ID); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.Ok(ctx)
}

// listChatSessions godoc
//
//	@Summary		List Chat Sessions
//	@Description	List all chat sessions owned by the authenticated user. If title query parameter is provided, filter sessions by title.
//	@Tags			ChatSession
//	@Accept			json
//	@Produce		json
//	@Param			title	query		string												false	"Chat session title to filter by"
//	@Success		200		{object}	response.ResponseBase[models.ChatSessionListResp]	"List of chat sessions"
//	@Failure		400		{object}	response.ResponseBase[any]							"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]							"Invalid or expired token"
//	@Failure		500		{object}	response.ResponseBase[any]							"Internal server error"
//	@Router			/chat-session [get]
func (this *chatSessionApi) listChatSessions(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	// Parse optional title query parameter
	args, err := utils.BindAndValidate[models.ChatSessionListReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	var total int64
	var sessions []models.ChatSessionInfo

	if args.Title != "" {
		total, sessions, err = chatSessionService.ListChatSessionsByTitle(ctx.Request().Context(), currentUser.ID, args.Title)
	} else {
		total, sessions, err = chatSessionService.ListChatSessionsByOwnerID(ctx.Request().Context(), currentUser.ID)
	}
	if err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.OkWithData(ctx, models.ChatSessionListResp{
		Total:    total,
		Sessions: sessions,
	})
}

// getChatSessionInfo godoc
//
//	@Summary		Get Chat Session
//	@Description	Get a specific chat session by ID
//	@Tags			ChatSession
//	@Accept			json
//	@Produce		json
//	@Param			session_id	path		int												true	"Chat session ID"
//	@Success		200			{object}	response.ResponseBase[models.ChatSessionInfo]	"Chat session details"
//	@Failure		400			{object}	response.ResponseBase[any]						"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]						"Invalid or expired token"
//	@Failure		404			{object}	response.ResponseBase[any]						"Chat session not found"
//	@Failure		500			{object}	response.ResponseBase[any]						"Internal server error"
//	@Router			/chat-session/{session_id} [get]
func (this *chatSessionApi) getChatSessionInfo(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatSessionInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch dbChatSession, err := chatSessionService.GetChatSessionInfoByID(ctx.Request().Context(), args.ID, currentUser.ID); {
	case err == nil:
		return response.OkWithData(ctx, dbChatSession)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrChatSessionNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// updateChatSessionInfo godoc
//
//	@Summary		Update Chat Session
//	@Description	Update an existing chat session
//	@Tags			ChatSession
//	@Accept			json
//	@Produce		json
//	@Param			session	body		models.ChatSessionUpdateReq	true	"Chat session update request"
//	@Success		200		{object}	response.ResponseBase[any]	"Chat session updated successfully"
//	@Failure		400		{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		404		{object}	response.ResponseBase[any]	"Chat session not found"
//	@Failure		500		{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat-session/update [post]
func (this *chatSessionApi) updateChatSessionInfo(ctx *echo.Context) error {
	// Get user ID from token context
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatSessionUpdateReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := chatSessionService.UpdateChatSession(ctx.Request().Context(), args.ID, currentUser.ID, args.Title); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrChatSessionNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}

// deleteChatSession godoc
//
//	@Summary		Delete Chat Session
//	@Description	Delete a chat session by ID
//	@Tags			ChatSession
//	@Accept			json
//	@Produce		json
//	@Param			session_id	path		int							true	"Chat session ID"
//	@Success		200			{object}	response.ResponseBase[any]	"Chat session deleted successfully"
//	@Failure		400			{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		404			{object}	response.ResponseBase[any]	"Chat session not found"
//	@Failure		500			{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat-session/delete/{session_id} [post]
func (this *chatSessionApi) deleteChatSession(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatSessionInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch err := chatSessionService.DeleteChatSession(ctx.Request().Context(), args.ID, currentUser.ID); {
	case err == nil:
		return response.Ok(ctx)
	case errors.Is(err, service.ErrNotFound):
		Logger.Error(err)
		return response.ErrChatSessionNotFound()
	default:
		Logger.Error(err)
		return response.ErrUnknownError()
	}
}
