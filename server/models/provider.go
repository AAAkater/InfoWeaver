package models

import "time"

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

// ProviderAddModelsReq represents a request to add models to a provider
type ProviderAddModelsReq struct {
	ID              uint     `json:"id" validate:"required"`
	AvailableModels []string `json:"available_models" validate:"required,min=1,dive,min=1"`
}

// ModelInfo represents a model from a provider
type ModelInfo struct {
	ID      string `json:"id"`       // Model identifier (e.g., "gpt-4", "text-embedding-3-small")
	Object  string `json:"object"`   // Object type (usually "model")
	OwnedBy string `json:"owned_by"` // Owner of the model (e.g., "openai")
}

// ProviderModelsResp represents a list of models from a provider
type ProviderModelsResp struct {
	Models []ModelInfo `json:"models"`
}
