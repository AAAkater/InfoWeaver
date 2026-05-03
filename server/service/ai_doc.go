package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server/config"
	"server/utils"
	"time"
)

var AIDocServiceApp = new(AIDocService)

type AIDocService struct {
	httpClient *http.Client
}

func (this *AIDocService) getHTTPClient() *http.Client {
	if this.httpClient == nil {
		this.httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return this.httpClient
}

// SplitDocumentReq represents the request to the AI document split service
type SplitDocumentReq struct {
	FileID       uint   `json:"file_id"`
	DatasetID    uint   `json:"dataset_id"`
	MinioPath    string `json:"minio_path"`
	ChunkSize    int    `json:"chunk_size"`
	ChunkOverlap int    `json:"chunk_overlap"`
}

// SplitDocumentResp represents the response from the AI document split service
type SplitDocumentResp struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data SplitDocumentData `json:"data"`
}

// SplitDocumentData represents the data field in the split response
type SplitDocumentData struct {
	FileID      uint   `json:"file_id"`
	DatasetID   uint   `json:"dataset_id"`
	FileName    string `json:"file_name"`
	ChunksCount int    `json:"chunks_count"`
}

// SplitDocument calls the AI document split service for a single file.
// It returns the split result or an error if the call fails.
func (this *AIDocService) SplitDocument(ctx context.Context, fileID, datasetID uint, minioPath string, chunkSize, chunkOverlap int) (*SplitDocumentData, error) {
	reqBody := SplitDocumentReq{
		FileID:       fileID,
		DatasetID:    datasetID,
		MinioPath:    minioPath,
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal split request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/v1/documents/split", config.Settings.GetAIServerDSN())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create split request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := this.getHTTPClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call split service: %w", err)
	}
	defer resp.Body.Close()

	var result SplitDocumentResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode split response: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("split service returned error: code=%d, msg=%s", result.Code, result.Msg)
	}

	utils.Logger.Infof("Document split completed: file_id=%d, chunks_count=%d",
		result.Data.FileID, result.Data.ChunksCount)

	return &result.Data, nil
}
