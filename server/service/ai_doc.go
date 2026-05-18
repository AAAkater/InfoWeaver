package service

import (
	"context"
	"server/models"
)

var AIDocServiceApp = new(AIDocService)

type AIDocService struct{}

// SplitDocument calls the AI document split service for a single file.
// It returns the split result or an error if the call fails.
func (this *AIDocService) SplitDocument(ctx context.Context, fileID, datasetID uint, minioPath string, chunkSize, chunkOverlap int) (*models.SplitDocumentResp, error) {
	return callAIServerSplit(ctx, fileID, datasetID, minioPath, chunkSize, chunkOverlap)
}

// EmbedDocument calls the AI document embedding service.
// It returns the embedding result or an error if the call fails.
func (this *AIDocService) EmbedDocument(ctx context.Context, chunkIDs []uint, embeddingCfg models.EmbeddingConfig) (*models.EmbeddingResp, error) {
	return callAIServerEmbed(ctx, models.EmbeddingReq{
		ChunkIDs:        chunkIDs,
		EmbeddingConfig: embeddingCfg,
	})
}
