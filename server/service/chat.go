package service

import (
	"context"
	"server/db"
	"server/models"

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
