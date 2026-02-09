package models

import (
	"mime/multipart"
	"time"
)

type FileUploadReq struct {
	File multipart.FileHeader `form:"file" binding:"required"`
}

type FileUploadResp struct {
	OwnerID   uint   `json:"owner_id"`
	DatasetID uint   `json:"dataset_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
}

type ListFilesReq struct {
	DatasetID uint `query:"dataset_id" validate:"required"`
	Page      int  `query:"page" binding:"required,min=1"`
	PageSize  int  `query:"page_size" binding:"required,min=1,max=100"`
}

type SimpleFileInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SimpleFileInfoListResp struct {
	Total int64            `json:"total"`
	Files []SimpleFileInfo `json:"files"`
}

type DetailedFileInfoReq struct {
	ID uint `param:"file_id" binding:"required"`
}

type FileDownloadResp struct {
	URL string `json:"url"`
}

type DetailedFileInfo struct {
	ID        uint
	CreatedAt time.Time
	Size      int64
	Name      string
	Type      string
	UserID    uint
}

type FileInfoUpdate struct {
	Size      int64
	Name      string
	Type      string
	MinioPath string
	UserID    uint
}
