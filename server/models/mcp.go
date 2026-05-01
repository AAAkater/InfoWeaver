package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JsonMap is a custom type for storing arbitrary JSON objects as jsonb
type JsonMap map[string]any

// Value implements driver.Valuer interface for database storage
func (m JsonMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

// Scan implements sql.Scanner interface for database retrieval
func (m *JsonMap) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JsonMap: expected []byte")
	}
	return json.Unmarshal(bytes, m)
}

// McpInfo represents the MCP server configuration for API response (without sensitive fields)
type McpInfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Transport string    `json:"transport"`
	Command   string    `json:"command,omitempty"`
	Args      string    `json:"args,omitempty"`
	URL       string    `json:"url,omitempty"`
	Headers   JsonMap   `json:"headers,omitempty"`
	EnvVars   JsonMap   `json:"env_vars,omitempty"`
	Enabled   bool      `json:"enabled"`
}

// McpCreateReq represents a request to create a new MCP server
type McpCreateReq struct {
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Transport string  `json:"transport" validate:"required,oneof=stdio sse streamable_http"`
	Command   string  `json:"command" validate:"required_if=Transport stdio"`
	Args      string  `json:"args"`
	URL       string  `json:"url" validate:"required_if=Transport sse,required_if=Transport streamable_http"`
	Headers   JsonMap `json:"headers"`
	EnvVars   JsonMap `json:"env_vars"`
	Enabled   *bool   `json:"enabled"`
}

// McpUpdateReq represents a request to update an existing MCP server
type McpUpdateReq struct {
	ID        uint    `json:"id" validate:"required"`
	Name      string  `json:"name" validate:"required,min=1,max=100"`
	Transport string  `json:"transport" validate:"required,oneof=stdio sse streamable_http"`
	Command   string  `json:"command" validate:"required_if=Transport stdio"`
	Args      string  `json:"args"`
	URL       string  `json:"url" validate:"required_if=Transport sse,required_if=Transport streamable_http"`
	Headers   JsonMap `json:"headers"`
	EnvVars   JsonMap `json:"env_vars"`
	Enabled   *bool   `json:"enabled"`
}

// McpInfoReq represents a request to get/delete a specific MCP server by ID
type McpInfoReq struct {
	ID uint `param:"mcp_id" validate:"required"`
}

// McpListReq represents a request to list MCP servers with pagination
type McpListReq struct {
	Page     int `query:"page" validate:"omitempty,min=1"`
	PageSize int `query:"page_size" validate:"omitempty,min=1,max=100"`
}

// McpListResp represents a list of MCP servers for API response
type McpListResp struct {
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Mcps  []McpInfo `json:"mcps"`
}
