package v1

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/common/response"
	"server/service"
	"server/utils"
	"strings"

	"github.com/labstack/echo/v5"
)

func SetChatRouter(e *echo.Echo) {
	chatRouterGroup := e.Group(config.API_V1+"/chat", middleware.TokenMiddleware())
	chatHandler := &chatApi{}
	chatRouterGroup.POST("/send", chatHandler.sendChatMessage)
	chatRouterGroup.GET("/session/:session_id", chatHandler.listSessionMessages)
	chatRouterGroup.GET("/message/:message_id", chatHandler.getChatMessageInfo)
	chatRouterGroup.POST("/update/:message_id", chatHandler.updateChatMessageContent)
	chatRouterGroup.POST("/delete/:message_id", chatHandler.deleteChatMessage)
	chatRouterGroup.POST("/clear/:session_id", chatHandler.clearSessionMessages)
}

type chatApi struct{}

// sendChatMessage godoc
//
//	@Summary		Send Chat Message (Streaming)
//	@Description	Send a new chat message in a session with streaming LLM response via SSE
//	@Tags			Chat
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			message	body		models.SendChatStreamReq	true	"Chat stream request with LLM/RAG config"
//	@Success		200		{object}	response.ResponseBase[any]	"SSE streaming response"
//	@Failure		400		{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401		{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403		{object}	response.ResponseBase[any]	"Session not owned by user"
//	@Failure		404		{object}	response.ResponseBase[any]	"Session not found"
//	@Failure		500		{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat/send [post]
func (this *chatApi) sendChatMessage(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.SendChatStreamReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Verify session belongs to the user
	if owned, err := chatSessionService.CheckChatSessionOwnership(ctx.Request().Context(), args.SessionID,
		currentUser.ID); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	} else if !owned {
		return response.ErrChatSessionNotFound()
	}

	// Save the user message first
	if err := chatService.CreateChatMessage(ctx.Request().Context(), args.SessionID, args.Query, models.ROLE_USER); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}

	// Call the AI streaming service
	streamBody, err := chatService.SendChatStreamToAIServer(ctx.Request().Context(), *args)
	if err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	defer streamBody.Close()

	// Set up SSE response headers
	w := ctx.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	// Stream chunks from AI server to client, accumulating the full response text
	flusher, _ := w.(http.Flusher)
	scanner := bufio.NewScanner(streamBody)
	var fullResponse strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Forward the SSE line to the client
		fmt.Fprintf(w, "%s\n", line)
		if flusher != nil {
			flusher.Flush()
		}

		// Check if the client has disconnected — save partial response and exit
		select {
		case <-ctx.Request().Context().Done():
			assistantContent := fullResponse.String()
			if assistantContent != "" {
				if saveErr := chatService.CreateChatMessage(context.Background(), args.SessionID, assistantContent, models.ROLE_ASSISTANT); saveErr != nil {
					Logger.Error("failed to save assistant message on disconnect: ", saveErr)
				}
			}
			return nil
		default:
		}

		// Parse SSE data lines to accumulate the assistant's response content
		if after, ok := strings.CutPrefix(line, "data:"); ok {
			data := after
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				break
			}
			var chunk models.ChatStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err == nil {
				fullResponse.WriteString(chunk.Content)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		Logger.Error("error reading chat stream: ", err)
	}

	// Save the assistant message to the database
	assistantContent := fullResponse.String()
	if assistantContent != "" {
		if err := chatService.CreateChatMessage(ctx.Request().Context(), args.SessionID, assistantContent, models.ROLE_ASSISTANT); err != nil {
			Logger.Error("failed to save assistant message: ", err)
		}
	}

	return nil
}

// listSessionMessages godoc
//
//	@Summary		List Session Messages
//	@Description	Get all chat messages in a session in chronological order
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			session_id	path		int													true	"Chat session ID"
//	@Success		200			{object}	response.ResponseBase[models.ChatMessageListResp]	"List of chat messages"
//	@Failure		400			{object}	response.ResponseBase[any]							"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]							"Invalid or expired token"
//	@Failure		403			{object}	response.ResponseBase[any]							"Session not owned by user"
//	@Failure		404			{object}	response.ResponseBase[any]							"Session not found"
//	@Failure		500			{object}	response.ResponseBase[any]							"Internal server error"
//	@Router			/chat/session/{session_id} [get]
func (this *chatApi) listSessionMessages(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatSessionMessagesReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}
	// Verify session belongs to the user
	if owned, err := chatSessionService.CheckChatSessionOwnership(ctx.Request().Context(), args.SessionID,
		currentUser.ID); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	} else if !owned {
		return response.ErrChatSessionNotFound()
	}

	total, messages, err := chatService.ListChatMessagesBySessionID(ctx.Request().Context(), args.SessionID)
	if err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.OkWithData(ctx, models.ChatMessageListResp{
		Total:    total,
		Messages: messages,
	})
}

