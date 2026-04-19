package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ModelEnableMap is a custom type for storing model enable status as jsonb
type ModelEnableMap map[string]bool

// Value implements driver.Valuer interface for database storage
func (m ModelEnableMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements sql.Scanner interface for database retrieval
func (m *ModelEnableMap) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan ModelEnableMap: expected []byte")
	}
	return json.Unmarshal(bytes, m)
}

const (
	PROVIDER_MODE_OPENAI      = "openai"
	PROVIDER_MODE_OPENAI_RESP = "openai_response"
	PROVIDER_MODE_GEMINI      = "gemini"
	PROVIDER_MODE_ANTHROPIC   = "anthropic"
	PROVIDER_MODE_OLLAMA      = "ollama"
)

// ProviderInfo represents the configuration for a provider
type ProviderInfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Mode      string    `json:"mode"`
	BaseURL   string    `json:"base_url"`
}

type ProviderInfoReq struct {
	ID uint `param:"provider_id" validate:"required"`
}

// ProviderListResp represents a list of providers for API response
type ProviderListResp struct {
	Total     int64          `json:"total"`
	Providers []ProviderInfo `json:"providers"`
}

// ProviderCreateReq represents a request to create a new provider
type ProviderCreateReq struct {
	Name    string `json:"name" validate:"required,min=1,max=50"`
	BaseURL string `json:"base_url" validate:"required,url"`
	APIKey  string `json:"api_key" validate:"required,min=1"`
	Mode    string `json:"mode" validate:"required,oneof=openai openai_response gemini anthropic ollama"`
}

// ProviderUpdateReq represents a request to update a provider
type ProviderUpdateReq struct {
	ID      uint   `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required,min=1,max=50"`
	BaseURL string `json:"base_url" validate:"required,url"`
	APIKey  string `json:"api_key" validate:"required,min=1"`
	Mode    string `json:"mode" validate:"required,oneof=openai openai_response gemini anthropic ollama"`
}

// ProviderModelsReq represents a request to list models from a provider
type ProviderModelsReq struct {
	ID uint `param:"provider_id" validate:"required"`
}

// ProviderSetModelEnableReq represents a request to set enable status for a single model
type ProviderSetModelEnableReq struct {
	ID      uint   `json:"id" validate:"required"`       // Provider ID
	ModelID string `json:"model_id" validate:"required"` // Model ID to enable/disable
	Enabled bool   `json:"enabled" validate:"required"`  // Enable status
}

// ModelInfo represents a model from a provider
type ModelInfo struct {
	ID      string `json:"id"`       // Model identifier (e.g., "gpt-4", "text-embedding-3-small")
	Object  string `json:"object"`   // Object type (usually "model")
	OwnedBy string `json:"owned_by"` // Owner of the model (e.g., "openai")
	Enabled bool   `json:"enabled"`  // Enable status for this model
}

// ProviderModelsResp represents a list of models from a provider
type ProviderModelsResp struct {
	Models []ModelInfo `json:"models"`
}
