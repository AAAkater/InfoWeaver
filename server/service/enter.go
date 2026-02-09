package service

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUnknown       = errors.New("Unknown error")
	ErrNotFound      = gorm.ErrRecordNotFound
	ErrDuplicatedKey = gorm.ErrDuplicatedKey
)
