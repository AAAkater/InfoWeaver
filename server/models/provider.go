package models

const (
	PROVIDER_MODE_OPENAI      = "openai"
	PROVIDER_MODE_OPENAI_RESP = "openai_response"
	PROVIDER_MODE_GEMINI      = "gemini"
	PROVIDER_MODE_ANTHROPIC   = "anthropic"
	PROVIDER_MODE_OLLAMA      = "ollama"
)

// ProviderInfo represents the configuration for a provider
type ProviderInfo struct {
	Name    string `json:"name"`
	Mode    string `json:"mode"`
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
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
