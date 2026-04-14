package service

import (
	"context"
	"server/db"
	"server/models"

	"gorm.io/gorm"
)

var ChatSessionServiceApp = new(ChatSessionService)

type ChatSessionService struct{}

// CreateNewChatSession creates a new chat session
func (this *ChatSessionService) CreateNewChatSession(ctx context.Context, title string, ownerID uint) error {
	dbChatSession := models.ChatSession{
		Title:   title,
		OwnerID: ownerID,
	}
	return gorm.G[models.ChatSession](db.PgSqlDB).Create(ctx, &dbChatSession)
}

// GetChatSessionInfoByID retrieves a chat session by ID and owner ID
func (this *ChatSessionService) GetChatSessionInfoByID(ctx context.Context, id uint, ownerID uint) (dbChatSession *models.ChatSessionInfo, err error) {
	result := db.PgSqlDB.Model(&models.ChatSession{}).
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&dbChatSession)

	return dbChatSession, result.Error
}

// ListChatSessionsByOwnerID retrieves all chat sessions for a specific user
func (this *ChatSessionService) ListChatSessionsByOwnerID(ctx context.Context, ownerID uint) (total int64, sessions []models.ChatSessionInfo, err error) {
	result := db.PgSqlDB.Model(&models.ChatSession{}).
		Where("owner_id = ?", ownerID).
		Find(&sessions)

	return result.RowsAffected, sessions, result.Error
}

// ListChatSessionsByTitle retrieves chat sessions by title with fuzzy matching
func (this *ChatSessionService) ListChatSessionsByTitle(ctx context.Context, ownerID uint, title string) (total int64, sessions []models.ChatSessionInfo, err error) {
	result := db.PgSqlDB.Model(&models.ChatSession{}).
		Where("title LIKE ? AND owner_id = ?", "%"+title+"%", ownerID).
		Find(&sessions)

	return result.RowsAffected, sessions, result.Error
}

// UpdateChatSession updates an existing chat session
func (this *ChatSessionService) UpdateChatSession(ctx context.Context, id uint, ownerID uint, title string) error {
	newChatSessionInfo := models.ChatSession{
		Title: title,
	}

	rowsAffected, err := gorm.G[models.ChatSession](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Updates(ctx, newChatSessionInfo)
	// id not found
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

// DeleteChatSession deletes a chat session by ID and owner ID
func (this *ChatSessionService) DeleteChatSession(ctx context.Context, id uint, ownerID uint) error {
	rowsAffected, err := gorm.G[models.ChatSession](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Delete(ctx)
	// id not found
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

// CheckChatSessionOwnership checks if a chat session belongs to a user
func (this *ChatSessionService) CheckChatSessionOwnership(ctx context.Context, id uint, ownerID uint) (exists bool, err error) {
	cnt, err := gorm.G[models.ChatSession](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Count(ctx, "*")
	return cnt > 0, err
}
