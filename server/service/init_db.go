package service

import (
	"server/models"

	"gorm.io/gorm"
)

// InitializeDBTables auto-migrates the schema
func InitializeDBTables(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.File{}, &models.Document{}, &models.Memory{})
}
