package models

import (
	"mime/multipart"
	"time"
)

type FileUploadReq struct {
	File multipart.FileHeader `form:"file" binding:"required"`
}

type FileUploadResp struct {
	OwnerID uint   `json:"owner_id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
}

type FileInfoListReq struct {
	Page     int `query:"page" binding:"required,min=1"`
	PageSize int `query:"page_size" binding:"required,min=1,max=100"`
}

type FileInfo struct {
	ID        uint
	CreatedAt time.Time
	Size      int64
	Name      string
	Type      string
	MinioPath string
	UserID    uint
}

type FileInfoUpdate struct {
	Size      int64
	Name      string
	Type      string
	MinioPath string
	UserID    uint
}
