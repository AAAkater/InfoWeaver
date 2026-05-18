package service

import (
	"context"
	"fmt"
	"io"
	"server/config"
	"server/models"
	"server/utils"
	"sync"
	"time"

	"resty.dev/v3"
)

var (
	aiRestyOnce sync.Once
	aiResty     *resty.Client
)

// getAIServerClient returns the shared resty client, initialized on first call.
func getAIServerClient() *resty.Client {
	aiRestyOnce.Do(func() {
		aiResty = resty.New().
			SetBaseURL(config.Settings.GetAIServerDSN()).
			SetHeader("Content-Type", "application/json")
	})
	return aiResty
}

// aiServerPostJSON sends a JSON POST request to the AI server and
// unmarshals the response body into result. ctx controls the deadline.
func aiServerPostJSON(ctx context.Context, path string, reqBody, result any) error {
	resp, err := getAIServerClient().R().
		SetContext(ctx).
		SetBody(reqBody).
		SetResult(result).
		Post(path)
	if err != nil {
		return fmt.Errorf("failed to call AI server: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("AI server returned status %d", resp.StatusCode())
	}
	return nil
}

// aiServerPostStream sends a JSON POST to the AI server and returns the raw
// response body for streaming. The caller must close the body.
func aiServerPostStream(ctx context.Context, path string, reqBody any) (io.ReadCloser, error) {
	resp, err := getAIServerClient().R().
		SetContext(ctx).
		SetBody(reqBody).
		SetDoNotParseResponse(true). // keep body raw for streaming
		Post(path)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI server: %w", err)
	}
	if resp.StatusCode() != 200 {
		resp.RawResponse.Body.Close()
		return nil, fmt.Errorf("AI server returned status %d", resp.StatusCode())
	}
	return resp.RawResponse.Body, nil
}

// AI server path constants
const (
	aiPathDocSplit   = "/ai/v1/documents/split"
	aiPathDocEmbed   = "/ai/v1/documents/embedding"
	aiPathChatStream = "/ai/v1/chat/chat/stream"
)

// AIServerDocTimeout is the timeout for document operations (split, embedding).
const AIServerDocTimeout = 30 * time.Second

// withDocCtx wraps ctx with AIServerDocTimeout for short-lived AI server calls.
func withDocCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, AIServerDocTimeout)
}

// --- Convenience wrappers for AIDocService ---

// callAIServerSplit sends a split request with a doc-level timeout.
func callAIServerSplit(ctx context.Context, fileID, datasetID uint, minioPath string, chunkSize, chunkOverlap int) (*models.SplitDocumentResp, error) {
	ctx, cancel := withDocCtx(ctx)
	defer cancel()

	var result struct {
		Code int                      `json:"code"`
		Msg  string                   `json:"msg"`
		Data models.SplitDocumentResp `json:"data"`
	}
	if err := aiServerPostJSON(ctx, aiPathDocSplit, map[string]any{
		"file_id":       fileID,
		"dataset_id":    datasetID,
		"minio_path":    minioPath,
		"chunk_size":    chunkSize,
		"chunk_overlap": chunkOverlap,
	}, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("split service returned error: code=%d, msg=%s", result.Code, result.Msg)
	}
	utils.Logger.Infof("Document split completed: file_id=%d, chunks_count=%d",
		result.Data.FileID, result.Data.ChunksCount)
	return &result.Data, nil
}

// callAIServerEmbed sends an embedding request with a doc-level timeout.
func callAIServerEmbed(ctx context.Context, req models.EmbeddingReq) (*models.EmbeddingResp, error) {
	ctx, cancel := withDocCtx(ctx)
	defer cancel()

	var result struct {
		Code int                  `json:"code"`
		Msg  string               `json:"msg"`
		Data models.EmbeddingResp `json:"data"`
	}
	if err := aiServerPostJSON(ctx, aiPathDocEmbed, req, &result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("embedding service returned error: code=%d, msg=%s", result.Code, result.Msg)
	}
	utils.Logger.Infof("Document embedding completed: chunks_count=%d", result.Data.ChunksCount)
	return &result.Data, nil
}
