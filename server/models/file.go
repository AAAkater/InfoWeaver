package models

import "time"

type FileInfo struct {
	ID        uint
	CreatedAt time.Time
	Size      int64
	Name      string
	Type      string
	MinioPath string
	UserID    uint
}

type UpdateFileInfo struct {
	Size      int64
	Name      string
	Type      string
	MinioPath string
	UserID    uint
}
