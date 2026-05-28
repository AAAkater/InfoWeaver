package service

import (
	"context"
	"fmt"
	"server/db"
	"server/models"
	"server/utils"
)

var AIDocServiceApp = new(AIDocService)

type AIDocService struct{}

// SplitDocument calls the AI document split service for a single file.
// It returns the split result or an error if the call fails.
func (this *AIDocService) SplitDocument(ctx context.Context, fileID, datasetID uint, minioPath string, chunkSize, chunkOverlap int) (*models.SplitDocumentResp, error) {
	return callAIServerSplit(ctx, fileID, datasetID, minioPath, chunkSize, chunkOverlap)
}

// EmbedDocument resolves the embedding configuration from the database (via the
// chunks' owning dataset and its provider), then calls the AI embedding service.
func (this *AIDocService) EmbedDocument(ctx context.Context, chunkIDs []uint, ownerID uint) (*models.EmbeddingResp, error) {
	// Step 1: get the owning dataset and provider for these chunks.
	type chunkFileInfo struct {
		FileID                uint
		DatasetID             uint
		DatasetEmbeddingModel string
		DatasetSearchType     string
		DatasetProviderID     uint
		ProviderBaseURL       string
		ProviderAPIKey        string
		ProviderMode          string
	}
	var info chunkFileInfo
	result := db.PgSqlDB.WithContext(ctx).
		Model(&models.Chunk{}).
		Select("files.id AS file_id, files.dataset_id, datasets.embedding_model AS dataset_embedding_model, datasets.search_type AS dataset_search_type, datasets.provider_id AS dataset_provider_id, providers.base_url AS provider_base_url, providers.api_key AS provider_api_key, providers.mode AS provider_mode").
		Joins("JOIN files ON files.id = chunks.file_id").
		Joins("JOIN datasets ON datasets.id = files.dataset_id").
		Joins("JOIN providers ON providers.id = datasets.provider_id").
		Where("chunks.id IN ? AND files.user_id = ?", chunkIDs, ownerID).
		First(&info)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to resolve embedding config: %w", result.Error)
	}

	// Step 2: decrypt the provider API key.
	apiKey, err := utils.DecryptAPIKey(info.ProviderAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt provider API key: %w", err)
	}

	// Step 3: build the embedding config and call the AI server.
	embeddingCfg := models.EmbeddingConfig{
		ModelName:    info.DatasetEmbeddingModel,
		BaseURL:      info.ProviderBaseURL,
		APIKey:       apiKey,
		ProviderType: info.ProviderMode,
		EmbedType:    info.DatasetSearchType,
	}
	return callAIServerEmbed(ctx, chunkIDs, embeddingCfg)
}
