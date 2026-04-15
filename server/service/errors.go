package service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUnknown       = errors.New("Unknown error")
	ErrNotFound      = gorm.ErrRecordNotFound
	ErrDuplicatedKey = gorm.ErrDuplicatedKey
	ErrOpenFile      = errors.New("Failed to open file")
	ErrUploadFile    = errors.New("Failed to upload file to MinIO")
	ErrSaveFileInfo  = errors.New("Failed to save file record to database")
)
