package models

import "time"

// ChatSessionCreateReq represents a request to create a new chat session
type ChatSessionCreateReq struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
}

// ChatSessionUpdateReq represents a request to update a chat session
type ChatSessionUpdateReq struct {
	ID    uint   `json:"id" validate:"required"`
	Title string `json:"title" validate:"required,min=1,max=200"`
}

// ChatSessionInfo represents the chat session information for API response
type ChatSessionInfo struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	OwnerID   uint      `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ChatSessionListResp represents a list of chat sessions for API response
type ChatSessionListResp struct {
	Total    int64             `json:"total"`
	Sessions []ChatSessionInfo `json:"sessions"`
}

// ChatSessionInfoReq represents a request to get/delete a specific chat session
type ChatSessionInfoReq struct {
	ID uint `param:"session_id" validate:"required"`
}

// ChatSessionListReq represents a request to list chat sessions with optional filters
type ChatSessionListReq struct {
	Title string `query:"title" validate:"omitempty,min=1,max=200"`
}
