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
