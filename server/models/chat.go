package models

import "time"

// Message role constants
const (
	ROLE_USER      = "user"
	ROLE_ASSISTANT = "assistant"
	ROLE_SYSTEM    = "system"
)

// ChatMessageCreateReq represents a request to send a new chat message
type ChatMessageCreateReq struct {
	SessionID uint   `json:"session_id" validate:"required"`
	Content   string `json:"content" validate:"required,min=1"`
	Role      string `json:"role" validate:"required,oneof=user assistant system"` // Message role
}

// ChatMessageInfo represents a single chat message for API response
type ChatMessageInfo struct {
	ID        uint      `json:"id"`
	SessionID uint      `json:"session_id"`
	Content   string    `json:"content"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatMessageListResp represents a list of chat messages for API response
type ChatMessageListResp struct {
	Total    int64             `json:"total"`
	Messages []ChatMessageInfo `json:"messages"`
}

// ChatSessionMessagesReq represents a request to get all messages in a session
type ChatSessionMessagesReq struct {
	SessionID uint `param:"session_id" validate:"required"`
}

// ChatMessageInfoReq represents a request to get/delete a specific message
type ChatMessageInfoReq struct {
	ID uint `param:"message_id" validate:"required"`
}

// SamplingParams represents LLM sampling parameters
type SamplingParams struct {
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
	MaxTokens        int     `json:"max_tokens"`
	PresencePenalty  float64 `json:"presence_penalty"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
}

// LLMConfig represents the LLM model configuration for the chat request
type LLMConfig struct {
	ModelName      string         `json:"model_name"`
	APIKey         string         `json:"api_key"`
	BaseURL        string         `json:"base_url"`
	ProviderType   string         `json:"provider_type"`
	SamplingParams SamplingParams `json:"sampling_params"`
}

// RetrievalConfig represents the retrieval configuration for RAG
type RetrievalConfig struct {
	TopK int `json:"top_k"`
}

// SendChatStreamReq represents a request to send a chat message with streaming LLM response
type SendChatStreamReq struct {
	Query           string          `json:"query" validate:"required,min=1"`
	DatasetID       uint            `json:"dataset_id"`
	SessionID       uint            `json:"session_id" validate:"required"`
	LLMConfig       LLMConfig       `json:"llm_config" validate:"required"`
	EmbeddingConfig EmbeddingConfig `json:"embedding_config" validate:"required"` // Defined in file.go
	RetrievalConfig RetrievalConfig `json:"retrieval_config" validate:"required"`
	SystemPrompt    string          `json:"system_prompt"`
}

// ChatStreamChunk represents a single SSE chunk from the AI chat stream service
type ChatStreamChunk struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
