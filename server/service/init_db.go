package service

import (
	"server/config"
	"server/models"

	"gorm.io/gorm"
)

// InitializeDBTables auto-migrates the schema and creates initial admin account
func InitializeDBTables(db *gorm.DB) error {
	// Auto-migrate database schema
	if err := db.AutoMigrate(&models.User{}, &models.File{}, &models.Document{}, &models.Memory{}); err != nil {
		return err
	}

	// Create initial admin account if not exists
	if err := CreateInitialAdmin(db); err != nil {
		return err
	}

	return nil
}

// CreateInitialAdmin creates the initial admin account if none exists
func CreateInitialAdmin(db *gorm.DB) error {
	var admin models.User
	result := db.Where("role = ?", "admin").First(&admin)
	if result.Error == gorm.ErrRecordNotFound {
		// Create admin user with default credentials
		if err := db.Create(&models.User{
			Username: config.Settings.SYSTEM_ADMIN_NAME,
			Email:    config.Settings.SYSTEM_ADMIN_EMAIL,
			Password: config.Settings.SYSTEM_ADMIN_PASSWORD,
			Role:     "admin",
		}).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	}
	return nil
}
