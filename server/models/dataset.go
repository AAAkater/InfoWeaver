package models

import "time"

type DatasetCreateReq struct {
	Icon           string `json:"icon" validate:"required,emoji"`
	Name           string `json:"name" validate:"required,min=1,max=100"`
	Description    string `json:"description" validate:"max=500"`
	SearchType     string `json:"search_type" validate:"required,oneof=sparse dense hybrid"`
	EmbeddingModel string `json:"embedding_model" validate:"required"`
	ProviderID     uint   `json:"provider_id" validate:"required"`
}

type DatasetUpdateReq struct {
	ID             uint   `json:"id" validate:"required"`
	Icon           string `json:"icon" validate:"emoji"`
	Name           string `json:"name" validate:"required,min=1,max=100"`
	Description    string `json:"description" validate:"max=500"`
	SearchType     string `json:"search_type" validate:"omitempty,oneof=sparse dense hybrid"`
	EmbeddingModel string `json:"embedding_model" validate:"omitempty"`
	ProviderID     uint   `json:"provider_id" validate:"omitempty"`
}

type DatasetInfo struct {
	ID             uint   `json:"id"`
	Icon           string `json:"icon"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	SearchType     string `json:"search_type"`
	EmbeddingModel string `json:"embedding_model"`
	ProviderID     uint   `json:"provider_id"`
	OwnerID        uint   `json:"owner_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type DatasetListResp struct {
	Total    int64         `json:"total"`
	Datasets []DatasetInfo `json:"datasets"`
}

type DatasetInfoReq struct {
	ID uint `param:"dataset_id" validate:"required"`
}

type DatasetListReq struct {
	Name string `query:"name" validate:"omitempty,min=1,max=100"`
}

// --- Chunk list under a dataset ---

type DatasetChunkListReq struct {
	DatasetID uint `param:"dataset_id" validate:"required"`
	Page      int  `query:"page" validate:"omitempty,min=1"`
	PageSize  int  `query:"page_size" validate:"omitempty,min=1,max=100"`
}

type ChunkInfo struct {
	ID            uint           `json:"id"`
	Content       string         `json:"content"`
	ChunkMetadata map[string]any `json:"chunk_metadata"`
	Status        string         `json:"status"`
	VectorID      string         `json:"vector_id"`
	FileID        uint           `json:"file_id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type DatasetChunkListResp struct {
	Total  int64       `json:"total"`
	Chunks []ChunkInfo `json:"chunks"`
}
