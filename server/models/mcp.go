package models

import "time"

// McpInfo represents the MCP server configuration for API response (without sensitive fields)
type McpInfo struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Name      string         `json:"name"`
	Transport string         `json:"transport"`
	Command   string         `json:"command,omitempty"`
	Args      string         `json:"args,omitempty"`
	URL       string         `json:"url,omitempty"`
	Headers   map[string]any `json:"headers,omitempty"`
	EnvVars   map[string]any `json:"env_vars,omitempty"`
	Enabled   bool           `json:"enabled"`
}

// McpCreateReq represents a request to create a new MCP server
type McpCreateReq struct {
	Name      string         `json:"name" validate:"required,min=1,max=100"`
	Transport string         `json:"transport" validate:"required,oneof=stdio sse streamable_http"`
	Command   string         `json:"command" validate:"required_if=Transport stdio"`
	Args      string         `json:"args"`
	URL       string         `json:"url" validate:"required_if=Transport sse,required_if=Transport streamable_http"`
	Headers   map[string]any `json:"headers"`
	EnvVars   map[string]any `json:"env_vars"`
	Enabled   *bool          `json:"enabled"`
}

// McpUpdateReq represents a request to update an existing MCP server
type McpUpdateReq struct {
	ID        uint           `json:"id" validate:"required"`
	Name      string         `json:"name" validate:"required,min=1,max=100"`
	Transport string         `json:"transport" validate:"required,oneof=stdio sse streamable_http"`
	Command   string         `json:"command" validate:"required_if=Transport stdio"`
	Args      string         `json:"args"`
	URL       string         `json:"url" validate:"required_if=Transport sse,required_if=Transport streamable_http"`
	Headers   map[string]any `json:"headers"`
	EnvVars   map[string]any `json:"env_vars"`
	Enabled   *bool          `json:"enabled"`
}

// McpInfoReq represents a request to get/delete a specific MCP server by ID
type McpInfoReq struct {
	ID uint `param:"mcp_id" validate:"required"`
}

// McpListResp represents a list of MCP servers for API response
type McpListResp struct {
	Total int64     `json:"total"`
	Mcps  []McpInfo `json:"mcps"`
}
