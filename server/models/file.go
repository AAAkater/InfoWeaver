package models

import (
	"mime/multipart"
	"time"
)

type FileUploadReq struct {
	Files []multipart.FileHeader `form:"files" binding:"required"`
}

type FileUploadInfo struct {
	ID        uint   `json:"id"`
	OwnerID   uint   `json:"owner_id"`
	DatasetID uint   `json:"dataset_id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Size      int64  `json:"size"`
	MinioPath string `json:"-"`
}

type MultiFileUploadResp struct {
	Files []FileUploadInfo `json:"files"`
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

// FileUploadMessage represents the message sent to RabbitMQ when a file is uploaded
type FileUploadMessage struct {
	Event     string    `json:"event"`
	FileID    uint      `json:"file_id"`
	MinioPath string    `json:"minio_path"`
	Timestamp time.Time `json:"timestamp"`
	DatasetID uint      `json:"dataset_id"`
}

// SplitDocumentReq represents the request body for the document split API
type SplitDocumentReq struct {
	FileID       uint   `json:"file_id" validate:"required"`
	DatasetID    uint   `json:"dataset_id" validate:"required"`
	MinioPath    string `json:"minio_path" validate:"required"`
	ChunkSize    int    `json:"chunk_size" validate:"required,min=64,max=4096"`
	ChunkOverlap int    `json:"chunk_overlap" validate:"required,min=0,max=2048"`
}

// SplitDocumentResp represents the data returned by the document split API
type SplitDocumentResp struct {
	FileID      uint   `json:"file_id"`
	DatasetID   uint   `json:"dataset_id"`
	FileName    string `json:"file_name"`
	ChunksCount int    `json:"chunks_count"`
}