// getChatMessageInfo godoc
//
//	@Summary		Get Chat Message
//	@Description	Get a specific chat message by ID
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			message_id	path		int												true	"Chat message ID"
//	@Success		200			{object}	response.ResponseBase[models.ChatMessageInfo]	"Chat message details"
//	@Failure		400			{object}	response.ResponseBase[any]						"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]						"Invalid or expired token"
//	@Failure		403			{object}	response.ResponseBase[any]						"Message not owned by user"
//	@Failure		404			{object}	response.ResponseBase[any]						"Message not found"
//	@Failure		500			{object}	response.ResponseBase[any]						"Internal server error"
//	@Router			/chat/message/{message_id} [get]
func (this *chatApi) getChatMessageInfo(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatMessageInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get the message
	message, err := chatService.GetChatMessageByID(ctx.Request().Context(), args.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.ErrChatMessageNotFound()
		}
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	if message == nil {
		return response.ErrChatMessageNotFound()
	}

	// Verify the session belongs to the user
	owned, err := chatSessionService.CheckChatSessionOwnership(ctx.Request().Context(), message.SessionID, currentUser.ID)
	if err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	if !owned {
		return response.ErrChatMessageNotOwned()
	}

	return response.OkWithData(ctx, models.ChatMessageInfo{
		ID:        message.ID,
		SessionID: message.SessionID,
		Content:   message.Content,
		Role:      message.Role,
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	})
}

// updateChatMessageContent godoc
//
//	@Summary		Update Chat Message Content
//	@Description	Update the content of a chat message
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			message_id	path		int							true	"Chat message ID"
//	@Param			content		body		object						true	"Content to update"
//	@Success		200			{object}	response.ResponseBase[any]	"Message updated successfully"
//	@Failure		400			{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403			{object}	response.ResponseBase[any]	"Message not owned by user"
//	@Failure		404			{object}	response.ResponseBase[any]	"Message not found"
//	@Failure		500			{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat/update/{message_id} [post]
func (this *chatApi) updateChatMessageContent(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatMessageInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get the message to verify ownership
	message, err := chatService.GetChatMessageByID(ctx.Request().Context(), args.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.ErrChatMessageNotFound()
		}
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	if message == nil {
		return response.ErrChatMessageNotFound()
	}

	// Verify the session belongs to the user
	owned, err := chatSessionService.CheckChatSessionOwnership(ctx.Request().Context(), message.SessionID, currentUser.ID)
	if err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	if !owned {
		return response.ErrChatMessageNotOwned()
	}

	// Parse the content from request body
	var body struct {
		Content string `json:"content"`
	}
	if err := ctx.Bind(&body); err != nil {
		return response.BadRequestWithMsg("Invalid request body")
	}

	if err := chatService.UpdateChatMessageContent(ctx.Request().Context(), args.ID, body.Content); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.ErrChatMessageNotFound()
		}
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.Ok(ctx)
}

// deleteChatMessage godoc
//
//	@Summary		Delete Chat Message
//	@Description	Delete a specific chat message by ID
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			message_id	path		int							true	"Chat message ID"
//	@Success		200			{object}	response.ResponseBase[any]	"Message deleted successfully"
//	@Failure		400			{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403			{object}	response.ResponseBase[any]	"Message not owned by user"
//	@Failure		404			{object}	response.ResponseBase[any]	"Message not found"
//	@Failure		500			{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat/delete/{message_id} [post]
func (this *chatApi) deleteChatMessage(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatMessageInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get the message to verify ownership
	message, err := chatService.GetChatMessageByID(ctx.Request().Context(), args.ID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.ErrChatMessageNotFound()
		}
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	if message == nil {
		return response.ErrChatMessageNotFound()
	}

	// Verify session belongs to the user
	if owned, err := chatSessionService.CheckChatSessionOwnership(ctx.Request().Context(), args.ID,
		currentUser.ID); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	} else if !owned {
		return response.ErrChatSessionNotFound()
	}

	if err := chatService.DeleteChatMessage(ctx.Request().Context(), args.ID); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return response.ErrChatMessageNotFound()
		}
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.Ok(ctx)
}

// clearSessionMessages godoc
//
//	@Summary		Clear Session Messages
//	@Description	Delete all chat messages in a session
//	@Tags			Chat
//	@Accept			json
//	@Produce		json
//	@Param			session_id	path		int							true	"Chat session ID"
//	@Success		200			{object}	response.ResponseBase[any]	"All messages cleared successfully"
//	@Failure		400			{object}	response.ResponseBase[any]	"Invalid request parameters"
//	@Failure		401			{object}	response.ResponseBase[any]	"Invalid or expired token"
//	@Failure		403			{object}	response.ResponseBase[any]	"Session not owned by user"
//	@Failure		404			{object}	response.ResponseBase[any]	"Session not found"
//	@Failure		500			{object}	response.ResponseBase[any]	"Internal server error"
//	@Router			/chat/clear/{session_id} [post]
func (this *chatApi) clearSessionMessages(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}
	args, err := utils.BindAndValidate[models.ChatSessionMessagesReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Verify session belongs to the user
	if owned, err := chatSessionService.CheckChatSessionOwnership(ctx.Request().Context(), args.SessionID,
		currentUser.ID); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	} else if !owned {
		return response.ErrChatSessionNotFound()
	}

	if err := chatService.DeleteAllChatMessagesBySessionID(ctx.Request().Context(), args.SessionID); err != nil {
		Logger.Error(err)
		return response.ErrUnknownError()
	}
	return response.Ok(ctx)
}
