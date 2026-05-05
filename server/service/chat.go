package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/config"
	"server/db"
	"server/models"
	"time"

	"gorm.io/gorm"
)

var ChatServiceApp = new(ChatService)

type ChatService struct{}

// CreateChatMessage creates a new chat message in a session
func (this *ChatService) CreateChatMessage(ctx context.Context, sessionID uint, content string, role string) error {
	dbMemory := models.Memory{
		SessionID: sessionID,
		Content:   content,
		Role:      role,
	}
	return gorm.G[models.Memory](db.PgSqlDB).Create(ctx, &dbMemory)
}

// GetChatMessageByID retrieves a specific chat message by ID
func (this *ChatService) GetChatMessageByID(ctx context.Context, id uint) (dbMemory *models.Memory, err error) {
	result := db.PgSqlDB.Model(&models.Memory{}).
		Where("id = ?", id).
		First(&dbMemory)

	return dbMemory, result.Error
}

// ListChatMessagesBySessionID retrieves all chat messages for a specific session in chronological order
func (this *ChatService) ListChatMessagesBySessionID(ctx context.Context, sessionID uint) (total int64, messages []models.ChatMessageInfo, err error) {
	result := db.PgSqlDB.Model(&models.Memory{}).
		Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages)

	return result.RowsAffected, messages, result.Error
}

// UpdateChatMessageContent updates the content of a chat message
func (this *ChatService) UpdateChatMessageContent(ctx context.Context, id uint, content string) error {
	rowsAffected, err := gorm.G[models.Memory](db.PgSqlDB).
		Where("id = ?", id).
		Updates(ctx, models.Memory{Content: content})

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

// DeleteChatMessage deletes a chat message by ID
func (this *ChatService) DeleteChatMessage(ctx context.Context, id uint) error {
	rowsAffected, err := gorm.G[models.Memory](db.PgSqlDB).
		Where("id = ?", id).
		Delete(ctx)

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

// DeleteAllChatMessagesBySessionID deletes all chat messages in a session
func (this *ChatService) DeleteAllChatMessagesBySessionID(ctx context.Context, sessionID uint) error {
	_, err := gorm.G[models.Memory](db.PgSqlDB).
		Where("session_id = ?", sessionID).
		Delete(ctx)
	return err
}

// chatStreamHTTPClient is a dedicated HTTP client for chat streaming requests.
// It uses a longer timeout than the default AIDocService client because SSE streams can be long-lived.
var chatStreamHTTPClient = &http.Client{
	Timeout: 5 * time.Minute,
}

// SendChatStreamToAIServer sends a chat streaming request to the AI server and returns the response body reader.
// The caller is responsible for closing the response body.
func (this *ChatService) SendChatStreamToAIServer(ctx context.Context, req models.SendChatStreamReq) (io.ReadCloser, error) {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal chat stream request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/v1/chat/chat/stream", config.Settings.GetAIServerDSN())
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create chat stream request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := chatStreamHTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call chat stream service: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("chat stream service returned status %d", resp.StatusCode)
	}

	return resp.Body, nil
}
